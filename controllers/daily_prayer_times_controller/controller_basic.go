package daily_prayer_times_controller

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"app/controller"
	"app/models/daily_prayer_times"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ControllerBasic struct {
	*controller.Dependencies
}

// Scopes ---------------------------------------------------------------------

func (c *ControllerBasic) scope(ctx echo.Context) *daily_prayer_times.WhereScope {
	scopes := c.Utils.CtxScopes(ctx)
	ctxUser := c.Utils.CtxUser(ctx)

	var admin, public bool
	for _, v := range scopes {
		switch v {
		case "admin":
			admin = true
		case "public":
			public = true
		}
	}

	ws := &daily_prayer_times.WhereScope{
		IsAdmin:     admin,
		IsPublic:    public && !admin,
		QueryParams: ctx.QueryParams(),
	}

	if ctxUser != nil {
		ws.UserID = &ctxUser.ID
	}

	if fd := ctx.QueryParam("from_day"); fd != "" {
		if d, err := strconv.Atoi(fd); err == nil && d >= 1 && d <= 31 {
			ws.FromDay = &d
		}
	}
	if fm := ctx.QueryParam("from_month"); fm != "" {
		if m, err := strconv.Atoi(fm); err == nil && m >= 1 && m <= 12 {
			ws.FromMonth = &m
		}
	}

	// Parse to_day, to_month
	if td := ctx.QueryParam("to_day"); td != "" {
		if d, err := strconv.Atoi(td); err == nil && d >= 1 && d <= 31 {
			ws.ToDay = &d
		}
	}
	if tm := ctx.QueryParam("to_month"); tm != "" {
		if m, err := strconv.Atoi(tm); err == nil && m >= 1 && m <= 12 {
			ws.ToMonth = &m
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
	var result daily_prayer_times.Model
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
	// ws := c.scope(ctx)
	var result daily_prayer_times.Model

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}

	if valid := result.MergeAndValidate(v); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}

	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := c.Models.DailyPrayerTimes.CreateOne(&result, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, result)
}

func (c *ControllerBasic) Update(ctx echo.Context) error {
	var result daily_prayer_times.Model
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

	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := c.Models.DailyPrayerTimes.UpdateOne(&result, ws, tx); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return c.APIErr.Database(
				ctx,
				errors.New(v.T.ConflictError()),
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

func (c *ControllerBasic) Destroy(ctx echo.Context) error {
	var result daily_prayer_times.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)

	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	if err := c.Models.DailyPrayerTimes.DeleteOne(&result, ws, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	if err := tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Rolling(ctx echo.Context) error {
	cityIDStr := ctx.QueryParam("city_id")
	if cityIDStr == "" {
		return c.APIErr.BadRequest(ctx, errors.New("city_id is required"))
	}

	cityID, err := uuid.Parse(cityIDStr)
	if err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}

	now := time.Now()
	day := now.Day()
	month := int(now.Month())

	if dateParam := ctx.QueryParam("date"); dateParam != "" {
		t, err := time.Parse("2006-01-02", dateParam)
		if err != nil {
			return c.APIErr.BadRequest(ctx, errors.New("invalid date format, use YYYY-MM-DD"))
		}
		day = t.Day()
		month = int(t.Month())
	}

	times, err := c.Models.DailyPrayerTimes.GetRollingPrayerTimes(ctx.Request().Context(), cityID, day, month)
	if err != nil {
		return c.APIErr.Database(ctx, err, nil)
	}

	return ctx.JSON(http.StatusOK, map[string]any{
		"data": times,
		"meta": map[string]any{
			"from_date": fmt.Sprintf("%d-%02d-%02d", now.Year(), month, day),
			"total":     len(times),
			"note":      "30 days starting from today (inclusive)",
		},
	})
}
