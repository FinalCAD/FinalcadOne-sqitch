# FinalcadOne-sqitch

Create configuration file and a secret file for sqitch containing postgres connection credentials using rds IAM authentification.
By default configuration file will be named `sqitch.conf` in current directory.
 Password to load in environment variable will be located by default in `secret.conf`

### Environment variables

| Name                                  | Description                                                 | Default     |
|---------------------------------------|-------------------------------------------------------------|-------------|
| AWS_REGION                            | AWS region                                                  | us-east-1   |
| AWS_PROFILE                           | AWS profile (Can be null if service use IRSA authentication)| empty       |
| POSTGRES_URI                          | Postgres URI endpoint                                       | empty       |
| POSTGRES_DB                           | Database name                                               | empty       |
| CONFIG_SQITCH_PATH                    | config filepath + name                                      | sqitch.conf |
| SECRET_SQITCH_PATH                    | config filepath + name                                      | secret.conf |

## User environment variable

IAM database authentication is enable if POSTGRES_IAM_USER is not empty, you can only use one method.

Postgres users

| Name                                  | Description                                                 | Default     |
|---------------------------------------|-------------------------------------------------------------|-------------|
| POSTGRES_USER                         | User                                                        | empty       |
| POSTGRES_PASSWORD                     | Password                                                    | empty       |

IAM database authentication

| Name                                  | Description                                                 | Default     |
|---------------------------------------|-------------------------------------------------------------|-------------|
| POSTGRES_IAM_USER                     | IAM user, activate token auth ge                            | empty       |

## Usage in other container

```yaml
# add binary from image
COPY --from XXX /sqitch-config .
# in entrypoint.sh
# load secret in env vars
source ./secret.conf
```
