package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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

		// ============================================================
		// 處理收到的 events
		log.Println("---------------------------------------------------------------")
		for _, event := range cb.Events {
			formattedEvent, _ := json.MarshalIndent(event, "", "  ")
			log.Printf("event: %s", string(formattedEvent))

			switch e := event.(type) {
			case webhook.FollowEvent:
				log.Printf("[Follow] UserID: %+v", e.Source.(webhook.UserSource).UserId)
			case webhook.UnfollowEvent:
				log.Printf("[Unfollow] UserID: %+v", e.Source.(webhook.UserSource).UserId)
			case webhook.JoinEvent:
				log.Printf("[Join] GroupID: %+v", e.Source.(webhook.GroupSource).GroupId)
			case webhook.LeaveEvent:
				log.Printf("[Leave] GroupID: %+v", e.Source.(webhook.GroupSource).GroupId)
			case webhook.MemberJoinedEvent:
				log.Printf("[MemberJoined] Join users: %+v", e.Joined.Members)
			case webhook.MemberLeftEvent:
				log.Printf("[MemberLeft] Left users: %+v", e.Left.Members)
			case webhook.MessageEvent:
				log.Printf("[Message]: %+v", e.Message)

				// var userID, groupID string
				// switch source := e.Source.(type) {
				// case webhook.UserSource: // 來自個人聊天室
				// 	userID = source.UserId
				// case webhook.GroupSource: // 來自群組聊天室
				// 	userID = source.UserId
				// 	groupID = source.GroupId
				// }
				// log.Printf("[Message] GroupID: %+v, UserID: %+v", groupID, userID)

				// switch message := e.Message.(type) {
				// case webhook.TextMessageContent:
				// 	log.Printf("[Message][Text] %+v", message.Text)
				// case webhook.ImageMessageContent:
				// 	log.Printf("[Message][Image] MessageID: %+v", message.Id)
				// case webhook.AudioMessageContent:
				// 	log.Printf("[Message][Audio] MessageID: %+v", message.Id)
				// case webhook.VideoMessageContent:
				// 	log.Printf("[Message][Video] MessageID: %+v", message.Id)
				// case webhook.FileMessageContent:
				// 	log.Printf("[Message][File] MessageID: %+v", message.Id)
				// case webhook.LocationMessageContent:
				// 	log.Printf("[Message][Location] Longitude: %v, Latitude: %v", message.Longitude, message.Latitude)
				// case webhook.StickerMessageContent:
				// 	log.Printf("[Message][Sticker] PackageID: %v, StickerID: %v", message.PackageId, message.StickerId)
				// }
			}
		}
		// ============================================================
	})

	// 監聽 8080 port
	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// 1. message id 可以用來下載 image、audio、video
