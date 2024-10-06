package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateApp(t *testing.T) {
	setUpTestEnvironment()
	defer tearDownTestEnvironment()

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

func TestGetUsers(t *testing.T) {
	setUpTestEnvironment()
	defer tearDownTestEnvironment()

	// Insert a test user
	_, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", "testuser@example.com", "hashedpassword")
	if err != nil {
		t.Fatal(err)
	}

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

	if len(response) != 1 {
		t.Errorf("handler returned unexpected number of users: got %v want %v", len(response), 1)
	}
}

func setUpTestEnvironment() {
	db = initTestDB()
	ipPool = NewIPPool()
}

func TestDeleteUser(t *testing.T) {
	// Create a test user first
	user := User{Username: "testuser@example.com", Password: "testpassword"}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, string(hashedPassword))
	if err != nil {
		t.Fatal(err)
	}
	userID, _ := result.LastInsertId()

	// Now delete the user
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/users/%d", userID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id:[0-9]+}", deleteUser).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// Verify that the user was deleted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = ?", userID).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("user was not deleted: got %v records, want 0", count)
	}
}

func TestUpdateApp(t *testing.T) {
	// Create a test app first
	app := App{Name: "TestApp", Description: "Original description"}
	result, err := db.Exec("INSERT INTO apps (name, description) VALUES (?, ?)", app.Name, app.Description)
	if err != nil {
		t.Fatal(err)
	}
	appID, _ := result.LastInsertId()

	// Now update the app
	updatedApp := App{Name: "UpdatedTestApp", Description: "Updated description"}
	requestBody, _ := json.Marshal(updatedApp)
	req, err := http.NewRequest("PUT", fmt.Sprintf("/apps/%d", appID), bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/apps/{id:[0-9]+}", updateApp).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify that the app was updated
	var updatedAppFromDB App
	err = db.QueryRow("SELECT name, description FROM apps WHERE id = ?", appID).Scan(&updatedAppFromDB.Name, &updatedAppFromDB.Description)
	if err != nil {
		t.Fatal(err)
	}
	if updatedAppFromDB.Name != updatedApp.Name || updatedAppFromDB.Description != updatedApp.Description {
		t.Errorf("app was not updated correctly: got %v, want %v", updatedAppFromDB, updatedApp)
	}
}

func TestUpdateWorkspace(t *testing.T) {
	// Create a test workspace first
	workspace := Workspace{Name: "TestWorkspace", UserID: 1, Subdomain: "testsubdomain", IPs: []string{"10.0.0.1"}}
	result, err := db.Exec("INSERT INTO workspaces (name, user_id, subdomain, ips) VALUES (?, ?, ?, ?)",
		workspace.Name, workspace.UserID, workspace.Subdomain, strings.Join(workspace.IPs, ","))
	if err != nil {
		t.Fatal(err)
	}
	workspaceID, _ := result.LastInsertId()

	// Now update the workspace
	updatedWorkspace := Workspace{Name: "UpdatedTestWorkspace", UserID: 2}
	requestBody, _ := json.Marshal(updatedWorkspace)
	req, err := http.NewRequest("PUT", fmt.Sprintf("/workspaces/%d", workspaceID), bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/workspaces/{id:[0-9]+}", updateWorkspace).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify that the workspace was updated
	var updatedWorkspaceFromDB Workspace
	err = db.QueryRow("SELECT name, user_id FROM workspaces WHERE id = ?", workspaceID).Scan(&updatedWorkspaceFromDB.Name, &updatedWorkspaceFromDB.UserID)
	if err != nil {
		t.Fatal(err)
	}
	if updatedWorkspaceFromDB.Name != updatedWorkspace.Name || updatedWorkspaceFromDB.UserID != updatedWorkspace.UserID {
		t.Errorf("workspace was not updated correctly: got %v, want %v", updatedWorkspaceFromDB, updatedWorkspace)
	}
}

func TestCreateUser(t *testing.T) {
	setUpTestEnvironment()
	defer tearDownTestEnvironment()

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

func initTestDB() *sql.DB {
	os.Remove("./testdiscovery.db")

	db, err := sql.Open("sqlite3", "./testdiscovery.db")
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

func TestDeleteApp(t *testing.T) {
	// Create a test app first
	app := App{Name: "TestApp", Description: "Test app for deletion"}
	result, err := db.Exec("INSERT INTO apps (name, description) VALUES (?, ?)", app.Name, app.Description)
	if err != nil {
		t.Fatal(err)
	}
	appID, _ := result.LastInsertId()

	// Now delete the app
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/apps/%d", appID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/apps/{id:[0-9]+}", deleteApp).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// Verify that the app was deleted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM apps WHERE id = ?", appID).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("app was not deleted: got %v records, want 0", count)
	}
}

func TestDeleteWorkspace(t *testing.T) {
	// Create a test workspace first
	workspace := Workspace{Name: "TestWorkspace", UserID: 1, Subdomain: "testsubdomain", IPs: []string{"10.0.0.1"}}
	result, err := db.Exec("INSERT INTO workspaces (name, user_id, subdomain, ips) VALUES (?, ?, ?, ?)",
		workspace.Name, workspace.UserID, workspace.Subdomain, strings.Join(workspace.IPs, ","))
	if err != nil {
		t.Fatal(err)
	}
	workspaceID, _ := result.LastInsertId()

	// Now delete the workspace
	req, err := http.NewRequest("DELETE", fmt.Sprintf("/workspaces/%d", workspaceID), nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/workspaces/{id:[0-9]+}", deleteWorkspace).Methods("DELETE")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	// Verify that the workspace was deleted
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM workspaces WHERE id = ?", workspaceID).Scan(&count)
	if err != nil {
		t.Fatal(err)
	}
	if count != 0 {
		t.Errorf("workspace was not deleted: got %v records, want 0", count)
	}
}

func TestCreateWorkspace(t *testing.T) {
	setUpTestEnvironment()
	defer tearDownTestEnvironment()

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
	setUpTestEnvironment()
	defer tearDownTestEnvironment()

	// Insert test workspaces
	_, err := db.Exec("INSERT INTO workspaces (name, user_id, subdomain, ips) VALUES (?, ?, ?, ?)", "workspace1", 1, "subdomain1", "10.0.0.1")
	if err != nil {
		t.Fatal(err)
	}
	_, err = db.Exec("INSERT INTO workspaces (name, user_id, subdomain, ips) VALUES (?, ?, ?, ?)", "workspace2", 1, "subdomain2", "10.0.0.2")
	if err != nil {
		t.Fatal(err)
	}

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

	if len(response) != 2 {
		t.Errorf("handler returned unexpected number of workspaces: got %v want %v", len(response), 2)
	}
}

func TestMain(m *testing.M) {
	// Set up
	setUpTestEnvironment()

	// Run tests
	code := m.Run()

	// Tear down
	tearDownTestEnvironment()

	os.Exit(code)
}

func tearDownTestEnvironment() {
	if db != nil {
		db.Close()
	}
	os.Remove("./testdiscovery.db")
}

func TestGetApps(t *testing.T) {
	setUpTestEnvironment()
	defer tearDownTestEnvironment()

	// Insert a test app
	_, err := db.Exec("INSERT INTO apps (name, description, git_hash, ip_port, endpoint, version, workspace_id) VALUES (?, ?, ?, ?, ?, ?, ?)",
		"testapp", "Test app", "abcdef", "10.0.0.1:8080", "/api", "1.0", 1)
	if err != nil {
		t.Fatal(err)
	}

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

	if len(response) != 1 {
		t.Errorf("handler returned unexpected number of apps: got %v want %v", len(response), 1)
	}
}

func TestUpdateUser(t *testing.T) {
	// Create a test user first
	user := User{Username: "testuser@example.com", Password: "testpassword"}
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, string(hashedPassword))
	if err != nil {
		t.Fatal(err)
	}
	userID, _ := result.LastInsertId()

	// Now update the user
	updatedUser := User{Username: "updateduser@example.com", Password: "updatedpassword"}
	requestBody, _ := json.Marshal(updatedUser)
	req, err := http.NewRequest("PUT", fmt.Sprintf("/users/%d", userID), bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/users/{id:[0-9]+}", updateUser).Methods("PUT")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	// Verify that the user was updated
	var updatedUserFromDB User
	err = db.QueryRow("SELECT username FROM users WHERE id = ?", userID).Scan(&updatedUserFromDB.Username)
	if err != nil {
		t.Fatal(err)
	}
	if updatedUserFromDB.Username != updatedUser.Username {
		t.Errorf("user was not updated correctly: got %v, want %v", updatedUserFromDB.Username, updatedUser.Username)
	}
}
