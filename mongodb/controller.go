package mongodb

import (
	"Aitunder/models"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"text/template"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/gomail.v2"
)

var log = logrus.New()

const connectionString = "mongodb+srv://zhangeldy:lemPrXZ1mCeuD0Gn@aitunder.bkn7epv.mongodb.net/?retryWrites=true&w=majority"
const dbName = "aitunder"
const colName = "users"

var collection *mongo.Collection
var projectCollection *mongo.Collection

func init() {
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Warn("Failed to open log file. Logging to standard output.")
	}
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	clientOption := options.Client().ApplyURI(connectionString)
	client, err := mongo.Connect(context.TODO(), clientOption)
	if err != nil {
		log.Warn(err)
		return
	}
	log.Info("MongoDB connection success")

	collection = client.Database(dbName).Collection(colName)
	projectCollection = client.Database(dbName).Collection("projects")
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	allUsers := getFullFromDB()

	json.NewEncoder(w).Encode(allUsers)
}

func AddUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	var user models.User
	_ = json.NewDecoder(r.Body).Decode(&user)
	user.ViewedBy = make([]primitive.ObjectID, 0)
	id, err := insertOneUser(user)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Email already used", "status": 400})
		return
	}

	cookie, err := r.Cookie("sessionID")
	if err != nil {
		cookie = &http.Cookie{
			Name:  "sessionID",
			Value: id,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
	} else {
		cookie.Value = id
		http.SetCookie(w, cookie)
	}
	log.Info("cookie created with id" + cookie.Value)

	err = sendVerificationEmail(user.Email, id)
	if err != nil {
		http.Error(w, "Failed to send verification email", http.StatusInternalServerError)
		log.Error("Failed to send verification email")
		return
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "User registered successfully. Please verify your email.", "status": 200})
	}
}

func AddProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	var project models.Project
	_ = json.NewDecoder(r.Body).Decode(&project)
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		log.Error("UnAuthorised try to add project")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Unauthorised", "status": 400})
		return
	}
	project.Owner, _ = primitive.ObjectIDFromHex(cookie.Value)
	project.ViewedBy = make([]primitive.ObjectID, 0)
	id, err := insertOneProject(project)
	if err != nil {
		http.Error(w, "Failed to add project", http.StatusInternalServerError)
		log.Error("Failed to add project, ", err)
		return
	}
	log.Info("new project with id: ", id)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Project added successfully", "status": 200})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		log.Error("Error reading request body")
		return
	}
	defer r.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error unmarshaling request body", http.StatusBadRequest)
		log.Error("Error unmarshaling request body")
		return
	}

	if data["email"] == "admin@admin.admin" && data["password"] == "Admin123!" {
		log.Info("login Admin")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Admin", "status": 200})
		return
	}

	email, ok := data["email"].(string)
	if !ok {
		http.Error(w, "Invalid email format", http.StatusBadRequest)
		log.Error("Invalid email format")
		return
	}

	password, ok := data["password"].(string)
	if !ok {
		http.Error(w, "Invalid password format", http.StatusBadRequest)
		log.Error("Invalid password format")
		return
	}

	user, err := getOneUserByEmail(email)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Wrong credentials", "status": 400})
		log.Error("Error checking user credentials")
		return
	}

	if !user.AccountActivated {
		http.Error(w, "Account not activated. Please verify your email.", http.StatusUnauthorized)
		return
	}

	cookie, err := r.Cookie("sessionID")
	if err != nil {
		cookie = &http.Cookie{
			Name:  "sessionID",
			Value: user.Id.Hex(),
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		log.Info("New session with id" + cookie.Value)
	}
	if cookie.Value == user.Id.Hex() {
		log.Info("login with cookies")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Login successful", "status": 200})
		return
	} else if user.Id.Hex() != "" {
		cookie.Value = user.Id.Hex()
		cookie.Path = "/"
		http.SetCookie(w, cookie)

		log.Info("New session with id " + cookie.Value)
	}

	if user != nil && user.Password == password {
		log.Info("Login successful")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Login successful", "status": 200})
		return
	} else {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Wrong credentials", "status": 400})
		log.Error("Invalid email or password")
		return
	}
}

var profileTemplate = template.Must(template.ParseFiles("./webPages/templates/home.html"))
var coWorkerTemplate = template.Must(template.ParseFiles("./webPages/templates/coWorkers.html"))
var projectTemplate = template.Must(template.ParseFiles("./webPages/templates/projects.html"))

func ServeProfile(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Error("Unauthorized access")
		return
	}

	user, err := getOneUserByID(cookie.Value)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error("Error retrieving users profiles ", err)
		return
	}
	projects, err := getProjectsByID(cookie.Value)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error("Error retrieving project ", err)
		return
	}

	data := models.UserProjectCombined{
		User:     user,
		Projects: projects,
	}
	err = profileTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error("Error rendering profile template")
		return
	}
}

func ServerCardUsers(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Error("Unauthorized access")
		return
	}
	var projects []models.Project
	projects, _ = getProjectsByID(cookie.Value)
	if projects == nil {
		http.Redirect(w, r, "/project", http.StatusFound)
		return
	}
	user, err := getRandomUser(cookie.Value)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error("Error retrieving user profile", err)
		return
	}

	err = setViewedUs(cookie.Value, user.ID.Hex())
	if err != nil {
		log.Error("cannot set project viewed ", err)
	}

	data := map[string]interface{}{
		"ID":   user.ID.Hex(),
		"Name": user.Name,
		"Profile": map[string]interface{}{
			"Major":             user.Profile.Major,
			"Bio":               user.Profile.Bio,
			"AcademicInterests": user.Profile.AcademicInterests,
			"Skills":            user.Profile.Skills,
		},
	}

	err = coWorkerTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error("Error rendering card-user template")
		return
	}
}

func ServeCardProjects(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Error("Unauthorized access")
		return
	}
	project, err := getRandomProject(cookie.Value)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error("Error retrieving project", err)
		return
	}
	err = setViewedProj(cookie.Value, project.Id.Hex())
	if err != nil {
		log.Error("cannot set project viewed ", err)
	}
	err = projectTemplate.Execute(w, project)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error("Error rendering card-user template")
		return
	}
}

func AddUserProfile(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Allow-Control-Allow-Methods", "POST")

	cookie, err := r.Cookie("sessionID")
	if err != nil {
		log.Error("No session started")
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "No Cookie for ID", "status": 400})
	}

	var profile models.Profile
	err = json.NewDecoder(r.Body).Decode(&profile)
	if err != nil {
		http.Error(w, "Error decoding request body", http.StatusBadRequest)
		log.Error("Error decoding request body")
		return
	}
	err = addProfileToUser(cookie.Value, profile)
	if err != nil {
		json.NewEncoder(w).Encode(map[string]interface{}{"message": "Error updating profile", "status": 400})
		log.Error("Error updating user profile")
		return
	}
	log.Info("Profile updated for id " + cookie.Value)
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Profile updated successfully", "status": 200})
}

func sendVerificationEmail(email, userID string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "aitunderapp.notifications@gmail.com")
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Account Verification")
	mailer.SetBody("text/html", fmt.Sprintf("Please click the following link to verify your account: <a href='http://localhost:8080/verify?userID=%s'>Verify</a>", userID))
	dialer := gomail.NewDialer("smtp.gmail.com", 587, "aitunderapp.notifications@gmail.com", "hbgr gnxq enfr zmtn")
	if err := dialer.DialAndSend(mailer); err != nil {
		return err
	}

	return nil
}

func VerifyAccount(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	userID := r.URL.Query().Get("userID")
	if userID == "" {
		log.Error("User ID not provided")
		http.Error(w, "User ID not provided", http.StatusBadRequest)
		return
	}
	if err := updateOneUserByID(userID); err != nil {
		http.Error(w, "Failed to verify account", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/profile", http.StatusFound)
}

func SendNotificationToUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		log.Error("Error reading request body")
		return
	}
	defer r.Body.Close()

	notificationMessage := string(body)

	authorizedUsers, err := getAllAuthorizedUsers()
	if err != nil {
		http.Error(w, "Error getting authorized users", http.StatusInternalServerError)
		log.Error("Error getting authorized users from the database")
		return
	}

	for _, user := range authorizedUsers {
		err := sendNotificationEmail(user.Email, notificationMessage)
		if err != nil {
			log.Error("Failed to send notification email to user:", user.Email)
			continue
		}
	}

	// Respond with success message
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "Notification sent to all authorized users", "status": 200})
}

func LikeUser(w http.ResponseWriter, r *http.Request) {
	var data map[string]string
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		log.Error("Error reading request body")
		return
	}
	defer r.Body.Close()

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Error("cannot unmarshal request")
		return
	}

	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		log.Error("Unauthorized access")
		return
	}
	loggedInUserID := cookie.Value

	err = addLikedUser(loggedInUserID, data["userID"])
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Error("Error adding liked user", err)
		return
	}

	// You can send a success response if needed
	json.NewEncoder(w).Encode(map[string]interface{}{"message": "User liked", "status": 200})
	log.Error("user ", data["userID"], " was liked by ", loggedInUserID)
}
