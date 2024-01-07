package main

import (
	"Aitunder/mongodb"
	"fmt"
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
	default:
		http.NotFound(w, r)
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

func main() {
	http.HandleFunc("/main", pageHandler)
	http.HandleFunc("/login", pageHandler)
	http.HandleFunc("/home", pageHandler)
	http.HandleFunc("/api/signUp", mongodb.AddUser)
	http.HandleFunc("/api/login", mongodb.LoginHandler)
	fmt.Println("Server is running on http://localhost:8080/main")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
