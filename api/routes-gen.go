package api

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"app/models/permission"

	config "bitbucket.org/sadeemTechnology/backend-config"
	"github.com/goccy/go-json"
	"github.com/labstack/echo/v4"
)

// routesGen generates permissions for all routes, initializes superadmin, and
// creates routes.json in uploads path
func (app *Application) routesGen(routes []*echo.Route) {
	// TODO:
	// 1. filter sensitve routes
	// 2. serve routes permission with (api-key) protection

	if _, err := app.DB.ExecContext(
		context.Background(),
		`
            WITH sequenced AS (
                SELECT
                    id,
                    ROW_NUMBER() OVER (ORDER BY id) AS new_sequence_id
                FROM
                    permissions
            )

            UPDATE permissions
            SET id = sequenced.new_sequence_id
            FROM sequenced
            WHERE permissions.id = sequenced.id;
        `,
	); err != nil {
		app.Logger.Fatal().Msgf("permissions sequence: %s", err.Error())
	}

	if _, err := app.DB.ExecContext(
		context.Background(),
		`
            SELECT setval(
                'permissions_id_seq',
                (SELECT MAX(id) FROM permissions)
            )
        `,
	); err != nil {
		app.Logger.Fatal().Msgf("permissions sequence: %s", err.Error())
	}

	perms := &[]permission.Model{}
	for _, v := range routes {
		nameSplit := strings.Split(v.Name, ":")

		if len(nameSplit) != 3 {
			err := fmt.Errorf(
				"route: %s is not configured correctly, "+
					"must be in format: table:action:scope1,scope2",
				v.Name,
			)
			app.Logger.Fatal().Msg(err.Error())
		} else {
			ignoredScopes := []string{"file", "dir", "web"}
			x := slices.Contains(ignoredScopes, nameSplit[0])
			if !x {
				scopeSplit := strings.Split(nameSplit[2], ",")
				for _, vv := range scopeSplit {
					*perms = append(*perms, permission.Model{
						Method: v.Method,
						Path:   v.Path,
						Model:  nameSplit[0],
						Action: nameSplit[1],
						Scope:  vv,
					})
				}
			}
		}
	}
	if len(*perms) > 0 {
		deleted, err := app.Models.Permission.DeleteUnused(perms)
		if err != nil {
			app.Logger.Fatal().Msg(err.Error())
		}
		if deleted > 0 {
			app.Logger.Info().Msgf("deleted permissions: %d", deleted)
		}
		affected, err := app.Models.Permission.BulkCreate(perms)
		if err != nil {
			app.Logger.Fatal().Msg(err.Error())
		}
		if affected > 0 {
			app.Logger.Info().Msgf("new permissions: %d", affected)
		}
		superAdminID, err := app.Models.Role.CreateSuperAdmin()
		if err != nil {
			app.Logger.Fatal().Msg(err.Error())
		}
		app.Logger.Info().Msgf("superadmin role id: %d", superAdminID)
		rolePerms, err := app.Models.Role.GrantAllPermissions(superAdminID)
		if err != nil {
			app.Logger.Fatal().Msg(err.Error())
		}
		if rolePerms > 0 {
			app.Logger.Info().Msgf(
				"superadmin was granted %d new permissions",
				rolePerms,
			)
		}

		if err := app.Models.User.CreateSuperAdmin(superAdminID); err != nil {
			app.Logger.Fatal().Msg(err.Error())
		}

		basicRoleID, err := app.Models.Role.CreateBasic()
		if err != nil {
			app.Logger.Fatal().Msg(err.Error())
		}

		basicRoleOwn, err := app.Models.Role.GrantByScope(basicRoleID, "own")
		if err != nil {
			app.Logger.Fatal().Msg(err.Error())
		}
		if basicRoleOwn > 0 {
			app.Logger.
				Info().
				Msgf("basic was granted %d new own permissions", basicRoleOwn)
		}

		basicRolePublic, err := app.Models.Role.GrantByScope(
			basicRoleID,
			"public",
		)

		if err != nil {
			app.Logger.Fatal().Msg(err.Error())
		}
		if basicRolePublic > 0 {
			app.Logger.Info().Msgf(
				"basic was granted %d new public permissions",
				basicRolePublic,
			)
		}

	}
	// customerRoleID, err := app.Models.Role.CreateCustomer()
	// if err != nil {
	// 	app.Logger.Fatal().Msg(err.Error())
	// }
	// customerRole, err := app.Models.Role.GrantByScope(customerRoleID, "customer")
	// if err != nil {
	// 	app.Logger.Fatal().Msg(err.Error())
	// }
	// if customerRole > 0 {
	// 	app.Logger.
	// 		Info().
	// 		Msgf("customer was granted %d new permissions", customerRole)
	// }
	// if err != nil {
	// 	app.Logger.Fatal().Msg(err.Error())
	// }

	adminID, err := app.Models.Role.CreateAdmin()
	if err != nil {
		app.Logger.Fatal().Msg(err.Error())
	}
	app.Logger.Info().Msgf("superadmin role id: %d", adminID)
	adminRoleScope, err := app.Models.Role.GrantByScope(adminID, "admin")
	if err != nil {
		app.Logger.Fatal().Msg(err.Error())
	}
	if adminRoleScope > 0 {
		app.Logger.
			Info().
			Msgf("admin was granted %d new permissions", adminID)
	}

	if err := app.DB.SelectContext(
		context.Background(),
		&app.Permissions,
		`
            SELECT *
            FROM permissions
            ORDER BY "model" ASC,"action" ASC,"scope" ASC
        `,
	); err != nil {
		app.Logger.Fatal().Msgf("permissions preloading: %s", err.Error())
	}
	app.Logger.Info().Msgf(
		"preloaded %d permissions successfully",
		len(app.Permissions),
	)

	// NOTE: shows created permissions on running the binary

	// for _, v := range perms {
	// 	app.Logger.Info().Msgf(
	// 		"%6s:%-30s | %-18s | %-10s | %s",
	// 		v.Method,
	// 		v.Path,
	// 		v.Model,
	// 		v.Action,
	// 		v.Scope,
	// 	)
	// }

	data, err := json.MarshalIndent(routes, "", "  ")
	if err != nil {
		app.Logger.Fatal().Msg(err.Error())
	}
	path := config.GetUploadsPath("routes.json")
	if err := os.WriteFile(path, data, 0o644); err != nil {
		app.Logger.Fatal().Msg(err.Error())
	}
}
