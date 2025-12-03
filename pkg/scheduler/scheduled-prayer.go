package scheduler

import (
	"app/models/fcm_notification"
	"app/pkg/firebase_utils"
	"context"
	"fmt"
	"strings"
	"time"
)

const (
	prayerJobName = "scheduled-prayer-notifications"
	libyaLocation = "Africa/Tripoli"
)

var prayerTimeColumns = map[string]string{
	"fajr":    "fajr_first_time",
	"dhuhr":   "dhuhr_time",
	"asr":     "asr_time",
	"maghrib": "maghrib_time",
	"isha":    "isha_time",
}
var prayerNameArabic = map[string]string{
	"fajr":    "الفجر",
	"dhuhr":   "الظهر",
	"asr":     "العصر",
	"maghrib": "المغرب",
	"isha":    "العشاء",
}

func ScheduledPrayerNotifications(cfg *Config) {
	startTime := time.Now()

	cfg.Logger.Info().Msgf("%s: started", prayerJobName)

	defer func() {
		cfg.Logger.Info().
			Str("duration", time.Since(startTime).String()).
			Msgf("%s: completed", prayerJobName)
	}()

	loc, err := time.LoadLocation(libyaLocation)
	if err != nil {
		cfg.Logger.Error().Err(err).Msg("Failed to load Libya timezone")
		return
	}

	now := time.Now().In(loc)
	todayDay := now.Day()
	todayMonth := int(now.Month())
	targetTime := now.Format("15:04:00")

	for techName, dbColumn := range prayerTimeColumns {

		arabicName := prayerNameArabic[techName]

		citiesDue, err := cfg.Models.DailyPrayerTimes.GetDueCities(
			context.Background(), dbColumn, todayDay, todayMonth, targetTime,
		)

		if err != nil {
			cfg.Logger.Error().Err(err).Msgf("%s: Failed to query cities for %s. Error: %v", prayerJobName, techName, err)
			continue
		}

		if len(citiesDue) == 0 {
			continue
		}

		for _, city := range citiesDue {
			topicName := strings.ToLower(fmt.Sprintf("%s_%s", strings.ReplaceAll(city.Name, " ", "_"), techName))

			title := fmt.Sprintf("حان الآن موعد صلاة %s", arabicName)
			body := fmt.Sprintf("%s (%s) في %s", arabicName, city.Time[:5], city.Name)
			// Data Payload
			data := map[string]string{
				"prayer": techName,
				"city":   city.Name,
				"time":   city.Time,
				"topic":  topicName,
			}
			sendTopicNotification(
				cfg,
				topicName,
				title,
				body,
				data,
			)
		}
	}
}

func sendTopicNotification(
	cfg *Config,
	topic, title, body string,
	data map[string]string,
) {
	message := firebase_utils.BuildTopicMessage(
		&title,
		&body,
		&topic,
		data,
	)
	response, err := cfg.Utils.FBM.Send(context.Background(), message)

	tx, logErr := cfg.DB.Beginx()
	if logErr != nil {
		cfg.Logger.Error().Err(logErr).Str("topic", topic).Msg("failed to start transaction for logging")
		return
	}
	defer tx.Rollback()

	fcmRecord := fcm_notification.Model{
		Title:  title,
		Body:   body,
		Data:   data,
		Topic:  &topic,
		IsSent: true,
	}

	if err != nil {
		fcmRecord.IsSent = false
		errorResponse := err.Error()
		fcmRecord.Response = &errorResponse
		cfg.Logger.Error().Err(err).Str("topic", topic).Msg("failed to send FCM message")
	} else {
		fcmRecord.Response = &response
	}

	if logErr := cfg.Models.FcmNotification.CreateOne(&fcmRecord); logErr != nil {
		cfg.Logger.Error().Err(logErr).Str("topic", topic).Msg("failed to create FCM record")
		return
	}

	if logErr := tx.Commit(); logErr != nil {
		cfg.Logger.Error().Err(logErr).Str("topic", topic).Msg("failed to commit log transaction")
	}
}
