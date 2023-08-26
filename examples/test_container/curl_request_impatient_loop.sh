#!/bin/bash

end_time=$((SECONDS + (5 * 60)))  # Set end time to current time + 60 seconds
output_file="output.txt"  # File to append the results

while [ $SECONDS -lt $end_time ]; do
    # Perform your curl request here
    curl http://192.168.2.40:31468/invoke/testnonblocking?delay=5000 >> "$output_file" &

    # Add a delay between each request if desired
    sleep 1
done