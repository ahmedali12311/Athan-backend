package apierrors

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	finder "bitbucket.org/sadeemTechnology/backend-finder"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
)

type PGErrContent struct {
	Code       string `json:"code"`
	Table      string `json:"table"`
	Column     string `json:"column"`
	Constraint string `json:"constraint"`
	Details    string `json:"details"`
	Message    string `json:"message"`
}

// PGErrMessage is shared because it exists in translations
type PGErrMessage struct {
	DB PGErrContent `json:"db"`
}

func (e *APIErrors) Database(
	ctx echo.Context,
	err error,
	m finder.Model,
) error {
	t := e.Utils.CtxT(ctx)
	dbe := PGErrMessage{
		DB: PGErrContent{
			Message: err.Error(),
		},
	}
	message := "database error."
	status := http.StatusBadRequest
	errType := Conflict
	if errors.Is(err, sql.ErrNoRows) {
		status = http.StatusNotFound
		errType = NotFound
		message = t.ModelNotFound(t.ModelName(m.ModelName()))
	}
	if strings.Contains(err.Error(), "SQLSTATE") {
		status = http.StatusConflict
		pgErr, ok := err.(*pgconn.PgError) //nolint: errorlint // must not wrap
		if ok {
			dbe.DB.Code = pgErr.Code
			dbe.DB.Details = pgErr.Detail
			dbe.DB.Table = pgErr.TableName
			dbe.DB.Message = pgErr.Message
			dbe.DB.Column = pgErr.ColumnName
			dbe.DB.Constraint = pgErr.ConstraintName

			message = t.PGError(
				ctx.Request().Method,
				pgErr.Code,
				pgErr.ConstraintName,
			) + " " + pgErr.Message
		}
	}
	return e.respond(ctx, status, errType, message, dbe)
}
