package websocket

import (
	"Aitunder/internal/models"
	"Aitunder/internal/mongodb"
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	clients sync.Map
)

type Message struct {
	RecipientID string `json:"recipient"`
	SenderName  string `json:"senderName"`
	Text        string `json:"text"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	cookie, err := r.Cookie("sessionID")
	if err != nil {
		log.Println("Failed to get session ID from cookie:", err)
		return
	}
	userID := cookie.Value

	params := r.URL.Query()
	otherUserID := params.Get("otherUserID")
	if otherUserID == "" {
		log.Println("otherUserID not provided")
		return
	}

	messages, err := mongodb.FetchChatHistory(userID, otherUserID)
	if err != nil {
		log.Println("Error fetching chat history:", err)
		return
	}

	for _, msg := range messages {
		message := string(msg.SenderName) + ": " + msg.Message
		err := ws.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("Error sending chat history message:", err)
			return
		}
	}

	clients.Store(userID, ws)

	for {
		var msg Message
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			clients.Delete(userID)
			break
		}

		go handleMessage(userID, msg)
	}
}

func handleMessage(senderID string, msg Message) {

	recipientID := msg.RecipientID
	senderName := msg.SenderName
	text := msg.Text
	msgToSend := senderName + ": " + text

	senderObjectID, err := primitive.ObjectIDFromHex(senderID)
	if err != nil {
		log.Println("Invalid senderID:", senderID)
		return
	}
	recipientObjectID, err := primitive.ObjectIDFromHex(recipientID)
	if err != nil {
		log.Println("Invalid recipientID:", recipientID)
		return
	}

	chatMessage := models.ChatMessage{
		RecipientID: recipientObjectID,
		SenderID:    senderObjectID,
		SenderName:  senderName,
		Message:     text,
		Time:        primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err = mongodb.ChatsCollection.InsertOne(context.Background(), chatMessage)
	if err != nil {
		log.Println("Error inserting chat message:", err)
	}

	if conn, ok := clients.Load(recipientID); ok {
		wsConn := conn.(*websocket.Conn)
		err := wsConn.WriteMessage(websocket.TextMessage, []byte(msgToSend))
		if err != nil {
			log.Println("Error writing message to recipient:", err)
		}
	} else {
		log.Println("Recipient not connected:", recipientID)
	}

	if conn, ok := clients.Load(senderID); ok {
		wsConn := conn.(*websocket.Conn)
		err := wsConn.WriteMessage(websocket.TextMessage, []byte(msgToSend))
		if err != nil {
			log.Println("Error writing message to sender:", err)
		}
	} else {
		log.Println("Sender not connected:", senderID)
	}
}
