package fcm_notification_controller

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	"app/controller"
	"app/pkg/firebase_utils"

	fcm "app/models/fcm_notification"
	"app/models/user"
	"app/models/user_notification"

	"github.com/labstack/echo/v4"
)

type ControllerBasic struct {
	*controller.Dependencies
}

func (c *ControllerBasic) Index(ctx echo.Context) error {
	indexResponse, err := c.Models.FcmNotification.GetAll(ctx, nil)
	if err != nil {
		return c.APIErr.Database(ctx, err, nil)
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *ControllerBasic) Show(ctx echo.Context) error {
	var result fcm.Model

	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	if err := c.Models.FcmNotification.GetOne(&result, nil); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Store(ctx echo.Context) error {
	var result fcm.Model
	ctxUser := c.Utils.CtxUser(ctx)
	result.SenderID = &ctxUser.ID

	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	if valid := result.MergeAndValidate(v); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}

	message := firebase_utils.BuildTopicMessage(
		&result.Title,
		&result.Body,
		result.Topic,
		result.Data,
	)
	response, err := c.Utils.FBM.Send(context.Background(), message)
	if err != nil {
		return err
	}
	result.Response = &response

	if err := c.Models.FcmNotification.CreateOne(&result); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	userNotification := user_notification.Model{
		UserID:     nil,
		IsRead:     true,
		IsNotified: true,
		Title:      result.Title,
		Body:       result.Body,
		Response:   result.Response,
		Data:       result.Data,
	}

	if err := c.Models.UserNotification.Create(&userNotification); err != nil {
		return c.APIErr.Database(ctx, err, &userNotification)
	}
	return ctx.JSON(http.StatusCreated, result)
}

func (c *ControllerBasic) Update(ctx echo.Context) error {
	var result fcm.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	if err := c.Models.FcmNotification.GetOne(&result, nil); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}
	if valid := result.MergeAndValidate(v); !valid {
		return c.APIErr.InputValidation(ctx, v)
	}

	if err := c.Models.FcmNotification.UpdateOne(
		&result,
		nil,
		nil,
	); err != nil {
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
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Destroy(ctx echo.Context) error {
	var result fcm.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	if err := c.Models.FcmNotification.GetOne(&result, nil); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	if err := c.Models.FcmNotification.DeleteOne(
		&result,
		nil,
		nil,
	); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) NotifyUser(ctx echo.Context) error {
	var result fcm.Model
	// Initialize a new Validator instance.
	t := c.Utils.CtxT(ctx)
	ctxUser := c.Utils.CtxUser(ctx)
	result.SenderID = &ctxUser.ID
	v, err := c.GetValidator(ctx, result.ModelName())
	if err != nil {
		return err
	}

	userID := v.Data.GetUUID("user_id")
	v.Exists(userID, "user_id", "id", "users", true)
	v.AssignString("title", &result.Title, 0, 500)
	v.AssignString("body", &result.Body, 0, 500)
	result.Topic = nil

	if !v.Valid() {
		return c.APIErr.InputValidation(ctx, v)
	}
	token, err := c.Models.User.GetFCMToken(
		userID,
		fcm.TokenTypeStandard,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || err.Error() == "sql: no rows in result set" {
			return c.APIErr.Database(ctx, errors.New(t.UserTokenNotFound()), &result)
		}

		return c.APIErr.Database(
			ctx,
			err,
			&result,
		)
	}
	result.Data = map[string]string{
		"user_id": userID.String(),
	}
	message := firebase_utils.BuildTokenMessage(
		&result.Title,
		&result.Body,
		token,
		result.Data,
	)
	response, err := c.Utils.FBM.Send(context.Background(), message)
	if err != nil {
		return c.APIErr.Firebase(ctx, err)
	}
	result.Response = &response
	topic := "notify-user"
	result.Topic = &topic
	if err := c.Models.FcmNotification.CreateOne(&result); err != nil {
		return c.APIErr.Database(
			ctx,
			err,
			&result,
		)
	}

	userNotification := user_notification.Model{
		User: user.MinimalModel{
			ID: userID,
		},
		IsRead:     false,
		IsNotified: false,
		Title:      result.Title,
		Body:       result.Body,
		Response:   result.Response,
		Data:       result.Data,
	}

	if err := c.Models.UserNotification.Create(&userNotification); err != nil {
		return c.APIErr.Database(ctx, err, &userNotification)
	}
	return ctx.JSON(http.StatusCreated, result)
}
