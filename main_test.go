package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func TestDeleteUser(t *testing.T) {
	clearDatabase()
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

func clearDatabase() {
	db.Exec("DELETE FROM apps")
	db.Exec("DELETE FROM workspaces")
	db.Exec("DELETE FROM users")
}

func tearDownTestEnvironment() {
	if db != nil {
		db.Close()
	}
	os.Remove("./testdiscovery.db")
}

func TestCreateUser(t *testing.T) {
	clearDatabase()
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

func tearDownTestDB() {
	if db != nil {
		db.Close()
	}
	// os.Remove("./testdiscovery.db")
}

func TestCreateAppRole(t *testing.T) {
	clearDatabase()
	requestBody := []byte(`{"user_id":1,"role":"developer","app_id":1}`)
	req, err := http.NewRequest("POST", "/app-roles", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/app-roles", createAppRole).Methods("POST")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response AppRole
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Role != "developer" {
		t.Errorf("handler returned unexpected role: got %v want %v", response.Role, "developer")
	}
}

func TestDeleteApp(t *testing.T) {
	clearDatabase()
	// Create a test workspace first
	_, err := db.Exec("INSERT INTO workspaces (name, user_id, subdomain, ips) VALUES (?, ?, ?, ?)", "TestWorkspace", 1, "testsubdomain", "10.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	// Create a test app first
	app := App{Name: "TestApp", Description: "Test app for deletion", IPPort: "10.0.0.1:8080", WorkspaceID: 1}
	result, err := db.Exec("INSERT INTO apps (name, description, ip_port, workspace_id) VALUES (?, ?, ?, ?)",
		app.Name, app.Description, app.IPPort, app.WorkspaceID)
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

func TestCreateWorkspaceRole(t *testing.T) {
	clearDatabase()
	requestBody := []byte(`{"user_id":1,"role":"admin","workspace_id":1}`)
	req, err := http.NewRequest("POST", "/workspace-roles", bytes.NewBuffer(requestBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/workspace-roles", createWorkspaceRole).Methods("POST")
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response WorkspaceRole
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	if err != nil {
		t.Fatal(err)
	}

	if response.Role != "admin" {
		t.Errorf("handler returned unexpected role: got %v want %v", response.Role, "admin")
	}
}

func TestUpdateUser(t *testing.T) {
	clearDatabase()
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

func TestGetUsers(t *testing.T) {
	clearDatabase()
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

func TestGetWorkspaces(t *testing.T) {
	clearDatabase()
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

func setupTestDB() {
	os.Remove("./testdiscovery.db")
	var err error
	db, err = initDB("./testdiscovery.db")
	if err != nil {
		panic(err)
	}
}

func TestUpdateWorkspace(t *testing.T) {
	clearDatabase()
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

func TestCreateApp(t *testing.T) {
	clearDatabase()
	requestBody := []byte(`{"name":"testapp","description":"Test app","git_hash":"abcdef","ip_port":"10.0.0.1:8080","endpoint":"/api","version":"1.0","workspace_id":1,"input_schema":"test input","output_schema":"test output"}`)
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

	if response.InputSchema != "test input" {
		t.Errorf("handler returned unexpected input schema: got %v want %v", response.InputSchema, "test input")
	}

	if response.OutputSchema != "test output" {
		t.Errorf("handler returned unexpected output schema: got %v want %v", response.OutputSchema, "test output")
	}
}

func TestUpdateApp(t *testing.T) {
	clearDatabase()
	// Create a test workspace first
	_, err := db.Exec("INSERT INTO workspaces (name, user_id, subdomain, ips) VALUES (?, ?, ?, ?)", "TestWorkspace", 1, "testsubdomain", "10.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	// Create a test app first
	app := App{Name: "TestApp", Description: "Original description", IPPort: "10.0.0.1:8080", WorkspaceID: 1, InputSchema: "original input", OutputSchema: "original output"}
	result, err := db.Exec("INSERT INTO apps (name, description, ip_port, workspace_id, input_schema, output_schema) VALUES (?, ?, ?, ?, ?, ?)",
		app.Name, app.Description, app.IPPort, app.WorkspaceID, app.InputSchema, app.OutputSchema)
	if err != nil {
		t.Fatal(err)
	}
	appID, _ := result.LastInsertId()

	// Now update the app
	updatedApp := App{Name: "UpdatedTestApp", Description: "Updated description", IPPort: "10.0.0.2:8080", WorkspaceID: 1, InputSchema: "updated input", OutputSchema: "updated output"}
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
	err = db.QueryRow("SELECT name, description, ip_port, input_schema, output_schema FROM apps WHERE id = ?", appID).Scan(&updatedAppFromDB.Name, &updatedAppFromDB.Description, &updatedAppFromDB.IPPort, &updatedAppFromDB.InputSchema, &updatedAppFromDB.OutputSchema)
	if err != nil {
		t.Fatal(err)
	}
	if updatedAppFromDB.Name != updatedApp.Name || updatedAppFromDB.Description != updatedApp.Description || updatedAppFromDB.IPPort != updatedApp.IPPort || updatedAppFromDB.InputSchema != updatedApp.InputSchema || updatedAppFromDB.OutputSchema != updatedApp.OutputSchema {
		t.Errorf("app was not updated correctly: got %v, want %v", updatedAppFromDB, updatedApp)
	}
}

func TestGetApps(t *testing.T) {
	clearDatabase()
	// Insert a test app
	_, err := db.Exec("INSERT INTO apps (name, description, git_hash, ip_port, endpoint, version, workspace_id, input_schema, output_schema) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"testapp", "Test app", "abcdef", "10.0.0.1:8080", "/api", "1.0", 1, "test input", "test output")
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

	if response[0].InputSchema != "test input" {
		t.Errorf("handler returned unexpected input schema: got %v want %v", response[0].InputSchema, "test input")
	}

	if response[0].OutputSchema != "test output" {
		t.Errorf("handler returned unexpected output schema: got %v want %v", response[0].OutputSchema, "test output")
	}
}

func TestMain(m *testing.M) {
	// Set up
	setupTestDB()
	ipPool = NewIPPool()

	// Run tests
	code := m.Run()

	// Tear down
	tearDownTestDB()

	os.Exit(code)
}

func TestDeleteWorkspace(t *testing.T) {
	clearDatabase()
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
	clearDatabase()
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
