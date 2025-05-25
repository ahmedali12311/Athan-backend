package config

import (
	"cmp"
	"log"
	"os"
	"path"
	"runtime/debug"
	"strconv"
	"time"
)

type CommitInfo struct {
	FullSHA1 string
	Time     string
}

type Settings struct {
	Timezone             string
	Port                 string
	Domain               string
	RootDIR              string
	CryptoKey            string
	Env                  string
	DSN                  string
	RootPath             string
	UploadsPath          string
	MigrationsRoot       string
	AppCode              string
	AppDesc              string
	AppName              string
	AllowedOrigins       []string
	JwtExpiry            time.Duration
	CRC32Poly            uint32
	GoogleServiceAccount string
	CommitInfo           CommitInfo
}

const (
	MaxFormMemory = "16mb"
)

var (
	JwtSecret            = []byte(os.Getenv("JWT_SECRET_KEY"))
	DOMAIN               = os.Getenv("DOMAIN")
	Timezone             = os.Getenv("TIMEZONE")
	RootDIR              = os.Getenv("ROOT_DIR")
	MigrationsRoot       = os.Getenv("MIGRATIONS_ROOT")
	SeedersRoot          = os.Getenv("SEEDERS_ROOT")
	ENV                  = os.Getenv("ENV")
	DSN                  = os.Getenv("CONNECTION_STRING")
	CryptoKey            = []byte(os.Getenv("CRYPTO_KEY"))
	Port                 = os.Getenv("PORT")
	AppCode              = os.Getenv("APP_CODE")
	AppDesc              = os.Getenv("APP_DESC")
	AppName              = os.Getenv("APP_NAME")
	JwtExpiry            = getJWTExpiryHours()
	GoogleServiceAccount = getServiceAccountPath()
	OTP_URL              = os.Getenv("OTP_URL")
	OTP_JWT              = os.Getenv("OTP_JWT")
	OTP_KEY              = os.Getenv("OTP_KEY")
	OTP_ENV              = os.Getenv("OTP_ENV")

	// CRC32Poly this polynomial ensures the checksum reproduces the same
	// hashes as orion v1 and v2
	CRC32Poly = uint32(0xEDB88320)
)

func TimeNow() time.Time {
	return time.Now().UTC()
}

func GetSettings() *Settings {
	return &Settings{
		Env:       ENV,
		DSN:       DSN,
		Port:      Port,
		Domain:    DOMAIN,
		RootDIR:   RootDIR,
		AppCode:   AppCode,
		AppDesc:   AppDesc,
		AppName:   AppName,
		CRC32Poly: CRC32Poly,
		Timezone:  Timezone,

		MigrationsRoot:       MigrationsRoot,
		GoogleServiceAccount: GoogleServiceAccount,

		CommitInfo:     getCommitInfo(),
		RootPath:       GetRootPath(""),
		UploadsPath:    GetUploadsPath(""),
		JwtExpiry:      getJWTExpiryHours(),
		AllowedOrigins: getAllowedOrigins(ENV),
	}
}

func getAllowedOrigins(env string) []string {
	origins := []string{
		"http://localhost",
	}
	if env == "production" {
		origins = []string{
			"https://production-url.com",
		}
	}
	return origins
}

func GetPrivatePath(dir string) string {
	return path.Join(RootDIR, "private", dir)
}

func GetUploadsPath(dir string) string {
	return path.Join(RootDIR, "public", "uploads", dir)
}

func GetRootPath(dir string) string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	return path.Join(path.Dir(ex), RootDIR, dir)
}

func GetSchemaURL(schemaName string) string {
	return DOMAIN + "/schemas/" + schemaName + ".json"
}

func GetFontPath(font string) string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	return path.Join(path.Dir(ex), "private", "fonts", font)
}

func GetLogoPath() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	return path.Join(path.Dir(ex), "private", "logo.png")
}

func getJWTExpiryHours() time.Duration {
	expEnv := cmp.Or(os.Getenv("JWT_EXPIRY_IN_HOURS"), "720")
	if expInt, err := strconv.ParseFloat(expEnv, 64); err != nil {
		return 1 * time.Hour
	} else {
		return time.Duration(expInt) * time.Hour
	}
}

func getServiceAccountPath() string {
	ex, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	}
	return path.Join(path.Dir(ex), "json_google", "api_key.json")
}

func getCommitInfo() CommitInfo {
	var commitInfo CommitInfo

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commitInfo.FullSHA1 = setting.Value
			}
			if setting.Key == "vcs.time" {
				commitInfo.Time = setting.Value
			}
		}
		return commitInfo
	}
	return commitInfo
}

// SQLSelectURLPath returns a conditionally selected image path column as in
//
//		 CASE WHEN nullif(banners.img, '') is not null
//	     THEN FORMAT('http://localhost:8056/%s', banners.img)
//		 ELSE null
//		 END as img
func SQLSelectURLPath(tableName, colName, aliasName string) string {
	return `
    CASE
        WHEN nullif(` + tableName + `.` + colName + `, '') is not null
        THEN FORMAT('` + DOMAIN + `/%s', ` + tableName + `.` + colName + `)
        ELSE null
    END as ` + aliasName
}
