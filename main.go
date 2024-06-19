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
    fmt.Fprintf(w, "Welcome to NoteApp")
}

func setupRoutes() *mux.Router {
    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    return r
}

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    r := setupToJsonRoutes()
    fmt.Printf("Starting server on port %s\n", port)
    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), r))
}