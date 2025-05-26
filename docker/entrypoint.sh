#!/bin/sh

set -e

CONFIG_FILE="/app/config/config.yml" # adjust as needed
APP="./gocron"

log_message() {
    local message="$1"
    local timestamp=$(date +"%Y-%m-%dT%H:%M:%S.%N%:z")
    echo "{\"time\":\"${timestamp}\",\"message\":\"${message}\"}"
}

cat logo.txt

while true; do
    log_message "Starting GoCron..."
    # Start the app in background
    $APP &
    PID=$!
    log_message "GoCron started with PID: $PID"

    log_message "Waiting for $CONFIG_FILE to exist..."
    while [ ! -f "$CONFIG_FILE" ]; do
        sleep 1
    done
    log_message "$CONFIG_FILE now exists. Proceeding..."

    log_message "Watching for changes in $CONFIG_FILE..."
    # inotifywait quiet mode: -qq suppresses all output
    inotifywait -qq -e modify "$CONFIG_FILE"
    log_message "Detected change in $CONFIG_FILE."

    log_message "Config changed. Reloading..."

    # Stop the app
    log_message "Killing GoCron (PID: $PID)..."
    kill $PID
    log_message "Waiting for GoCron (PID: $PID) to terminate..."
    wait $PID
    log_message "GoCron (PID: $PID) terminated. Restarting loop..."
done

log_message "Entrypoint script exited. This should not happen in the loop." # This line should ideally not be reached
