package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"net/mail"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"-"` // Password is never sent in JSON responses
}

type Workspace struct {
	ID        int      `json:"id"`
	Name      string   `json:"name"`
	UserID    int      `json:"user_id"`
	Subdomain string   `json:"subdomain"`
	IPs       []string `json:"ips"`
}

type App struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	GitHash      string `json:"git_hash"`
	IPPort       string `json:"ip_port"`
	Endpoint     string `json:"endpoint"`
	Version      string `json:"version"`
	WorkspaceID  int    `json:"workspace_id"`
	InputSchema  string `json:"input_schema"`
	OutputSchema string `json:"output_schema"`
}

type WorkspaceRole struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Role        string `json:"role"`
	WorkspaceID int    `json:"workspace_id"`
}

type AppRole struct {
	ID     int    `json:"id"`
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	AppID  int    `json:"app_id"`
}

var (
	db         *sql.DB
	ipPool     *IPPool
	mutex      sync.Mutex
	subdomains map[string]bool
)

type IPPool struct {
	available []net.IP
	inUse     map[string]bool
	mutex     sync.Mutex
}

func createWorkspaceRole(w http.ResponseWriter, r *http.Request) {
	var role WorkspaceRole
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO workspace_roles (user_id, role, workspace_id) VALUES (?, ?, ?)",
		role.UserID, role.Role, role.WorkspaceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	role.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

func deleteApp(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, err := db.Exec("DELETE FROM apps WHERE id = ?", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func init() {
	rand.Seed(time.Now().UnixNano())
	subdomains = make(map[string]bool)
}

func NewIPPool() *IPPool {
	pool := &IPPool{
		available: make([]net.IP, 0),
		inUse:     make(map[string]bool),
	}

	// Generate IPs for 10.0.0.0/16 and 172.16.0.0/16
	for i := 0; i <= 255; i++ {
		for j := 0; j <= 255; j++ {
			pool.available = append(pool.available, net.IPv4(10, 0, byte(i), byte(j)))
			pool.available = append(pool.available, net.IPv4(172, 16, byte(i), byte(j)))
		}
	}

	return pool
}

func (p *IPPool) ReleaseIP(ip string) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.inUse[ip] {
		delete(p.inUse, ip)
		p.available = append(p.available, net.ParseIP(ip))
	}
}

func updateApp(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var app App
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE apps SET name = ?, description = ?, git_hash = ?, ip_port = ?, endpoint = ?, version = ?, workspace_id = ?, input_schema = ?, output_schema = ? WHERE id = ?",
		app.Name, app.Description, app.GitHash, app.IPPort, app.Endpoint, app.Version, app.WorkspaceID, app.InputSchema, app.OutputSchema, params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(app)
}

func updateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate email
	_, err := mail.ParseAddress(user.Username)
	if err != nil {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE users SET username = ?, password = ? WHERE id = ?", user.Username, string(hashedPassword), params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = "" // Don't send password back
	json.NewEncoder(w).Encode(user)
}

func getApps(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, description, git_hash, ip_port, endpoint, version, workspace_id, input_schema, output_schema FROM apps")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	apps := []App{}
	for rows.Next() {
		var a App
		if err := rows.Scan(&a.ID, &a.Name, &a.Description, &a.GitHash, &a.IPPort, &a.Endpoint, &a.Version, &a.WorkspaceID, &a.InputSchema, &a.OutputSchema); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		apps = append(apps, a)
	}

	json.NewEncoder(w).Encode(apps)
}

func createApp(w http.ResponseWriter, r *http.Request) {
	var app App
	if err := json.NewDecoder(r.Body).Decode(&app); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO apps (name, description, git_hash, ip_port, endpoint, version, workspace_id, input_schema, output_schema) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)",
		app.Name, app.Description, app.GitHash, app.IPPort, app.Endpoint, app.Version, app.WorkspaceID, app.InputSchema, app.OutputSchema)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	app.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(app)
}

func createAppRole(w http.ResponseWriter, r *http.Request) {
	var role AppRole
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result, err := db.Exec("INSERT INTO app_roles (user_id, role, app_id) VALUES (?, ?, ?)",
		role.UserID, role.Role, role.AppID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	role.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(role)
}

func deleteWorkspaceRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, err := db.Exec("DELETE FROM workspace_roles WHERE id = ?", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func deleteUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, err := db.Exec("DELETE FROM users WHERE id = ?", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func createUser(w http.ResponseWriter, r *http.Request) {
	var user User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate email
	_, err := mail.ParseAddress(user.Username)
	if err != nil {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	user.ID = int(id)

	// Create default workspace for the user
	subdomain := generateSubdomain()
	ip, err := ipPool.AllocateIP()
	if err != nil {
		http.Error(w, "Failed to allocate IP", http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO workspaces (name, user_id, subdomain, ips) VALUES (?, ?, ?, ?)",
		"default", user.ID, subdomain, ip)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	user.Password = "" // Don't send password back
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func getWorkspace(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var workspace Workspace
	var ips string
	err := db.QueryRow("SELECT id, name, user_id, subdomain, ips FROM workspaces WHERE id = ?", params["id"]).
		Scan(&workspace.ID, &workspace.Name, &workspace.UserID, &workspace.Subdomain, &ips)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	workspace.IPs = strings.Split(ips, ",")
	json.NewEncoder(w).Encode(workspace)
}

func getUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var user User
	err := db.QueryRow("SELECT id, username FROM users WHERE id = ?", params["id"]).Scan(&user.ID, &user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func (p *IPPool) AllocateIP() (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.available) == 0 {
		return "", fmt.Errorf("no available IPs")
	}

	ip := p.available[0]
	p.available = p.available[1:]
	ipStr := ip.String()
	p.inUse[ipStr] = true

	return ipStr, nil
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, username FROM users")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Username); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, u)
	}

	json.NewEncoder(w).Encode(users)
}

func updateWorkspace(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var workspace Workspace
	if err := json.NewDecoder(r.Body).Decode(&workspace); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE workspaces SET name = ?, user_id = ? WHERE id = ?",
		workspace.Name, workspace.UserID, params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(workspace)
}

func updateAppRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var role AppRole
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE app_roles SET user_id = ?, role = ?, app_id = ? WHERE id = ?",
		role.UserID, role.Role, role.AppID, params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(role)
}

func initDB(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS workspaces (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			user_id INTEGER,
			subdomain TEXT NOT NULL UNIQUE,
			ips TEXT NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS apps (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			description TEXT,
			git_hash TEXT,
			ip_port TEXT NOT NULL,
			endpoint TEXT,
			version TEXT,
			workspace_id INTEGER,
			input_schema TEXT,
			output_schema TEXT,
			FOREIGN KEY(workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS workspace_roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			role TEXT NOT NULL,
			workspace_id INTEGER,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(workspace_id) REFERENCES workspaces(id) ON DELETE CASCADE
		);
		CREATE TABLE IF NOT EXISTS app_roles (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			role TEXT NOT NULL,
			app_id INTEGER,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(app_id) REFERENCES apps(id) ON DELETE CASCADE
		);
	`)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getWorkspaces(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, name, user_id, subdomain, ips FROM workspaces")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	workspaces := []Workspace{}
	for rows.Next() {
		var ws Workspace
		var ips string
		if err := rows.Scan(&ws.ID, &ws.Name, &ws.UserID, &ws.Subdomain, &ips); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ws.IPs = strings.Split(ips, ",")
		workspaces = append(workspaces, ws)
	}

	json.NewEncoder(w).Encode(workspaces)
}

func deleteAppRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	_, err := db.Exec("DELETE FROM app_roles WHERE id = ?", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func deleteWorkspace(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	// Get the workspace's IPs before deleting
	var ips string
	err := db.QueryRow("SELECT ips FROM workspaces WHERE id = ?", params["id"]).Scan(&ips)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Delete the workspace
	_, err = db.Exec("DELETE FROM workspaces WHERE id = ?", params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Release the IPs
	for _, ip := range strings.Split(ips, ",") {
		ipPool.ReleaseIP(ip)
	}

	w.WriteHeader(http.StatusNoContent)
}

func getWorkspaceRoles(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, user_id, role, workspace_id FROM workspace_roles")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	roles := []WorkspaceRole{}
	for rows.Next() {
		var role WorkspaceRole
		if err := rows.Scan(&role.ID, &role.UserID, &role.Role, &role.WorkspaceID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		roles = append(roles, role)
	}

	json.NewEncoder(w).Encode(roles)
}

func getAppRoles(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, user_id, role, app_id FROM app_roles")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	roles := []AppRole{}
	for rows.Next() {
		var role AppRole
		if err := rows.Scan(&role.ID, &role.UserID, &role.Role, &role.AppID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		roles = append(roles, role)
	}

	json.NewEncoder(w).Encode(roles)
}

func generateSubdomain() string {
	mutex.Lock()
	defer mutex.Unlock()

	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	for {
		result := make([]byte, 8)
		for i := range result {
			result[i] = chars[rand.Intn(len(chars))]
		}
		subdomain := string(result)
		if !subdomains[subdomain] {
			subdomains[subdomain] = true
			return subdomain
		}
	}
}

func updateWorkspaceRole(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var role WorkspaceRole
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE workspace_roles SET user_id = ?, role = ?, workspace_id = ? WHERE id = ?",
		role.UserID, role.Role, role.WorkspaceID, params["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(role)
}

func main() {
	initDB("./discovery.db")
	defer db.Close()

	ipPool = NewIPPool()

	r := mux.NewRouter()

	// User routes
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users", getUsers).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", getUser).Methods("GET")
	r.HandleFunc("/users/{id:[0-9]+}", updateUser).Methods("PUT")
	r.HandleFunc("/users/{id:[0-9]+}", deleteUser).Methods("DELETE")

	// Workspace routes
	r.HandleFunc("/workspaces", createWorkspace).Methods("POST")
	r.HandleFunc("/workspaces", getWorkspaces).Methods("GET")
	r.HandleFunc("/workspaces/{id:[0-9]+}", getWorkspace).Methods("GET")
	r.HandleFunc("/workspaces/{id:[0-9]+}", updateWorkspace).Methods("PUT")
	r.HandleFunc("/workspaces/{id:[0-9]+}", deleteWorkspace).Methods("DELETE")

	// App routes
	r.HandleFunc("/apps", createApp).Methods("POST")
	r.HandleFunc("/apps", getApps).Methods("GET")
	r.HandleFunc("/apps/{id:[0-9]+}", getApp).Methods("GET")
	r.HandleFunc("/apps/{id:[0-9]+}", updateApp).Methods("PUT")
	r.HandleFunc("/apps/{id:[0-9]+}", deleteApp).Methods("DELETE")

	// Workspace role routes
	r.HandleFunc("/workspace-roles", createWorkspaceRole).Methods("POST")
	r.HandleFunc("/workspace-roles", getWorkspaceRoles).Methods("GET")
	r.HandleFunc("/workspace-roles/{id:[0-9]+}", updateWorkspaceRole).Methods("PUT")
	r.HandleFunc("/workspace-roles/{id:[0-9]+}", deleteWorkspaceRole).Methods("DELETE")

	// App role routes
	r.HandleFunc("/app-roles", createAppRole).Methods("POST")
	r.HandleFunc("/app-roles", getAppRoles).Methods("GET")
	r.HandleFunc("/app-roles/{id:[0-9]+}", updateAppRole).Methods("PUT")
	r.HandleFunc("/app-roles/{id:[0-9]+}", deleteAppRole).Methods("DELETE")

	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func getApp(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	var app App
	err := db.QueryRow("SELECT id, name, description, git_hash, ip_port, endpoint, version, workspace_id, input_schema, output_schema FROM apps WHERE id = ?", params["id"]).
		Scan(&app.ID, &app.Name, &app.Description, &app.GitHash, &app.IPPort, &app.Endpoint, &app.Version, &app.WorkspaceID, &app.InputSchema, &app.OutputSchema)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(app)
}

func createWorkspace(w http.ResponseWriter, r *http.Request) {
	var workspace Workspace
	if err := json.NewDecoder(r.Body).Decode(&workspace); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	workspace.Subdomain = generateSubdomain()
	ip, err := ipPool.AllocateIP()
	if err != nil {
		http.Error(w, "Failed to allocate IP", http.StatusInternalServerError)
		return
	}
	workspace.IPs = []string{ip}

	result, err := db.Exec("INSERT INTO workspaces (name, user_id, subdomain, ips) VALUES (?, ?, ?, ?)",
		workspace.Name, workspace.UserID, workspace.Subdomain, strings.Join(workspace.IPs, ","))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	workspace.ID = int(id)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workspace)
}
