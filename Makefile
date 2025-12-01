.PHONY:start
start:
	CONFIG_PATH=./configs/local.yaml go run ./cmd/sso/main.go

DEFAULT_GOAL := start