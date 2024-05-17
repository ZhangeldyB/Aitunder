package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type ChatMessage struct {
	RecipientID primitive.ObjectID `json:"recipientID,omitempty" bson:"recipientID,omitempty"`
	SenderID    primitive.ObjectID `json:"senderID,omitempty" bson:"senderID,omitempty"`
	SenderName  string             `json:"senderName,omitempty" bson:"senderName,omitempty"`
	Message     string             `json:"message,omitempty" bson:"message,omitempty"`
	Time        primitive.DateTime `json:"time,omitempty" bson:"time,omitempty"`
}
