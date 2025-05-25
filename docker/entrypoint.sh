#!/bin/sh

set -e

CONFIG_FILE="/app/config/config.yml" # adjust as needed
APP="./gocron"

# Function to generate log messages with Echo-like format
log_message() {
    level="$1"
    message="$2"
    timestamp=$(date +"%Y-%m-%dT%H:%M:%S.%N")
    tz_offset=$(date +"%:::z")
    echo "{\"time\":\"${timestamp}${tz_offset}\",\"level\":\"${level}\",\"prefix\":\"entrypoint\",\"message\":\"${message}\"}"
}

cat logo.txt

log_message "INFO" "Entrypoint script starting..."

while true; do
    log_message "INFO" "Starting GoCron..."
    # Start the app in background
    $APP &
    PID=$!
    log_message "INFO" "GoCron started with PID: $PID"

    log_message "INFO" "Watching for changes in $CONFIG_FILE..."
    # inotifywait quiet mode: -qq suppresses all output
    inotifywait -qq -e modify "$CONFIG_FILE"
    log_message "INFO" "Detected change in $CONFIG_FILE."

    log_message "INFO" "Config changed. Reloading..."

    # Stop the app
    log_message "INFO" "Killing GoCron (PID: $PID)..."
    kill $PID
    log_message "INFO" "Waiting for GoCron (PID: $PID) to terminate..."
    wait $PID
    log_message "INFO" "GoCron (PID: $PID) terminated. Restarting loop..."
done

log_message "INFO" "Entrypoint script exited. This should not happen in the loop." # This line should ideally not be reached
