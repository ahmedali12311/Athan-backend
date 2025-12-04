package schemas

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"strings"

	config "bitbucket.org/sadeemTechnology/backend-config"
	category "bitbucket.org/sadeemTechnology/backend-model-category"
	setting "bitbucket.org/sadeemTechnology/backend-model-setting"

	"github.com/rs/zerolog"
	js "github.com/santhosh-tekuri/jsonschema/v5"
)

var (

	// properties
	//go:embed properties/gender.json
	genderSF []byte
	//go:embed properties/id_obj.json
	idObjSF []byte
	//go:embed properties/uuid_obj_null.json
	uuidObjNullSF []byte
	//go:embed properties/uuid_obj.json
	uuidObjSF []byte
	//go:embed properties/line_string.json
	lineStringSF []byte
	//go:embed properties/point.json
	pointSF []byte
	//go:embed properties/point_null.json
	pointNullSF []byte

	// models
	//go:embed user.json
	userSF []byte
	//go:embed role.json
	roleSF []byte
	//go:embed permission.json
	permissionSF []byte
	//go:embed wallet_transaction.json
	WalletTrxSF []byte
	//go:embed fcm_notification.json
	fcmNotificationSF []byte
	//go:embed user_notification.json
	userNotificationSF []byte
	//go:embed city.json
	citySF []byte
	//go:embed adhkar.json
	adhkarSF []byte
	//go:embed daily_prayer_time.json
	dailyPrayerTimesSF []byte
	//go:embed hadith.json
	hadithsSF []byte
	//go:embed special_topic.json
	specialTopicsSF []byte
)

type schemaRegisterar struct {
	Compiler *js.Compiler
	Logger   *zerolog.Logger
}

func (s *schemaRegisterar) register(url, schema string) *js.Schema {
	if err := s.Compiler.AddResource(
		url,
		strings.NewReader(schema),
	); err != nil {
		msg := fmt.Sprintf(
			"json-schema: error compiling %s: %s",
			url,
			err.Error(),
		)
		s.Logger.Fatal().Msg(msg)
	}
	locationSchema, err := s.Compiler.Compile(url)
	if err != nil {
		msg := fmt.Sprintf(
			"json-schema: error compiling %s: %s",
			url,
			err.Error(),
		)
		s.Logger.Fatal().Msg(msg)
	}
	return locationSchema
}

func BuildSchemas(logger *zerolog.Logger) map[string]*js.Schema {
	compiler := js.NewCompiler()
	compiler.Draft = js.Draft2020
	compiler.AssertFormat = true
	compiler.AssertContent = true

	schema := schemaRegisterar{
		Compiler: compiler,
		Logger:   logger,
	}
	// SF stands for Shadow Fiend, also properties need to be generated first
	propertyFiles := map[string][]byte{
		// properties
		"properties/gender.json": genderSF,
		// id/uuid objects
		"properties/id_obj.json":        idObjSF,
		"properties/uuid_obj.json":      uuidObjSF,
		"properties/uuid_obj_null.json": uuidObjNullSF,
		// GIS stuff
		"properties/point.json":       pointSF,
		"properties/point_null.json":  pointNullSF,
		"properties/line_string.json": lineStringSF,
	}
	modelFiles := map[string][]byte{
		"category.json": category.Schema,
		"setting.json":  setting.Schema,

		"permission.json":         permissionSF,
		"role.json":               roleSF,
		"user.json":               userSF,
		"wallet_transaction.json": WalletTrxSF,
		"fcm_notification.json":   fcmNotificationSF,
		"user_notification.json":  userNotificationSF,
		"city.json":               citySF,
		"adhkar.json":             adhkarSF,
		"daily_prayer_time.json":  dailyPrayerTimesSF,
		"hadith.json":             hadithsSF,
		"special_topic.json":      specialTopicsSF,
	}
	domain := config.DOMAIN + "/schemas"
	SchemaMap := map[string]*js.Schema{}

	hrefProps := []string{}
	hrefModels := []string{}
	for k, v := range propertyFiles {
		replaced := strings.ReplaceAll(string(v), "%DOMAIN%", domain)
		createSchemaFile(logger, k, replaced)
		loc := domain + "/" + k
		SchemaMap[loc] = schema.register(loc, replaced)
		hrefProps = append(hrefProps, k)
	}

	for k, v := range modelFiles {
		replaced := strings.ReplaceAll(string(v), "%DOMAIN%", domain)
		loc := domain + "/" + k
		createSchemaFile(logger, k, replaced)
		SchemaMap[loc] = schema.register(loc, replaced)
		hrefModels = append(hrefModels, k)
	}

	if err := generateHTMLPage(hrefProps, hrefModels); err != nil {
		logger.Fatal().Err(err).Msg("error parsing schema html template")
	}
	return SchemaMap
}

func createSchemaFile(logger *zerolog.Logger, filename, data string) {
	filePath := config.GetRootPath("public/schemas/" + filename)
	f, err := os.Create(filePath) //nolint: gosec //shoori ya t3ama
	if err != nil {
		msg := fmt.Sprintf(
			"json-schema: error generating %s: %s",
			filePath,
			err.Error(),
		)
		logger.Fatal().Msg(msg)
		return
	}
	if _, err := f.WriteString(data); err != nil {
		logger.Fatal().Msg(err.Error())
		f.Close()
		return
	}
	if err := f.Close(); err != nil {
		logger.Fatal().Msg(err.Error())
		return
	}
}

func generateHTMLPage(props, models []string) error {
	tmplPath := config.GetPrivatePath("schemas.html")
	filePath := config.GetRootPath("public/schemas/index.html")
	t, err := template.
		New("schemas.html").
		ParseFiles(tmplPath)
	if err != nil {
		return fmt.Errorf("error parsing schema html template: %w", err)
	}
	f, err := os.Create(filePath) //nolint: gosec //shoori ya t3ama
	if err != nil {
		return fmt.Errorf("error creating html file: %w", err)
	}
	x := struct {
		Props  []string
		Models []string
	}{
		Props:  props,
		Models: models,
	}
	if err := t.Execute(f, x); err != nil {
		return fmt.Errorf("error executing schema html template: %w", err)
	}
	if err := f.Close(); err != nil {
		return fmt.Errorf("error closing schema file: %w", err)
	}
	return nil
}
