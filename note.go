package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    _ "github.com/lib/pq"
    "github.com/joho/godotenv"
)

type Note struct {
    ID      int64
    Title   string
    Content string
}

var db *sql.DB

func init() {
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
    db, err = sql.Open(
        "postgres",
        fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
            os.Getenv("DB_HOST"),
            os.Getenv("DB_PORT"),
            os.Getenv("DB_USER"),
            os.Getenv("DB_PASS"),
            os.Getenv("DB_NAME"),
        ),
    )
    if err != nil {
        log.Fatalf("Could not connect to database: %v", err)
    }
    err = db.Ping()
    if err != nil {
        log.Fatalf("Database connection test failed: %v", err)
    }
}

func CreateNote(title, content string) (*Note, error) {
    note := &Note{
        Title:   title,
        Content: content,
    }
    query := `INSERT INTO notes (title, content) VALUES ($1, $2) RETURNING id`
    err := db.QueryRow(query, note.Title, note.Content).Scan(&note.ID)
    if err != nil {
        return nil, fmt.Errorf("CreateNote failed: %v", err)
    }
    return note, nil
}

func GetNote(id int64) (*Note, error) {
    note := &Note{}
    query := `SELECT id, title, content FROM notes WHERE id = $1`
    row := db.QueryRow(query, id)
    err := row.Scan(&note.ID, &note.Title, &note.Content)
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, fmt.Errorf("GetNote No Result Found for ID: %d", id)
        }
        return nil, fmt.Errorf("GetNote failed: %v", err)
    }
    return note, nil
}

func UpdateNote(id int64, title, content string) (*Note, error) {
    note := &Note{
        ID:      id,
        Title:   title,
        Content: content,
    }
    query := `UPDATE notes SET title = $2, content = $3 WHERE id = $1`
    _, err := db.Exec(query, note.ID, note.Title, note.Content)
    if err != nil {
        return nil, fmt.Errorf("UpdateNote failed: %v", err)
    }
    return note, nil
}

func DeleteNote(id int64) error {
    query := `DELETE FROM notes WHERE id = $1`
    _, err := db.Exec(query, id)
    if err != nil {
        return fmt.Errorf("DeleteNote failed: %v", err)
    }
    return nil
}

func ListNotes() ([]*Note, error) {
    var notes []*Note
    query := `SELECT id, title, content FROM notes`
    rows, err := db.Query(query)
    if err != nil {
        return nil, fmt.Errorf("ListNotes failed: %v", err)
    }
    defer rows.Close()
    for rows.Next() {
        var note Note
        if err := rows.Scan(&note.ID, &note.Title, &note.Content); err != nil {
            return nil, fmt.Errorf("Failed to scan note: %v", err)
        }
        notes = append(notes, &note)
    }
    if err = rows.Err(); err != nil {
        return nil, fmt.Errorf("Rows iterating error: %v", err)
    }
    return notes, nil
}