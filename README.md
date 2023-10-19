# FinalcadOne-sqitch

Create configuration file for sqitch containing postgres connection credentials using rds IAM authentification.
By default configuration file will be named `sqitch.conf` in current directory.

### Environment variables

| Name                                  | Description                                                 | Default     |
|---------------------------------------|-------------------------------------------------------------|-------------|
| AWS_REGION                            | AWS region                                                  | us-east-1   |
| AWS_PROFILE                           | AWS profile (Can be null if service use IRSA authentication)| empty       |
| POSTGRES_IAM_USER                     | IAM user                                                    | empty       |
| POSTGRES_URI                          | Postgres URI endpoint                                       | empty       |
| POSTGRES_DB                           | Database name                                               | empty       |
| CONFIG_SQITCH_PATH                    | config filepath + name                                      | sqitch.conf |

## Usage in other container

```yaml
COPY --from XXX /sqitch-config .
```
