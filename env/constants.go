package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"runtime"
)

var (
	JwtSecret       = "MESSAGEWITH_JWT_SECRET"
	MockupIpAddress = "MESSAGEWITH_MOCKUP_IP_ADDRESS"
	DatabaseURI     = "MESSAGEWITH_DATABASE_URI"
	Domain          = "MESSAGEWITH_DOMAIN"
	SmtpHost        = "MESSAGEWITH_SMTP_HOST"
	SmtpPort        = "MESSAGEWITH_SMTP_PORT"
	SmtpUsername    = "MESSAGEWITH_SMTP_USERNAME"
	SmtpPassword    = "MESSAGEWITH_SMTP_PASSWORD"
	SmtpEmail       = "MESSAGEWITH_SMTP_EMAIL"
	RootDir         = "MESSAGEWITH_ROOTDIR"
)

func InitEnvConstants() {
	_, b, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(b), "..")
	err := godotenv.Load(fmt.Sprintf("%v/.env", root))
	if err != nil {
		panic(fmt.Errorf("create .env file"))
	}

	JwtSecret = os.Getenv(JwtSecret)
	MockupIpAddress = os.Getenv(MockupIpAddress)
	DatabaseURI = os.Getenv(DatabaseURI)
	Domain = os.Getenv(Domain)
	SmtpHost = os.Getenv(SmtpHost)
	SmtpPort = os.Getenv(SmtpPort)
	SmtpUsername = os.Getenv(SmtpUsername)
	SmtpPassword = os.Getenv(SmtpPassword)
	SmtpEmail = os.Getenv(SmtpEmail)
	RootDir = os.Getenv(RootDir)
}
