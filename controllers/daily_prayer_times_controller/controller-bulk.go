package daily_prayer_times_controller

import (
	"app/models/daily_prayer_times"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (c *ControllerBasic) BulkCSV(ctx echo.Context) error {
	var result daily_prayer_times.Model

	tx, err := c.Models.DB.Beginx()
	if err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	defer func() { _ = tx.Rollback() }()

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}

	prayerTimes, entryLog, parseErr := c.Models.DailyPrayerTimes.ParseCSV(v, tx)
	if parseErr != nil {
		if parseErr == daily_prayer_times.ErrValidation {
			return c.APIErr.InputValidation(ctx, v)
		}
		return c.APIErr.Database(ctx, parseErr, &result)
	}

	const chunkSize = 1000
	total := len(*prayerTimes)

	for i := 0; i < total; i += chunkSize {
		end := i + chunkSize
		if end > total {
			end = total
		}
		chunk := (*prayerTimes)[i:end]

		_, err := c.Models.DailyPrayerTimes.BulkCreate(&chunk, v, tx)
		if err != nil {
			return c.APIErr.Database(ctx, err, &result)
		}
		c.Utils.Logger.Info().
			Msgf("completed chunk: [%d / %d] ", i/chunkSize+1, total/chunkSize)
	}

	response := map[string]any{
		"commit": "pending",
		"errors": v.GetErrorMap(),
		"log":    entryLog,
	}
	if err := tx.Commit(); err != nil {
		response["commit"] = "failure"
		return ctx.JSON(http.StatusTeapot, response)
	}

	response["commit"] = "success"
	return ctx.JSON(http.StatusOK, response)
}
