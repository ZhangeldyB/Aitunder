package mongodb

import (
	"Aitunder/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	err := insertOneUser(user)

	cookie := http.Cookie{
		Name:  "sessionID",
		Value: user.Id.Hex(),
	}
	http.SetCookie(w, &cookie)
	fmt.Println("cookie created with id" + cookie.Value)

	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Email already used", "status": 400})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Account created successfully", "status": 200})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		fmt.Println("Error reading request body")
		return
	}
	defer r.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error unmarshaling request body", http.StatusBadRequest)
		fmt.Println("Error unmarshaling request body")
		return
	}

	email, ok := data["email"].(string)
	if !ok {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		fmt.Println("Invalid email format")
		return
	}

	password, ok := data["password"].(string)
	if !ok {
		http.Error(w, "Invalid password format", http.StatusBadRequest)
		fmt.Println("Invalid password format")
		return
	}

	user, err := getOneUserByEmail(email)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Wrong credentials", "status": 400})
		fmt.Println("Error checking user credentials")
		return
	}
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		cookie = &http.Cookie{
			Name:  "sessionID",
			Value: user.Id.Hex(),
		}
		http.SetCookie(w, cookie)
		fmt.Println("cookie created with id" + cookie.Value)
	}
	if cookie.Value == user.Id.Hex() {
		fmt.Println("Login successful cokies")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Login successful", "status": 200})
		return
	} else {
		cookie = &http.Cookie{
			Name:  "sessionID",
			Value: user.Id.Hex(),
		}
		http.SetCookie(w, cookie)
		fmt.Println("cookie created with id" + cookie.Value)
	}

	if user != nil && user.Password == password {
		fmt.Println("Login successful")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Login successful", "status": 200})
		return
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Wrong credentials", "status": 400})
		fmt.Println("Invalid email or password")
		return
	}
}

func AddUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	cookie, err := r.Cookie("sessionID")
	if err != nil {
		fmt.Println("no cookie for id")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "No Cookie for ID", "status": 400})
	}

	var profile models.Portfolio
	err = json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		fmt.Println("Error decoding request body")
		return
	}
	err = addProfileToUser(cookie.Value, profile)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Error updating profile", "status": 500})
		fmt.Println("Error updating user profile")
		return
	}
	fmt.Println("Profile updated successfully")
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Profile updated successfully", "status": 200})
}
