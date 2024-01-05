package mongodb

import (
	"Aitunder/models"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const connectionString = "mongodb+srv://zhangeldy:lemPrXZ1mCeuD0Gn@aitunder.bkn7epv.mongodb.net/?retryWrites=true&w=majority"
const dbName = "aitunder"
const colName = "users"

var collection *mongo.Collection

func init() {
	clientOption := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDB connection success")

	collection = client.Database(dbName).Collection(colName)

}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json") //x-www-form-urlencode
	allUsers := getAllUsers()
	json.NewEncoder(w).Encode(allUsers)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	insertOneUser(user)
	json.NewEncoder(w).Encode(user)
}

