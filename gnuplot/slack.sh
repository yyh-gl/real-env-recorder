curl -X POST $WEBHOOK_URL \
     -H 'Content-Type: application/json; charset=UTF-8' \
     -d '{
    "attachments": [
      {
          "text": "Hello world :tada:",
          "image_url": "https://growthseed.jp/wp-content/uploads/2016/12/peach-1.jpg"
      }
  ]
}'
