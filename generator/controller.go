package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Function to generate the controller file
func generateControllerFile(modelName string) error {
	// Define the directory and file name
	f := fmt.Sprintf("%s_controller", strings.ToLower(modelName))

	routerFile := filepath.Clean(
		path.Join("./", "controllers", f, "router.go"),
	)
	// Ensure the directory exists
	err := os.MkdirAll(filepath.Clean(
		path.Join("./", "controllers", f),
	), 0o755)
	if err != nil {
		return err
	}

	file, err := os.Create(routerFile)
	if err != nil {
		return err
	}
	defer file.Close()

	// Content template with a placeholder for the model name
	content := fmt.Sprintf(`package %s_controller

import (
	"app/controller"
	"models/%s"

	"bitbucket.org/sadeemTechnology/backend-model"
	"github.com/labstack/echo/v4"
)

type Controllers struct {
	*controller.Dependencies
}

func Get(d *controller.Dependencies) *Controllers {
	return &Controllers{d}
}

// scope ----------------------------------------------------------------------

func (c *Controllers) scope(ctx echo.Context) *%s.WhereScope {
	scopes := c.Utils.CtxScopes(ctx)
	var admin, public bool
	for _, v := range scopes {
		switch v {
		case "admin":
			admin = true
		case "public":
			public = true
		}
	}

	return &%s.WhereScope{
		IsAdmin:  admin,
		IsPublic: public && !admin,
	}
}

func (m *Controllers) SetRoutes(
	d *controller.RouterDependencies,
) {

	b := d.E.Group("/%s")

	b.GET("", m.Index, r).Name = "%s:index:admin,public"
	b.GET("/:id", m.Show, r).Name = "%s:show:admin,public"

	r := d.Requires(model.ScopeAdmin)

	b.POST("", m.Store, r).Name = "%s:store:admin"
	b.PUT("/:id", m.Update, r, a).Name = "%s:update:admin"
	b.DELETE("/:id", m.Destroy, r, a).Name = "%s:destroy:admin"
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
	)

	// Write the content to the file
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	log.Printf("generated: %s\n", routerFile)
	return nil
}
