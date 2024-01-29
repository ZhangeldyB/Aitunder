package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id       primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name     string             `json:"name,omitempty"`
	Email    string             `json:"email,omitempty"`
	Password string             `json:"password,omitempty"`
	Portfolio Portfolio          `json:"portfolio,omitempty"`
}


type Portfolio struct {
	UserID           primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	DateOfBirth      string             `json:"dateOfBirth,omitempty" bson:"dateOfBirth,omitempty"`
	Major            string             `json:"major,omitempty"`
	Bio              string             `json:"bio,omitempty"`
	AcademicInterests []string           `json:"academicInterests,omitempty"`
	Skills           []string           `json:"skills,omitempty"`
	SocialLinks      []string           `json:"socialLinks,omitempty"`
}