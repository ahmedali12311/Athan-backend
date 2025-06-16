package meta_controller

import (
	"net/http"
	"strconv"
	"time"

	"app/controller"
	"app/models/setting"
	"app/models/user"
	"bitbucket.org/sadeemTechnology/backend-config"

	"github.com/labstack/echo/v4"
)

type Controllers struct {
	*controller.Dependencies
}

func Get(d *controller.Dependencies) *Controllers {
	return &Controllers{d}
}

func (c *Controllers) SetPublicRoutes(
	d *controller.RouterDependencies,
) {
	d.E.GET("/meta", c.Index).Name = "meta:index:public"
}

func (c *Controllers) Index(ctx echo.Context) error {
	settings := []setting.Model{}
	if err := c.Models.Setting.GetForMeta(&settings); err != nil {
		return c.APIErr.Database(ctx, err, &setting.Model{})
	}
	nowDB, err := c.Models.Setting.GetDBTime()
	if err != nil {
		nowDB = ""
	}
	scopes := []string{}
	if err := c.Models.Permission.DistinctScopes(&scopes); err != nil {
		return c.APIErr.BadRequest(ctx, err)
	}

	settingsMap := make(map[string]string)
	for i := range settings {
		settingsMap[settings[i].Key] = settings[i].Value
	}

	commitCountInt, err := strconv.ParseInt(*c.Utils.CommitCount, 10, 64)
	if err != nil {
		commitCountInt = 0
	}

	meta := map[string]any{
		"domain":      c.Utils.Settings.Domain,
		"api_version": c.Utils.Version,
		"commit": map[string]any{
			"describe": c.Utils.CommitDescribe,
			"count":    commitCountInt,
			"sha1":     c.Utils.Settings.CommitInfo.FullSHA1,
			"time":     c.Utils.Settings.CommitInfo.Time,
		},
		"project_code": c.Utils.Settings.AppCode,
		"project_name": c.Utils.Settings.AppName,
		"enums": map[string]any{
			"locales": []string{
				"ar",
				"en",
			},
			"genders": user.AllGenders,
		},
		"settings": settingsMap,
		// "super_parents": category.SuperParentsMap, // FIX: get value
		"server_time": map[string]any{
			"pg_now": nowDB,
			"go_now": config.TimeNow(),
			"utc_2":  config.TimeNow().Add(time.Hour * 2),
		},
		"filter_op": map[string]string{
			"eq":    "=",
			"nq":    "!=",
			"gt":    ">",
			"gte":   ">=",
			"lt":    "<",
			"lte":   "<=",
			"ex":    "exactly - returns rows where relation has one value",
			":null": "na - returns rows where \"is null\" filter=column:null",
		},
		"scopes": scopes,
	}

	return ctx.JSON(http.StatusOK, meta)
}
