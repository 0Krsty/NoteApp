package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
)

type Note struct {
    ID      string `json:"id"`
    Title   string `json:"title"`
    Content string `json:"content"`
}

var notes []Note

func loadEnv() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func getNotes(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    pageStr := r.URL.Query().Get("page")
    limitStr := r.URL.Query().Get("limit")

    page, err := strconv.Atoi(pageStr)
    if err != nil || page < 1 {
        page = 1 // default to first page
    }

    limit, err := strconv.Atoi(limitStr)
    if err != nil || limit < 1 || limit > 100 { // Limit results to 100 to prevent excessive loads
        limit = 10 // default to 10 items per page
    }

    startIndex := (page - 1) * limit
    endIndex := startIndex + limit

    if startIndex > len(notes) {
        startIndex = len(notes)
    }
    if endIndex > len(notes) {
        endIndex = len(notes)
    }

    json.NewEncoder(w).Encode(notes[startIndex:endIndex])
}

func createNote(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var note Note
    err := json.NewDecoder(r.Body).Decode(&note)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    notes = append(notes, note)
    json.NewEncoder(w).Encode(note)
}

func createNotesInBatch(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var newNotes []Note
    err := json.NewDecoder(r.Body).Decode(&newNotes)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    notes = append(notes, newNotes...)
    json.NewEncoder(w).Encode(newNotes)
}

func getNote(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    for _, item := range notes {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
    http.NotFound(w, r)
}

func updateNote(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    for index, item := range notes {
        if item.ID == params["id"] {
            notes = append(notes[:index], notes[index+1:]...)
            var note Note
            err := json.NewDecoder(r.Body).Decode(&note)
            if err != nil {
                http.Error(w, err.Error(), http.StatusBadRequest)
                return
            }
            note.ID = params["id"]
            notes = append(notes, note)
            json.NewEncoder(w).Encode(note)
            return
        }
    }
    http.NotFound(w, r)
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    for index, item := range notes {
        if item.ID == params["id"] {
            notes = append(notes[:index], notes[index+1:]...)
            json.NewEncoder(w).Encode(notes)
            return
        }
    }
    http.NotFound(w, r)
}

func deleteNotesInBatch(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var idsToDelete map[string]bool
    err := json.NewDecoder(r.Body).Decode(&idsToDelete)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    filteredNotes := []Note{}
    for _, note := range notes {
        if !idsToDelete[note.ID] {
            filteredNotes = append(filteredNotes, note)
        }
    }
    notes = filteredNotes
    json.NewEncoder(w).Encode(notes)
}

func main() {
    loadEv()

    router := mux.NewRouter()

    router.HandleFunc("/api/notes", getNotes).Methods("GET")
    router.HandleFunc("/api/notes", createNote).Methods("POST")
    router.HandleFunc("/api/notes/batch", createNotesInBatch).Methods("POST")
    router.HandleFunc("/api/notes/{id}", getNote).Methods("GET")
    router.HandleFunc("/api/notes/{id}", updateNote).Methods("PUT")
    router.HandleFunc("/api/notes/{id}", deleteNote).Methods("DELETE")
    router.HandleFunc("/api/notes/batch", deleteNotesInBatch).Methods("DELETE")

    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router))
}