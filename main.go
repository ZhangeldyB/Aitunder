package main

import (
	"Aitunder/mongodb"
	"fmt"
	"io"
	"net/http"
	"os"
)

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
	http.HandleFunc("/main", pageHandler)
	http.HandleFunc("/login", pageHandler)
	http.HandleFunc("/home", pageHandler)
	http.HandleFunc("/admin", pageHandler)
	http.HandleFunc("/profile", pageHandler)
	http.HandleFunc("/api/signUp", mongodb.AddUser)
	http.HandleFunc("/api/login", mongodb.LoginHandler)
	http.HandleFunc("/api/test", testRequest)
	http.HandleFunc("/api/profile/add", mongodb.AddUserProfile)
	http.HandleFunc("/api/getAllUsers", mongodb.GetAllUsers)
	fmt.Println("Server is running on http://localhost:8080/main")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
