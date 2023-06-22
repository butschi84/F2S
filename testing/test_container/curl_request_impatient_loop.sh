#!/bin/bash

end_time=$((SECONDS + 60))  # Set end time to current time + 60 seconds

while [ $SECONDS -lt $end_time ]; do
    # Perform your curl request here
    curl http://192.168.2.40:32611/invoke/testblocking &

    # Add a delay between each request if desired
    sleep 1
done