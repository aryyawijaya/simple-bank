## Run server

```bash
# prod
$ make compose-up-prod

# dev
$ make compose-up-dev
```

## Log the app

```bash
$ make logs-app
```

## Test

```bash
$ make test
```

## Generate code

```bash
# Create new db migration
$ migrate create -ext sql -dir db/migration -seq <migration_name>
```
