#!/bin/bash

for ((i=1; i<=10; i++)); do
    bash curl_request_impatient_loop.sh &
done

# Wait for all background processes to finish
wait
