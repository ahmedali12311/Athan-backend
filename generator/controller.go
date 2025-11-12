package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func generateControllerFiles(modelName string) error {
	// Generate router.go
	if err := generateRouterFile(modelName); err != nil {
		return err
	}

	// Generate controller_basic.go
	if err := generateControllerBasicFile(modelName); err != nil {
		return err
	}
	// Generate controllers.go
	if err := generateControllersFile(modelName); err != nil {
		return err
	}
	return nil
}

func generateRouterFile(modelName string) error {
	controllerDir := fmt.Sprintf("%s_controller", strings.ToLower(modelName))
	routerFile := filepath.Clean(
		path.Join("./", "controllers", controllerDir, "router.go"),
	)

	err := os.MkdirAll(filepath.Dir(routerFile), 0o755)
	if err != nil {
		return err
	}

	file, err := os.Create(routerFile)
	if err != nil {
		return err
	}
	defer file.Close()

	content := fmt.Sprintf(`package %s_controller

import (
	"app/controller"
	"app/models/%s"

	"bitbucket.org/sadeemTechnology/backend-model"
	"github.com/labstack/echo/v4"
)


func (m *Controllers) SetRoutes(
	d *controller.RouterDependencies,
) {
	b := d.E.Group("/%s")

	b.GET("", m.Basic.Index).Name = "%s:index:admin,public"
	b.GET("/:id", m.Basic.Show).Name = "%s:show:admin,public"

	r := d.Requires(model.ScopeAdmin)

	b.POST("", m.Basic.Store, r).Name = "%s:store:admin"
	b.PUT("/:id", m.Basic.Update, r).Name = "%s:update:admin"
	b.DELETE("/:id", m.Basic.Destroy, r).Name = "%s:destroy:admin"
}
`,
		strings.ToLower(modelName),
		modelName,
		strings.ToLower(modelName),
		strings.ToLower(modelName),
		strings.ToLower(modelName),
		strings.ToLower(modelName),
		strings.ToLower(modelName),
		strings.ToLower(modelName),
	)

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	log.Printf("generated: %s\n", routerFile)
	return nil
}

func generateControllerBasicFile(modelName string) error {
	controllerDir := fmt.Sprintf("%s_controller", strings.ToLower(modelName))
	basicFile := filepath.Clean(
		path.Join("./", "controllers", controllerDir, "controller_basic.go"),
	)

	file, err := os.Create(basicFile)
	if err != nil {
		return err
	}
	defer file.Close()

	content := fmt.Sprintf(`package %s_controller

import (
	"database/sql"
	"errors"
	"net/http"
	"net/url"

	"app/controller"
	"app/models/%s"

	"github.com/labstack/echo/v4"
)

type ControllerBasic struct {
	*controller.Dependencies
}

// Scopes ---------------------------------------------------------------------

func (c *ControllerBasic) scope(
	ctx echo.Context,
) *%s.WhereScope {
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
	
	ws := &%s.WhereScope{
		IsAdmin:  admin,
		IsPublic: public && !admin,
		QueryParams: ctx.QueryParams(),
	}
	
	if ctxUser != nil {
		ws.UserID = &ctxUser.ID
	}
	
	return ws
}

// Actions --------------------------------------------------------------------

func (c *ControllerBasic) Index(ctx echo.Context) error {
	ws := c.scope(ctx)
	indexResponse, err := c.Models.%s.GetAll(ctx, ws)
	if err != nil {
		return c.APIErr.Database(ctx, err, nil)
	}
	return ctx.JSON(http.StatusOK, indexResponse)
}

func (c *ControllerBasic) Show(ctx echo.Context) error {
	var result %s.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)
	if err := c.Models.%s.GetOne(&result, ws); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}
	return ctx.JSON(http.StatusOK, result)
}

func (c *ControllerBasic) Store(ctx echo.Context) error {
	ws := c.scope(ctx)
	var result %s.Model
	
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

	if err := c.Models.%s.CreateOne(&result, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	if err = tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusCreated, result)
}

func (c *ControllerBasic) Update(ctx echo.Context) error {
	var result %s.Model
	if err := c.Utils.ReadUUIDParam(&result.ID, ctx); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}
	ws := c.scope(ctx)
	if err := c.Models.%s.GetOne(&result, ws); err != nil {
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

	if err := c.Models.%s.UpdateOne(&result, ws, tx); err != nil {
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
	var result %s.Model
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

	if err := c.Models.%s.DeleteOne(&result, ws, tx); err != nil {
		return c.APIErr.Database(ctx, err, &result)
	}

	if err := tx.Commit(); err != nil {
		return c.APIErr.InternalServer(ctx, err)
	}
	return ctx.JSON(http.StatusOK, result)
}
`,
		strings.ToLower(modelName),
		modelName,
		modelName,
		modelName,
		modelName,
		modelName,
		modelName,
		modelName,
		modelName,
		modelName,
		modelName,
		modelName,
		modelName,
		modelName,
	)

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	log.Printf("generated: %s\n", basicFile)
	return nil
}

func generateControllersFile(modelName string) error {
	controllerDir := fmt.Sprintf("%s_controller", strings.ToLower(modelName))
	routerFile := filepath.Clean(
		path.Join("./", "controllers", controllerDir, "router.go"),
	)

	err := os.MkdirAll(filepath.Dir(routerFile), 0o755)
	if err != nil {
		return err
	}

	file, err := os.Create(routerFile)
	if err != nil {
		return err
	}
	defer file.Close()

	content := fmt.Sprintf(`package %s_controller

import (
	"app/controller"
	"app/models/%s"

	"bitbucket.org/sadeemTechnology/backend-model"
	"github.com/labstack/echo/v4"
)

type Controllers struct {
	*controller.Dependencies
}

func Get(d *controller.Dependencies) *Controllers {
	return &Controllers{d}
}

`,
		modelName,
		modelName,
	)

	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	log.Printf("generated: %s\n", routerFile)
	return nil
}
