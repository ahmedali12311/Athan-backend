package api

import (
	"expvar"
	"net/http"
	"runtime"
	"strconv"

	"app/controller"

	config "bitbucket.org/sadeemTechnology/backend-config"

	"github.com/labstack/echo/v4"
)

func (app *Application) Routes(e *echo.Echo, isTest bool) http.Handler {
	// Application Variables --------------------------------------------------
	if !isTest {
		commitCountInt, err := strconv.ParseInt(CommitCount, 10, 64)
		if err != nil {
			commitCountInt = 0
		}
		expvar.NewString("version").Set(Version)
		expvar.NewString("environment").Set(app.Config.Env)
		expvar.NewString("app_code").Set(app.Config.AppCode)
		expvar.NewString("app_name").Set(app.Config.AppName)
		expvar.NewString("commit_sha1").Set(app.Config.CommitInfo.FullSHA1)
		expvar.NewString("commit_time").Set(app.Config.CommitInfo.Time)
		expvar.NewString("commit_describe").Set(CommitDescribe)
		expvar.NewInt("commit_count").Set(commitCountInt)
		// Publish the number of active goroutines.
		expvar.Publish("goroutines", expvar.Func(func() any {
			return runtime.NumGoroutine()
		}))
		// Publish the database connection pool statistics.
		expvar.Publish("database", expvar.Func(func() any {
			return app.DB.Stats()
		}))
		// Publish the current Unix timestamp.
		expvar.Publish("timestamp", expvar.Func(func() any {
			return config.TimeNow().Unix()
		}))
	}

	// Middlewares setup pre web / api scaffolding ----------------------------
	app.SetupMiddlewares(e, isTest)

	// API / V1 / standard setup ----------------------------------------------
	v1 := e.Group("/api/v1")
	deps := &controller.RouterDependencies{
		E:        v1,
		Requires: app.requires,
	}

	// API / V1 ---------------------------------------------------------------
	app.Controllers.Meta.SetPublicRoutes(deps)

	app.Controllers.User.SetOTPRoutes(deps)
	app.Controllers.User.SetAuthRoutes(deps)
	app.Controllers.User.SetAdminRoutes(deps)
	app.Controllers.User.SetProfileRoutes(deps)

	app.Controllers.Category.SetRoutes(deps)

	app.Controllers.Role.SetBasicRoutes(deps)
	app.Controllers.Token.SetBasicRoutes(deps)
	app.Controllers.Setting.SetBasicRoutes(deps)
	app.Controllers.Permission.SetBasicRoutes(deps)
	app.Controllers.Wallet.SetBasicRoutes(deps)
	app.Controllers.DailyPrayerTimes.SetBasicRoutes(deps)
	app.Controllers.Hadiths.SetRoutes(deps)
	app.Controllers.SpecialTopics.SetRoutes(deps)
	app.Controllers.Adhkars.SetRoutes(deps)
	app.Controllers.City.SetBasicRoutes(deps)

	// Generate Routes --------------------------------------------------------
	app.routesGen(e.Routes())
	return e
}
