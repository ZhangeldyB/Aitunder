package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func pageHandler(w http.ResponseWriter, r *http.Request) {
	htmlContent, err := os.ReadFile("login.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(htmlContent)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		http.Error(w, "Error unmarshling request body", http.StatusBadRequest)
		fmt.Print(err.Error())
	}

	fmt.Println(data)
}

func main() {
	http.HandleFunc("/main", pageHandler)
	http.HandleFunc("/login", loginHandler)
	fmt.Println("Server is running on http://localhost:8080/main")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}
}
