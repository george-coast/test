#!/bin/bash

# Variables
WAZUH_API_URL="https://api.jumpcloud.com"
API_TOKEN="jca_8pxpP6acJGxDFAyrtMeCa8tWfv6r54n9Bqhf"
SERVICE="password_manager"
LAST_TIMESTAMP_FILE="last_timestamp.txt"

# Function to get the last timestamp
get_last_timestamp() {
    if [[ -f $LAST_TIMESTAMP_FILE ]]; then
        cat $LAST_TIMESTAMP_FILE
    else
        echo "2020-01-01T14:00:00Z"  # Default start time if file doesn't exist
    fi
}

# Function to update the last timestamp
update_last_timestamp() {
    echo "$1" > $LAST_TIMESTAMP_FILE
}

# Main loop
while true; do
    # Get the last timestamp
    START_TIME=$(get_last_timestamp)

    # Query the API
    RESPONSE=$(curl -s -X POST "$WAZUH_API_URL" \
        -H "Content-Type: application/json" \
        -H "Authorization: Bearer $API_TOKEN" \
        -d '{
            "service": [
                "'"$SERVICE"'"
            ],
            "start_time": "'"$START_TIME"'"
        }')

    # Process the response (this is where you'd handle the events)
    echo "Response: $RESPONSE"

    # Update the last timestamp to the current time
    CURRENT_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
    update_last_timestamp "$CURRENT_TIME"

    # Wait for a specified interval before the next request (e.g., 60 seconds)
    sleep 60
done
