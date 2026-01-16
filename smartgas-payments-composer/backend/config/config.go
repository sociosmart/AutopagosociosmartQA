package config

import (
	"strconv"

	"github.com/golobby/config/v3"
	"github.com/golobby/config/v3/pkg/feeder"
)

var ConfigSettings *Config

type DB struct {
	Name string `env:"DB_NAME"`
	Port int    `env:"DB_PORT"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASS"`
	Host string `env:"DB_HOST"`
}

type SMTP struct {
	User     string `env:"SMTP_USER"`
	Password string `env:"SMTP_PASSWORD"`
	Host     string `env:"SMTP_HOST"`
	Port     uint   `env:"SMTP_PORT"`
}

type Invoicing struct {
	Rfc               string `env:"INVOICING_RFC"`
	CpEmitter         string `env:"INVOICING_CP_EMITTER"`
	Name              string `env:"INVOICING_NAME"`
	FiscalName        string `env:"INVOICING_FISCAL_NAME"`
	CertificateNumber string `env:"INVOICING_CERTIFICATE_NUMBER"`
}

type Config struct {
	Host                string `env:"HOST"`
	Port                int    `env:"PORT"`
	Debug               bool   `env:"DEBUG"`
	SecretKey           string `env:"SECRET_KEY"`
	SecretKeyRefresh    string `env:"SECRET_KEY_REFRESH"`
	JwtExpMinutes       uint   `env:"JWT_EXP_MINUTES"`
	JwtRefreshExpDays   uint   `env:"JWT_REFRESH_EXP_DAYS"`
	Tz                  string `env:"TZ"`
	TrustedProxies      string `env:"TRUSTED_PROXIES"`
	AllowedHosts        string `env:"ALLOWED_HOSTS"`
	StripeSecretKey     string `env:"STRIPE_SECRET_KEY"`
	SocioSmartUrl       string `env:"SOCIO_SMART_URL"`
	StripeWebhookSecret string `env:"STRIPE_WEBHOOK_SECRET"`
	SwitBaseUrl         string `env:"SWIT_BASE_URL"`
	SwitBusiness        string `env:"SWIT_BUSINESS"`
	SwitToken           string `env:"SWIT_TOKEN"`
	SwitApiKey          string `env:"SWIT_API_KEY"`

	DebitBaseUrl string `env:"DEBIT_BASE_URL"`
	DebitAppKey  string `env:"DEBIT_APP_KEY"`
	DebitApiKey  string `env:"DEBIT_API_KEY"`

	ConectiaUrl    string `env:"CONECTIA_URL"`
	ConectiaUrlApi string `env:"CONECTIA_URL_API"`
	ConectiaToken  string `env:"CONECTIA_TOKEN"`

	FromEmail string `env:"FROM_EMAIL"`

	SentryDsn   string `env:"SENTRY_DSN"`
	Environment string `env:"ENVIRONMENT"`

	DB DB

	SMTP SMTP

	Invoicing Invoicing
}

func NewConfig() (c Config, err error) {
	envFeeder := feeder.Env{}
	err = config.New().AddFeeder(envFeeder).AddStruct(&c).Feed()

	if err == nil {
		ConfigSettings = &c
	}

	return
}

func (cfg *Config) Setup() error {
	if cfg.Tz == "" {
		cfg.Tz = "UTC"
	}

	if cfg.Port == 0 {
		cfg.Port = 8008
	}

	if cfg.DB.Port == 0 {
		cfg.DB.Port = 3306
	}

	if cfg.JwtExpMinutes == 0 {
		cfg.JwtExpMinutes = 30
	}

	if cfg.JwtRefreshExpDays == 0 {
		cfg.JwtRefreshExpDays = 30
	}

	// avoiding empty
	if cfg.SecretKey == "" {
		cfg.SecretKey = "secret_key"
	}

	if cfg.SecretKeyRefresh == "" {
		cfg.SecretKeyRefresh = "secret_refresh_key"
	}

	if cfg.TrustedProxies == "" {
		cfg.TrustedProxies = "*"
	}

	if cfg.Host == "" {
		cfg.Host = "localhost:" + strconv.Itoa(cfg.Port)
	}

	if cfg.AllowedHosts == "" {
		cfg.AllowedHosts = "*"
	}

	if cfg.Environment == "" {
		cfg.Environment = "development"
	}

	return nil
}
