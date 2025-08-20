package privacy_policy

import (
	_ "embed"
	"fmt"
	"html/template"
	"os"
	"strings"

	config "bitbucket.org/sadeemTechnology/backend-config"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

//go:embed privacy-policy.html
var privacyPolicyHTML string

var (
	tmpl = template.Must(
		template.New("privacy-policy.html").Parse(privacyPolicyHTML),
	)
	caser = cases.Title(language.English)
)

type privacyData struct {
	AppName  string
	AppOwner string
}

// GenerateHTML builds a privacy policy using app name from config
// and passed destination and owner name
func GenerateHTML(dest string) error {
	f, err := os.Create(config.GetRootPath(dest))
	if err != nil {
		return fmt.Errorf("error creating html file: %w", err)
	}
	defer f.Close()

	appName := strings.ReplaceAll(config.AppName, "-", " ")
	appName = strings.Trim(appName, " ")
	appName = caser.String(appName)

	appOwner := caser.String(config.AppOwner)

	data := privacyData{
		AppName:  appName,
		AppOwner: appOwner,
	}
	if err := tmpl.Execute(f, data); err != nil {
		return fmt.Errorf(
			"error executing privacy policy html template: %w",
			err,
		)
	}
	return nil
}
