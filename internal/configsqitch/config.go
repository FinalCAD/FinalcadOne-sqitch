package configsqitch

import (
	_ "embed"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/FinalCAD/FinalcadOne-sqitch/internal/utils"
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
	Filepath         string `json:"filepath"`
}

const (
	DEFAULT_REGION             = "us-east-1"
	DEFAULT_TIMEOUT            = 5000
	DEFAULT_PORT               = "5432"
	DEFAULT_CONFIG_SQITCH_PATH = "sqitch.conf"
)

//go:embed sqitch.tmpl
var templateConfig string

func (c ConfigSqitch) String() string {
	return fmt.Sprintf("{PostgresUser: %s, PostgresURI: %s, PostgresPort: %s, PostgresDB: %s, Region: %s, Profile: %s, Timeout : %v}",
		c.PostgresUser, c.PostgresURI, c.PostgresPort, c.PostgresDB, c.Region, c.Profile, c.Timeout)
}

func WriteConfig(configSqitch *ConfigSqitch) error {
	outputFile, err := utils.CreateFile(configSqitch.Filepath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	tmpl := template.Must(template.New("configTemplate").Parse(templateConfig))
	return tmpl.Execute(outputFile, configSqitch)
}

func GetConfig() (*ConfigSqitch, error) {
	var err error
	configSqitch := ConfigSqitch{Timeout: DEFAULT_TIMEOUT}

	configSqitch.PostgresUser = utils.Getenv("POSTGRES_IAM_USER", "notset")
	configSqitch.PostgresURI = utils.Getenv("POSTGRES_URI", "notset")
	configSqitch.PostgresPort = utils.Getenv("POSTGRES_PORT", DEFAULT_PORT)
	configSqitch.PostgresDB = utils.Getenv("POSTGRES_DB", "notset")
	configSqitch.Region = utils.Getenv("AWS_REGION", DEFAULT_REGION)
	configSqitch.Profile = utils.Getenv("AWS_PROFILE", "")
	configSqitch.Filepath = utils.Getenv("CONFIG_SQITCH_PATH", DEFAULT_CONFIG_SQITCH_PATH)
	slog.Debug(fmt.Sprintf("ConfigSqitch struct: %v", configSqitch))

	configSqitch.PostgresPassword, err = Connect(&configSqitch)
	if err != nil {
		return nil, err
	}
	slog.Info("Successfully connected to database")

	return &configSqitch, nil
}
