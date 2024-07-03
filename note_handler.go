package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"

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
    json.NewEncoder(w).Encode(notes)
}

func createNote(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var note Note
    _ = json.NewDecoder(r.Body).Decode(&note)
    notes = append(notes, note)
    json.NewEncoder(w).Encode(note)
}

// New endpoint for creating notes in batch
func createNotesInBatch(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var newNotes []Note
    _ = json.NewDecoder(r.Body).Decode(&newNotes)
    notes = append(notes, newNotes...)
    json.NewEncoder(w).Encode(newNotes)
}

func getNote(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r) // Fixed typo from Vors to Vars
    for _, item := range notes {
        if item.ID == params["id"] {
            json.NewEncoder(w).Encode(item)
            return
        }
    }
}

func updateNote(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    for index, item := range notes {
        if item.ID == params["id"] {
            notes = append(notes[:index], notes[index+1:]...)
            var note Note
            _ = json.NewDecoder(r.Body).Decode(&note)
            note.ID = params["id"]
            notes = append(notes, note)
            json.NewEncoder(w).Encode(note)
            return
        }
    }
}

func deleteNote(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    params := mux.Vars(r)
    for index, item := range notes {
        if item.ID == params["id"] {
            notes = append(notes[:index], notes[index+1:]...)
            break
        }
    }
    json.NewEncoder(w).Encode(notes)
}

// New endpoint for deleting notes in batch by Ids
func deleteNotesInBatch(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    var idsToDelete map[string]bool
    _ = json.NewDecoder(r.Body).Decode(&idsToDelete) // Assuming request body is {"id1":true, "id2":true, ...}
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
    loadEnv()

    router := mux.NewRouter()

    router.HandleFunc("/api/notes", getNotes).Methods("GET")
    router.HandleFunc("/api/notes", createNote).Methods("POST")
    router.HandleFunc("/api/notes/batch", createNotesInBatch).Methods("POST") // Batch create
    router.HandleFunc("/api/notes/{id}", getNote).Methods("GET")
    router.HandleFunc("/api/notes/{id}", updateNote).Methods("PUT")
    router.HandleFunc("/api/notes/{id}", deleteNote).Methods("DELETE")
    router.HandleFunc("/api/notes/batch", deleteNotesInProjectionJobBatch).Methods("DELETE") // Batch delete

    log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), router))
}