.PHONY:start
start:
	CONFIG_PATH=./configs/local.yaml go run ./cmd/sso/main.go

.PHONY:migrate
migrate:
	go run ./cmd/migrator/main.go --storage-path=./storage/sso.db --migrations-path=./migrations

DEFAULT_GOAL := start