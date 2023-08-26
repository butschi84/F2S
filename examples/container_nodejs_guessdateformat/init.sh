#!/bin/bash

npm start &
./f2sfizzlet &

# wait for any process to exit
wait -n

# exit with status of process that exited first
exit $?