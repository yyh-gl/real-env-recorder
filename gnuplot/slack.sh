curl -X POST $WEBHOOK_URL_TO_02 \
     -H 'Content-Type: application/json; charset=UTF-8' \
     -d '{
    "attachments": [
      {
          "text": "気温レポート",
          "image_url": "https://super.hobigon.work/public/temperature.png"
      }
  ]
}'
