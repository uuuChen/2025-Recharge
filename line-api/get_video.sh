#!/bin/bash

# -----------------------------------------------------------------------------
# 使用
# . ./line-api/get_video.sh <message_id>
# -----------------------------------------------------------------------------

# 載入 .env 檔案
if [ -f .env ]; then
  export $(cat .env | grep -v '^#' | xargs)
fi

curl -L -H "Authorization: Bearer $LINE_CHANNEL_ACCESS_TOKEN" \
  -D - \
  "https://api-data.line.me/v2/bot/message/$1/content" \
  -o ./line-api/video-$1.mp4