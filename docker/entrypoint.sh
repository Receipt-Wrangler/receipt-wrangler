#!/bin/bash

# Source venv
source /app/receipt-wrangler-api/wranglervenv/bin/activate

cd /app/receipt-wrangler-api

# Forward TERM/INT to the children: bash runs as PID 1 here, and PID 1 ignores
# SIGTERM by default, so `docker stop` previously waited out the full grace
# period and then SIGKILLed everything.
trap 'kill $(jobs -p) 2>/dev/null' TERM INT

# Run the API and nginx as background jobs and exit the container if EITHER one
# exits, so the restart policy / orchestrator can recover it. Previously the API
# ran in the background while only nginx held the foreground, so a crashed API
# left the container "up" serving 502s and the restart policy never fired.
./api --env prod &
nginx -g 'daemon off;' &

wait -n
exit $?

