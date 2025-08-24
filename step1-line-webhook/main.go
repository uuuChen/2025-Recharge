package main

import (
	"log"
	"net/http"
)

func main() {
	// 開 api 來處理 LINE webhook
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		log.Printf("handle webhook request")
		w.WriteHeader(200)
	})

	// 監聽 8080 port
	log.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
