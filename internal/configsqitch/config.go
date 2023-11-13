package configsqitch

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log/slog"
	"text/template"
	"time"

	"github.com/FinalCAD/FinalcadOne-sqitch/internal/utils"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"

	_ "github.com/lib/pq"
)

type ConfigSqitch struct {
	PostgresUser     string `json:"postgresuser"`
	PostgresPassword string `json:"-"`
	PostgresURI      string `json:"postgresuri"`
	PostgresPort     string `json:"postgresport"`
	PostgresDB       string `json:"postgresdb"`
	Region           string `json:"region"`
	Profile          string `json:"profile"`
	Timeout          int    `json:"timeout"`
	ConfigFilepath   string `json:"configfilepath"`
	SecretFilepath   string `json:"secretfilepath"`
}

const (
	DEFAULT_REGION             = "us-east-1"
	DEFAULT_TIMEOUT            = 5000
	DEFAULT_PORT               = "5432"
	DEFAULT_CONFIG_SQITCH_PATH = "sqitch.conf"
	DEFAULT_SECRET_SQITCH_PATH = "secret.conf"
)

//go:embed sqitch.tmpl
var templateConfig string

func (c ConfigSqitch) String() string {
	return fmt.Sprintf("{PostgresUser: %s, PostgresURI: %s, PostgresPort: %s, PostgresDB: %s, Region: %s, Profile: %s, Timeout : %v}",
		c.PostgresUser, c.PostgresURI, c.PostgresPort, c.PostgresDB, c.Region, c.Profile, c.Timeout)
}

func GetConfig() (*ConfigSqitch, error) {
	var err error
	configSqitch := ConfigSqitch{Timeout: DEFAULT_TIMEOUT}

	configSqitch.PostgresURI = utils.Getenv("POSTGRES_URI", "notset")
	configSqitch.PostgresPort = utils.Getenv("POSTGRES_PORT", DEFAULT_PORT)
	configSqitch.PostgresDB = utils.Getenv("POSTGRES_DB", "notset")
	configSqitch.Region = utils.Getenv("AWS_REGION", DEFAULT_REGION)
	configSqitch.Profile = utils.Getenv("AWS_PROFILE", "")
	configSqitch.ConfigFilepath = utils.Getenv("CONFIG_SQITCH_PATH", DEFAULT_CONFIG_SQITCH_PATH)
	configSqitch.SecretFilepath = utils.Getenv("SECRET_SQITCH_PATH", DEFAULT_SECRET_SQITCH_PATH)

	configSqitch.PostgresUser = utils.Getenv("POSTGRES_IAM_USER", "")
	if configSqitch.PostgresUser == "" {
		configSqitch.PostgresUser = utils.Getenv("POSTGRES_USER", "notset")
		configSqitch.PostgresPassword = utils.Getenv("POSTGRES_PASSWORD", "notset")
	} else {
		configSqitch.PostgresPassword, err = configSqitch.Connect()
		if err != nil {
			return nil, err
		}
		slog.Info("Successfully connected to database")
	}
	slog.Debug(fmt.Sprintf("ConfigSqitch struct: %v", configSqitch))

	return &configSqitch, nil
}

func (c *ConfigSqitch) WriteConfig() error {
	outputFile, err := utils.CreateFile(c.ConfigFilepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	tmpl := template.Must(template.New("configTemplate").Parse(templateConfig))
	err = tmpl.Execute(outputFile, &c)
	if err != nil {
		return err
	}

	secretFile, err := utils.CreateFile(c.SecretFilepath)
	if err != nil {
		return err
	}
	defer secretFile.Close()
	_, err = secretFile.WriteString(fmt.Sprintf("SQITCH_PASSWORD=\"%s\"", c.PostgresPassword))

	return err
}

func (c *ConfigSqitch) Connect() (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(c.Timeout)*time.Millisecond)
	defer cancel()
	var err error
	var cfg aws.Config

	if c.Profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(c.Profile), config.WithRegion(c.Region))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(c.Region),
			config.WithRetryer(func() aws.Retryer { return retry.AddWithMaxAttempts(aws.NopRetryer{}, 1) }))
	}

	if err != nil {
		return "", fmt.Errorf("configuration error: %s", err.Error())
	}

	authenticationToken, err := auth.BuildAuthToken(
		context.TODO(), fmt.Sprintf("%s:%s", c.PostgresURI, c.PostgresPort), c.Region, c.PostgresUser, cfg.Credentials)
	if err != nil {
		slog.Debug(fmt.Sprintf("ConfigSqitch struct: %v", c))
		return "", fmt.Errorf("failed to create authentication token: %s", err.Error())
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		c.PostgresURI, c.PostgresPort, c.PostgresUser, authenticationToken, c.PostgresDB,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return "", err
	}

	err = db.Ping()
	if err != nil {
		return "", err
	}

	return authenticationToken, nil
}
