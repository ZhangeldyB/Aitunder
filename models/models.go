package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id               primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Name             string               `json:"name,omitempty"`
	Email            string               `json:"email,omitempty"`
	Password         string               `json:"password,omitempty"`
	Profile          Profile              `json:"profile,omitempty"`
	AccountActivated bool                 `json:"accountActivated,omitempty" bson:"accountActivated,omitempty"`
	ViewedUsers      []primitive.ObjectID `json:"viewedUsers,omitempty" bson:"viewedUsers,omitempty"`
}

type Profile struct {
	UserID            primitive.ObjectID `json:"userId,omitempty" bson:"userId,omitempty"`
	DateOfBirth       string             `json:"dateOfBirth,omitempty" bson:"dateOfBirth,omitempty"`
	Major             string             `json:"major,omitempty"`
	PhoneNumbber      string             `json:"phoneNumber,omitempty"`
	Bio               string             `json:"bio,omitempty"`
	AcademicInterests string             `json:"academicInterests,omitempty"`
	Skills            []string           `json:"skills,omitempty"`
	SocialLinks       []string           `json:"socialLinks,omitempty"`
}

type UserFull struct {
	ID      primitive.ObjectID `json:"_id" bson:"_id"`
	Name    string             `json:"name" bson:"name"`
	Email   string             `json:"email" bson:"email"`
	Profile struct {
		DateOfBirth       string   `json:"dateOfBirth" bson:"dateOfBirth"`
		Major             string   `json:"major" bson:"major"`
		PhoneNumber       string   `json:"phoneNumber,omitempty"`
		Bio               string   `json:"bio" bson:"bio"`
		AcademicInterests string   `json:"academicinterests" bson:"academicinterests"`
		Skills            []string `json:"skills" bson:"skills"`
		SocialLinks       []string `json:"sociallinks" bson:"sociallinks"`
	} `json:"profile" bson:"profile"`
	ViewedUsers []primitive.ObjectID `json:"viewedUsers,omitempty" bson:"viewedUsers,omitempty"`
}

type Project struct {
	Id          primitive.ObjectID   `json:"_id,omitempty" bson:"_id,omitempty"`
	Title       string               `json:"title,omitempty" bson:"title,omitempty"`
	Description string               `json:"description,omitempty" bson:"description,omitempty"`
	Skills      []string             `json:"skills,omitempty" bson:"skills,omitempty"`
	Owner       primitive.ObjectID   `json:"owner,omitempty" bson:"owner,omitempty"`
	ViewedBy    []primitive.ObjectID `json:"viewedBy,omitempty" bson:"viewedBy,omitempty"`
}

type UserProjectCombined struct {
	User     UserFull  `json:"user" bson:"user"`
	Projects []Project `json:"project" bson:"project"`
}
