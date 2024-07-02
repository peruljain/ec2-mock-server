package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "strconv"
    "time"

    _ "github.com/go-sql-driver/mysql"
)

type Server struct {
    ID     int    `json:"id"`
    Status string `json:"status"`
}

var db *sql.DB

func main() {
    var err error
    dsn := "root:root_password@tcp(127.0.0.1:3306)/server"
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        log.Fatalf("Error opening database: %v", err)
    }
    defer db.Close()

    http.HandleFunc("/create", createServerHandler)
    http.HandleFunc("/status", getStatusHandler)

    log.Println("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func createServerHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    result, err := db.Exec("INSERT INTO server (status) VALUES ('IN_PROGRESS')")
    if err != nil {
        http.Error(w, fmt.Sprintf("Error creating server entry: %v", err), http.StatusInternalServerError)
        return
    }

    id, err := result.LastInsertId()
    if err != nil {
        http.Error(w, fmt.Sprintf("Error getting last insert ID: %v", err), http.StatusInternalServerError)
        return
    }

    go updateStatusAfterDelay(int(id))

    server := Server{ID: int(id), Status: "IN_PROGRESS"}
    json.NewEncoder(w).Encode(server)
}

func updateStatusAfterDelay(id int) {
    time.Sleep(20 * time.Second)
    _, err := db.Exec("UPDATE server SET status = 'CREATED' WHERE id = ?", id)
    if err != nil {
        log.Printf("Error updating server status: %v", err)
    }
}

func getStatusHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    idStr := r.URL.Query().Get("id")
    if idStr == "" {
        http.Error(w, "ID is required", http.StatusBadRequest)
        return
    }

    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid ID format", http.StatusBadRequest)
        return
    }

    var server Server
    err = db.QueryRow("SELECT id, status FROM server WHERE id = ?", id).Scan(&server.ID, &server.Status)
    if err != nil {
        if err == sql.ErrNoRows {
            http.Error(w, "Server not found", http.StatusNotFound)
        } else {
            http.Error(w, fmt.Sprintf("Error querying database: %v", err), http.StatusInternalServerError)
        }
        return
    }

    json.NewEncoder(w).Encode(server)
}
