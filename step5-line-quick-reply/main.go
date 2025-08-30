package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/v8/linebot/messaging_api"
	"github.com/line/line-bot-sdk-go/v8/linebot/webhook"
)

// ============================================================
// 設定 quick reply
func buildQuickReply() *messaging_api.QuickReply {
	return &messaging_api.QuickReply{
		Items: []messaging_api.QuickReplyItem{ // quick reply 最多 13 個
			// 文字訊息 action
			{Action: messaging_api.MessageAction{Label: "貼圖", Text: "/貼圖"}},
			{Action: messaging_api.MessageAction{Label: "圖片", Text: "/圖片"}},
			{Action: messaging_api.MessageAction{Label: "影片", Text: "/影片"}},
			{Action: messaging_api.MessageAction{Label: "位置", Text: "/位置"}},

			// URI action
			{Action: messaging_api.UriAction{Label: "開 MYEDIT", Uri: "https://myedit.online/"}},

			// DatetimePicker action
			{Action: messaging_api.DatetimePickerAction{Label: "選日期時間", Data: "2025-01-01", Mode: "datetime"}},

			// 只有 quick-reply 可以用的 action
			{Action: messaging_api.CameraAction{Label: "開相機"}},
			{Action: messaging_api.CameraRollAction{Label: "開相簿"}},
			{Action: messaging_api.LocationAction{Label: "傳位置"}},
		},
	}
}

// ============================================================

func main() {
	// 載入 .env 檔案
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// 接收 line webhook 時用來驗證 x-line-signature
	channelSecret := os.Getenv("LINE_CHANNEL_SECRET")

	// 初始化 LINE bot
	channelAccessToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	bot, _ := messaging_api.NewMessagingApiAPI(channelAccessToken)
	if err != nil {
		log.Fatal("Error initializing LINE bot: ", err)
	}

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
				case webhook.TextMessageContent:
					log.Printf("[Message][Text] %+v", message.Text)

					var respMsg messaging_api.MessageInterface
					switch message.Text {
					case "/貼圖":
						respMsg = messaging_api.StickerMessage{
							PackageId:  "6632",
							StickerId:  "11825375",
							QuickReply: buildQuickReply(), // quick reply
						}
					case "/圖片":
						respMsg = messaging_api.ImageMessage{
							OriginalContentUrl: "https://leo-recharge.s3.ap-southeast-2.amazonaws.com/public/1.jpeg",
							PreviewImageUrl:    "https://leo-recharge.s3.ap-southeast-2.amazonaws.com/public/1.jpeg",
							QuickReply:         buildQuickReply(), // quick reply
						}
					case "/影片":
						respMsg = messaging_api.VideoMessage{
							OriginalContentUrl: "https://leo-recharge.s3.ap-southeast-2.amazonaws.com/public/video1.mp4",
							PreviewImageUrl:    "https://leo-recharge.s3.ap-southeast-2.amazonaws.com/public/video1-thumbnail.png",
							QuickReply:         buildQuickReply(), // quick reply
						}
					case "/位置":
						respMsg = messaging_api.LocationMessage{
							Title:      "cyberlink",
							Address:    "新北市新店區民權路100號16樓",
							Latitude:   24.983348636569882,
							Longitude:  121.53610649999997,
							QuickReply: buildQuickReply(), // quick reply
						}
					default:
						respMsg = messaging_api.TextMessage{
							Text:       message.Text,
							QuickReply: buildQuickReply(), // quick reply
						}
					}

					_, err := bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							respMsg,
						},
					})
					if err != nil {
						log.Printf("Error replying message: %v", err)
					}
				default:
					// 非文字訊息，回覆支援的指令
					_, err := bot.ReplyMessage(&messaging_api.ReplyMessageRequest{
						ReplyToken: e.ReplyToken,
						Messages: []messaging_api.MessageInterface{
							messaging_api.TextMessage{
								Text:       "支援的指令：\n/貼圖\n/圖片\n/影片\n/位置",
								QuickReply: buildQuickReply(),
							},
						},
					})
					if err != nil {
						log.Printf("Error replying message: %v", err)
					}
				}
			case webhook.PostbackEvent:
				formattedEvent, _ := json.MarshalIndent(e, "", "  ")
				log.Printf("[Postback] %+v", string(formattedEvent))
			}
		}
	})

	// 監聽 8080 port
	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
