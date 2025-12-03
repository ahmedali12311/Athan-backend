package daily_prayer_times_controller

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"app/controller"
	dailyprayertimes "app/models/daily-prayer-times"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ControllerBasic struct {
	*controller.Dependencies
}

// Scopes ---------------------------------------------------------------------

func (c *ControllerBasic) scope(ctx echo.Context) *dailyprayertimes.WhereScope {
	scopes := c.Utils.CtxScopes(ctx)

	var admin bool
	for _, v := range scopes {
		switch v {
		case "admin":
			admin = true
		}
	}

	ws := &dailyprayertimes.WhereScope{
		IsAdmin: admin,
	}

	queryParams := ctx.QueryParams()

	if cityIDStr := queryParams.Get("city_id"); cityIDStr != "" {
		cityID, err := uuid.Parse(cityIDStr)
		if err == nil {
			ws.CityID = &cityID
		}
	}
	if dayStr := queryParams.Get("day"); dayStr != "" {
		dayVal, err := strconv.ParseInt(dayStr, 10, 64)
		if err == nil {
			dayInt := int(dayVal)
			if dayInt >= 1 && dayInt <= 31 {
				ws.Day = &dayInt
			}
		}
	}

	if monthStr := queryParams.Get("month"); monthStr != "" {
		monthVal, err := strconv.ParseInt(monthStr, 10, 64)
		if err == nil {
			monthInt := int(monthVal)
			if monthInt >= 1 && monthInt <= 12 {
				ws.Month = &monthInt
			}
		}
	}
	// Filter by Date (e.g., ?date=2025-12-02)
	if dateStr := queryParams.Get("date"); dateStr != "" {
		t, err := time.Parse(time.DateOnly, dateStr)
		if err == nil {
			ws.Date = &t
		}
	}

	return ws
}

func (c *ControllerBasic) Index(ctx echo.Context) error {
	ws := c.scope(ctx)
	indexResponse, err := c.Models.DailyPrayerTimes.GetAll(ctx, ws)
	if err != nil {
		return c.APIErr.Database(ctx, err, nil)
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *ControllerBasic) Show(ctx echo.Context) error {
	var result dailyprayertimes.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}

	ws := c.scope(ctx)
	if err := c.Models.DailyPrayerTimes.GetOne(&result, ws); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Store(ctx echo.Context) error {
	var result dailyprayertimes.Model

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	if valid := result.MergeAndValidate(v); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}

	if !v.Valid() {
		return c.APIErr.InputValidation(ctx, v)
	}

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.DailyPrayerTimes.CreateOne(&result, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, result)
}

func (c *ControllerBasic) Destroy(ctx echo.Context) error {
	var result dailyprayertimes.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	ws := c.scope(ctx)

	if err := c.Models.DailyPrayerTimes.DeleteOne(&result, ws, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	} else {
		if err := tx.Commit(); err != nil {
			return c.APIErr.InternalServer(ctx, err)
		}
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Update(ctx echo.Context) error {
	t := c.Utils.CtxT(ctx)
	var result dailyprayertimes.Model

	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}

	ws := c.scope(ctx)
	if err := c.Models.DailyPrayerTimes.GetOne(&result, ws); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	if valid := result.MergeAndValidate(v); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}

	// Start transacting
	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	if err := c.Models.DailyPrayerTimes.UpdateOne(&result, ws, tx); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.APIErr.Database(
				ctx,
				errors.New(t.ConflictError()),
				&result,
			)
		default:
			return c.APIErr.Database(ctx, err, &result)
		}
	}
	if err := tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, result)
}
