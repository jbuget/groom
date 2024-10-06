#!/bin/bash

echo $GOOGLE_SERVICE_ACCOUNT_SECRET_FILE | base64 -d > /app/service_account.json
# Start default script for PHP apps
$HOME/bin/groom