package main

import (
	"Aitunder/internal/mongodb"
	"Aitunder/internal/websocket"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	pb "github.com/whateveer/payment-grpc/cmd/grpc-microservice/payment"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
	"google.golang.org/grpc"
)

var log = logrus.New()
var limiter = rate.NewLimiter(3, 5)

func init() {
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		log.Warn("Failed to open log file. Logging to standard output.")
	}

	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetLevel(logrus.InfoLevel)

	log.Info("Logging initialized")
}

func handleWithRateLimit(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		handler(w, r)
	}
}

func pageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		servePage(w, "./templates/registration.html")
	case "/login":
		servePage(w, "./templates/login.html")
	case "/home":
		servePage(w, "./templates/home.html")
	case "/admin":
		servePage(w, "./templates/admin.html")
	case "/profile":
		servePage(w, "./templates/profile.html")
	case "/project":
		servePage(w, "./templates/project.html")
	case "/upgrade":
		servePage(w, "./templates/payment.html")
	default:
		http.NotFound(w, r)
		defer r.Body.Close()
	}
}

func testRequest(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("sessionID")
	if err != nil {
		fmt.Println("there is no cookie")
	} else {
		fmt.Println("there is a cookie with id" + cookie.Value)
	}
	message, _ := io.ReadAll(r.Body)
	w.Write(message)
}

func servePage(w http.ResponseWriter, pagePath string) {
	htmlContent, err := os.ReadFile(pagePath)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(htmlContent)
}

func handlePayment(w http.ResponseWriter, r *http.Request) {
	var paymentRequest struct {
		CardNumber   string `json:"cardNumber"`
		ExpiryDate   string `json:"expiryDate"`
		Cvv          string `json:"cvv"`
		CustomerName string `json:"customerName"`
	}
	err := json.NewDecoder(r.Body).Decode(&paymentRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cookie, err := r.Cookie("sessionID")
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	userID := cookie.Value

	user, err := mongodb.GetOneUserByID(userID)
	if err != nil {
		http.Error(w, "Failed to get user details", http.StatusInternalServerError)
		return
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		http.Error(w, "Failed to connect to gRPC server", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	client := pb.NewPaymentServiceClient(conn)
	resp, err := client.ProcessPayment(context.Background(), &pb.PaymentRequest{
		CardNumber:   paymentRequest.CardNumber,
		ExpiryDate:   paymentRequest.ExpiryDate,
		Cvv:          paymentRequest.Cvv,
		CustomerName: paymentRequest.CustomerName,
		Email:        user.Email, // Add email from user details
		FullName:     user.Name,  // Add full name from user details
	})
	if err != nil {
		http.Error(w, "Payment failed", http.StatusInternalServerError)
		return
	}

	// Update user details
	err = mongodb.UpdateUserAfterPayment(userID)
	if err != nil {
		http.Error(w, "Failed to update user details", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, `{"status": "success", "transactionId": "%s"}`, resp.TransactionId)
}

func main() {
	http.HandleFunc("/login", pageHandler)
	http.HandleFunc("/admin", pageHandler)
	http.HandleFunc("/profile", pageHandler)
	http.HandleFunc("/project", pageHandler)
	http.HandleFunc("/upgrade", pageHandler)
	http.HandleFunc("/api/test", testRequest)
	http.HandleFunc("/api/pay", handlePayment)
	http.HandleFunc("/api/signUp", mongodb.AddUser)
	http.HandleFunc("/api/login", mongodb.LoginHandler)
	http.HandleFunc("/api/sendNotification", mongodb.SendNotificationToUsers)
	http.HandleFunc("/card/Co-Worker", mongodb.ServerCardUsers)
	http.HandleFunc("/card/Project", mongodb.ServeCardProjects)
	// http.HandleFunc("/card/LikeProject", mongodb.LikeProject)
	http.HandleFunc("/card/LikeUser", mongodb.LikeUser)
	http.HandleFunc("/home", mongodb.ServeProfile)
	http.HandleFunc("/api/profile/add", mongodb.AddUserProfile)
	http.HandleFunc("/api/project/add", mongodb.AddProject)
	http.HandleFunc("/api/getAllUsers", mongodb.GetAllUsers)
	http.HandleFunc("/verify", mongodb.VerifyAccount)
	http.HandleFunc("/ws", websocket.HandleConnections)
	http.HandleFunc("/unlike", mongodb.UnlikeUser)
	http.HandleFunc("/", handleWithRateLimit(pageHandler))
	fmt.Println("Server is running on http://localhost:8080/main")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error(err)
	}
}
