package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"name,omitempty"`
	Email     string             `json:"email,omitempty"`
	Password  string             `json:"password,omitempty"`
	Portfolio Profile            `json:"portfolio,omitempty"`
}

type Profile struct {
	UserID            primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	DateOfBirth       string             `json:"dateOfBirth,omitempty" bson:"dateOfBirth,omitempty"`
	Major             string             `json:"major,omitempty"`
	Bio               string             `json:"bio,omitempty"`
	AcademicInterests string             `json:"academicInterests,omitempty"`
	Skills            []string           `json:"skills,omitempty"`
	SocialLinks       []string           `json:"socialLinks,omitempty"`
}

type UserFull struct {
	ID       primitive.ObjectID `json:"_id" bson:"_id"`
	Name     string             `json:"name" bson:"name"`
	Email    string             `json:"email" bson:"email"`
	Password string             `json:"password" bson:"password"`
	Profile  struct {
		DateOfBirth       string   `json:"dateOfBirth" bson:"dateOfBirth"`
		Major             string   `json:"major" bson:"major"`
		Bio               string   `json:"bio" bson:"bio"`
		AcademicInterests string   `json:"academicinterests" bson:"academicinterests"`
		Skills            []string `json:"skills" bson:"skills"`
		SocialLinks       []string `json:"sociallinks" bson:"sociallinks"`
	} `json:"profile" bson:"profile"`
}
