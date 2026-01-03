#!/bin/bash

# Source venv
source /app/receipt-wrangler-api/wranglervenv/bin/activate

# Start api
cd /app/receipt-wrangler-api
./api --env prod &

# Start nginx
nginx -g 'daemon off;'

