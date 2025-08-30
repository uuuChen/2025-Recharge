package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

func main() {
	// 載入 .env 檔案
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 接收 line webhook 時用來驗證 x-line-signature
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")

	// ============================================================
	// 初始化 LINE bot
	channelAccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	bot, _ := messaging_api.NewMessagingApiAPI(channelAccessToken)
	if err != nil {
		log.Fatal("Error initializing LINE bot: ", err)
	}
	// ============================================================

	// 開 api 來處理 LINE webhook
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		// 驗證 signature 並解析 request
		cb, err := webhook.ParseRequest(channelSecret, req)
		if err != nil {
			if errors.Is(err, webhook.ErrInvalidSignature) {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		// 如果回 5** 或是 timeout 且有設定 redeliver，LINE 會 redeliver
		w.WriteHeader(200)

		// 處理收到的 events
		for _, event := range cb.Events {
			switch e := event.(type) {
			case webhook.MessageEvent:
				switch message := e.Message.(type) {
				// ============================================================
				case webhook.TextMessageContent:
					log.Printf("[Message][Text] %+v", message.Text)

					var respMsg messaging_api.MessageInterface
					switch message.Text {
					case "/貼圖":
						respMsg = messaging_api.StickerMessage{
							PackageId: "6632",
							StickerId: "11825375",
						}
					case "/圖片":
						respMsg = messaging_api.ImageMessage{
							OriginalContentUrl: "https://leo-recharge.s3.ap-southeast-2.amazonaws.com/public/1.jpeg",
							PreviewImageUrl:    "https://leo-recharge.s3.ap-southeast-2.amazonaws.com/public/1.jpeg",
						}
					case "/影片":
						respMsg = messaging_api.VideoMessage{
							OriginalContentUrl: "https://leo-recharge.s3.ap-southeast-2.amazonaws.com/public/video1.mp4",
							PreviewImageUrl:    "https://leo-recharge.s3.ap-southeast-2.amazonaws.com/public/video1-thumbnail.png",
						}
					case "/位置":
						respMsg = messaging_api.LocationMessage{
							Title:     "cyberlink",
							Address:   "新北市新店區民權路100號16樓",
							Latitude:  24.983348636569882,
							Longitude: 121.53610649999997,
						}
					default:
						respMsg = messaging_api.TextMessage{
							Text: message.Text,
						}
					}
					_, err := bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{ // 最多 5 則
							respMsg,
							respMsg,
							respMsg,
							respMsg,
							respMsg,
						},
					})
					if err != nil {
						log.Printf("Error replying message: %v", err)
					}

					// push message
					// xLineRetryID := uuid.New().String()
					// _, err = bot.PushMessage(&messaging_api.PushMessageRequest{
					// 	To: e.Source.(webhook.UserSource).UserId,
					// 	Messages: []messaging_api.MessageInterface{
					// 		messaging_api.TextMessage{
					// 			Text: "Push message",
					// 		},
					// 	},
					// }, xLineRetryID)
					// if err != nil {
					// 	log.Printf("Error pushing message: %v", err)
					// }
				}
				// ============================================================
			}
		}
	})

	// 監聽 8080 port
	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// 1. ReplyMessage 最多 5 則
// 2. ReplyToken 最多使用一次
