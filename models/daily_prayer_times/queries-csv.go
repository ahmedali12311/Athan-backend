package daily_prayer_times

import (
	"app/models/city"
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	finder "bitbucket.org/sadeemTechnology/backend-finder"
	validator "bitbucket.org/sadeemTechnology/backend-validator"
	"github.com/google/uuid"
)

var ErrValidation = errors.New("validation error")

func (m *Queries) ParseCSV(
	v *validator.Validator,
	tx finder.Connection,
) (*[]Model, []string, error) {
	prayerTimes := &[]Model{}
	var entryLog []string

	headers := []string{
		"Day (1-31)",
		"Fajr 1 Time",
		"Fajr 2 Time",
		"Sunrise Time",
		"Dhuhr Time",
		"Asr Time",
		"Maghrib Time",
		"Isha Time",
		"Month (1-12)",
	}

	records, err := v.ParseCSV("csv", headers, true)
	if err != nil {
		v.Check(false, "csv", err.Error())
		return nil, nil, err
	}

	tempModel := Model{City: city.MinimalModel{}}
	v.UnmarshalInto("city", &tempModel.City)

	var cityID uuid.UUID
	if tempModel.City.ID != nil {
		v.UUIDExistsInDB(tempModel.City.ID, "city_id", "id", "cities", true)
		cityID = *tempModel.City.ID
	} else {
		v.Check(false, "city", "must enter a city!")
	}
	if !v.Valid() {
		return nil, nil, ErrValidation
	}

	for i, r := range records {
		day, dayErr := strconv.Atoi(r[0])
		month, monthErr := strconv.Atoi(r[8])

		if dayErr != nil || day < 1 || day > 31 {
			v.Check(false, fmt.Sprintf("day.%d", i+1), "Invalid day (must be 1-31).")
			entryLog = append(entryLog, fmt.Sprintf("Row %d: Invalid day value", i+1))
			continue
		}
		if monthErr != nil || month < 1 || month > 12 {
			v.Check(false, fmt.Sprintf("month.%d", i+1), "Invalid month (must be 1-12).")
			entryLog = append(entryLog, fmt.Sprintf("Row %d: Invalid month value", i+1))
			continue
		}

		pt := Model{
			ID:             uuid.New(),
			CityID:         cityID,
			Day:            day,
			Month:          month,
			FajrFirstTime:  r[1],
			FajrSecondTime: r[2],
			SunriseTime:    r[3],
			DhuhrTime:      r[4],
			AsrTime:        r[5],
			MaghribTime:    r[6],
			IshaTime:       r[7],
			CreatedAt:      time.Now(),
		}

		*prayerTimes = append(*prayerTimes, pt)
	}

	if !v.Valid() {
		return nil, entryLog, ErrValidation
	}

	return prayerTimes, entryLog, nil
}

func (m *Queries) BulkCreate(
	prayerTimes *[]Model,
	_ *validator.Validator,
	conn finder.Connection,

) (*[]Model, error) {
	inserts := m.QB.
		Insert("daily_prayer_times").
		Columns("id", "city_id", "day", "month",
			"fajr_first_time", "fajr_second_time", "sunrise_time",
			"dhuhr_time", "asr_time", "maghrib_time", "isha_time", "created_at",
		)
	for _, itm := range *prayerTimes {
		inserts = inserts.Values(
			itm.ID,
			itm.CityID,
			itm.Day,
			itm.Month,
			itm.FajrFirstTime,
			itm.FajrSecondTime,
			itm.SunriseTime,
			itm.DhuhrTime,
			itm.AsrTime,
			itm.MaghribTime,
			itm.IshaTime,
			itm.CreatedAt,
		)
	}
	query, args, err := inserts.ToSql()
	if err != nil {
		return nil, err
	}
	if _, err := conn.ExecContext(
		context.Background(),
		query,
		args...,
	); err != nil {
		return nil, err
	}

	return prayerTimes, nil
}
