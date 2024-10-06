package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func TestCreateUser(t *testing.T) {
	db = initTestDB()
	defer db.Close()

	requestBody := []byte(`{"username":"test@example.com","password":"testpassword"}`)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users", createUser).Methods("POST")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response User
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Username != "test@example.com" {
		t.Errorf("handler returned unexpected username: got %v want %v", response.Username, "test@example.com")
	}
}

func TestGetUsers(t *testing.T) {
	db = initTestDB()
	defer db.Close()

	req, err := http.NewRequest("GET", "/users", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users", getUsers).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response []User
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if len(response) != 0 {
		t.Errorf("handler returned unexpected number of users: got %v want %v", len(response), 0)
	}
}

func TestCreateWorkspace(t *testing.T) {
	db = initTestDB()
	defer db.Close()

	requestBody := []byte(`{"name":"testworkspace","user_id":1}`)
	req, err := http.NewRequest("POST", "/workspaces", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/workspaces", createWorkspace).Methods("POST")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response Workspace
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Name != "testworkspace" {
		t.Errorf("handler returned unexpected workspace name: got %v want %v", response.Name, "testworkspace")
	}
}

func TestGetWorkspaces(t *testing.T) {
	db = initTestDB()
	defer db.Close()

	req, err := http.NewRequest("GET", "/workspaces", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/workspaces", getWorkspaces).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response []Workspace
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if len(response) != 0 {
		t.Errorf("handler returned unexpected number of workspaces: got %v want %v", len(response), 0)
	}
}

func TestCreateApp(t *testing.T) {
	db = initTestDB()
	defer db.Close()

	requestBody := []byte(`{"name":"testapp","description":"Test app","git_hash":"abcdef","ip_port":"10.0.0.1:8080","endpoint":"/api","version":"1.0","workspace_id":1}`)
	req, err := http.NewRequest("POST", "/apps", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/apps", createApp).Methods("POST")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response App
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Name != "testapp" {
		t.Errorf("handler returned unexpected app name: got %v want %v", response.Name, "testapp")
	}
}

func TestGetApps(t *testing.T) {
	db = initTestDB()
	defer db.Close()

	req, err := http.NewRequest("GET", "/apps", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/apps", getApps).Methods("GET")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response []App
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if len(response) != 0 {
		t.Errorf("handler returned unexpected number of apps: got %v want %v", len(response), 0)
	}
}

func initTestDB() *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`
		CREATE TABLE users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
		CREATE TABLE workspaces (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			user_id INTEGER,
			subdomain TEXT NOT NULL UNIQUE,
			ips TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE apps (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			git_hash TEXT,
			ip_port TEXT NOT NULL,
			endpoint TEXT,
			version TEXT,
			workspace_id INTEGER,
			FOREIGN KEY(workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
