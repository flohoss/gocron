#!/bin/sh

action=$1

case $action in
"install")
    go install github.com/swaggo/swag/cmd/swag@latest
    ;;
"init")
    swag init --dir internal/controller -g ../router/router.go --pd --requiredByDefault
    ;;
"format")
    swag fmt
    ;;
*)
    exit 0
    ;;
esac
