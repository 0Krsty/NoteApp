package main

import (
    "fmt"
    "log"
    "net/http"
    "os"

    "github.com/joho/godotenv"
    "github.com/gorilla/mux"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    if _, err := fmt.Fprintf(w, "Welcome to NoteApp"); err != nil {
        log.Printf("Error writing response: %v", err)
    }
}

func setupRoutes() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    return r
}

func main() {
    if err := godotenv.Load(); err != nil {
        log.Printf("Error loading .env file: %v\n", err)
    }

    port := os.Getenv("PORT")
    if port == "" {
        log.Println("Port not specified in .env file, using the default port: 8080")
        port = "8080"
    }

    r := setupRoutes()
    serverAddress := fmt.Sprintf(":%s", port)

    fmt.Printf("Starting server on port %s\n", port)
    if err := http.ListenAndServe(serverAddress, r); err != nil {
        log.Fatalf("Failed to start the server on port %s: %v", port, err)
    }
}