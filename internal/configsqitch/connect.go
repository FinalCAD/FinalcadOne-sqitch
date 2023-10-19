package configsqitch

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/retry"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/lib/pq"
)

func Connect(configSqitch *ConfigSqitch) (string, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(configSqitch.Timeout)*time.Millisecond)
	defer cancel()
	var err error
	var cfg aws.Config

	if configSqitch.Profile != "" {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithSharedConfigProfile(configSqitch.Profile), config.WithRegion(configSqitch.Region))
	} else {
		cfg, err = config.LoadDefaultConfig(ctx, config.WithRegion(configSqitch.Region),
			config.WithRetryer(func() aws.Retryer { return retry.AddWithMaxAttempts(aws.NopRetryer{}, 1) }))
	}

	if err != nil {
		return "", fmt.Errorf("configuration error: %s", err.Error())
	}

	authenticationToken, err := auth.BuildAuthToken(
		context.TODO(), configSqitch.PostgresURI+":"+configSqitch.PostgresPort, configSqitch.Region, configSqitch.PostgresUser, cfg.Credentials)
	if err != nil {
		return "", fmt.Errorf("failed to create authentication token: %s", err.Error())
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		configSqitch.PostgresURI, configSqitch.PostgresPort, configSqitch.PostgresUser, authenticationToken, configSqitch.PostgresDB,
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
