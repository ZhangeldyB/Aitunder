package main

import (
	"Aitunder/mongodb"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"golang.org/x/time/rate"
)

var log = logrus.New()
var limiter = rate.NewLimiter(3, 5) // Rate limit of 3 requests per second with a burst of 5 requests

func init() {
	file, err := os.OpenFile("logfile.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.SetOutput(file)
	} else {
		// If unable to open the log file, log to standard output
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
	case "/main":
		servePage(w, "webPages/registration.html")
	case "/login":
		servePage(w, "webPages/login.html")
	case "/home":
		servePage(w, "webPages/home.html")
	case "/admin":
		servePage(w, "webPages/admin.html")
	case "/profile":
		servePage(w, "webPages/profile.html")
	default:
		http.NotFound(w, r)
		defer r.Body.Close()
	}
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

func main() {
	// http.HandleFunc("/main", pageHandler)
	http.HandleFunc("/login", pageHandler)
	http.HandleFunc("/home", mongodb.ServeProfile)
	http.HandleFunc("/admin", pageHandler)
	http.HandleFunc("/profile", pageHandler)
	http.HandleFunc("/api/signUp", mongodb.AddUser)
	http.HandleFunc("/api/login", mongodb.LoginHandler)
	http.HandleFunc("/api/test", testRequest)
	http.HandleFunc("/api/profile/add", mongodb.AddUserProfile)
	http.HandleFunc("/api/getAllUsers", mongodb.GetAllUsers)
	http.HandleFunc("/verify", mongodb.VerifyAccount)
	http.HandleFunc("/main", handleWithRateLimit(pageHandler))
	http.Handle("/", http.FileServer(http.Dir("webPages/")))
	fmt.Println("Server is running on http://localhost:8080/main")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Error(err)
	}
}
