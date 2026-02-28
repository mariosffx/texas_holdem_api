#!/usr/bin/env bash
set -euo pipefail

export PATH="$PATH:/root/.toolkit/tools/.asdf/shims"
export PATH="$PATH:/root/.toolkit/bin/Linux/aarch64"

source /root/.bashrc

asdf plugin add golang https://github.com/asdf-community/asdf-golang.git
asdf install golang 1.26.0
asdf reshim
asdf reshim golang 1.26.0




echo "[post-deploy] Building project..."
npm install

echo "[post-deploy] Building Go application..."
go build -o ./texas_holdem_api . 2>&1 || {
  EXIT_CODE=$?
  echo "[ERROR] Go build failed with exit code $EXIT_CODE"
  exit $EXIT_CODE
}

chmod +x ./texas_holdem_api

CONFIG_FILE="ecosystem.config.cjs"
echo "[post-deploy] Starting/Reloading pm2 using $CONFIG_FILE"

APP_NAME="texas_holdem_api"
if pm2 describe "$APP_NAME" >/dev/null 2>&1; then
  echo "[post-deploy] Reloading existing pm2 app: $APP_NAME"
  pm2 reload "$CONFIG_FILE" --update-env
else
  echo "[post-deploy] Starting new pm2 app from: $CONFIG_FILE"
  pm2 start "$CONFIG_FILE"
fi


echo "[post-deploy] Saving pm2 process list..."

pm2 save

echo "[post-deploy] Completed successfully."
