package main

import (
	"Aitunder/internal/models"
	"Aitunder/internal/mongodb"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestCreateUserCardData(t *testing.T) {
	user := models.UserFull{
		ID:    primitive.NewObjectID(),
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Profile: struct {
			DateOfBirth       string   `json:"dateOfBirth" bson:"dateOfBirth"`
			Major             string   `json:"major" bson:"major"`
			PhoneNumber       string   `json:"phoneNumber,omitempty"`
			Bio               string   `json:"bio" bson:"bio"`
			AcademicInterests string   `json:"academicinterests" bson:"academicinterests"`
			Skills            []string `json:"skills" bson:"skills"`
			SocialLinks       []string `json:"sociallinks" bson:"sociallinks"`
		}{
			DateOfBirth:       "1990-01-01",
			Major:             "Computer Science",
			PhoneNumber:       "123456789",
			Bio:               "A bio about John Doe",
			AcademicInterests: "Programming",
			Skills:            []string{"Golang", "Java", "Python"},
			SocialLinks:       []string{"https://twitter.com/johndoe"},
		},
	}

	result := mongodb.CreateUserCardData(user)

	expectedResult := map[string]interface{}{
		"ID":   user.ID.Hex(),
		"Name": user.Name,
		"Profile": map[string]interface{}{
			"Major":             user.Profile.Major,
			"Bio":               user.Profile.Bio,
			"AcademicInterests": user.Profile.AcademicInterests,
			"Skills":            user.Profile.Skills,
		},
	}

	assert.True(t, reflect.DeepEqual(result, expectedResult), "Expected and actual results do not match")
}
