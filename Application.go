package main

import (
    "fmt"
    "log"
    "bytes"
    "net/http"
    "./controllers"
    "gopkg.in/mgo.v2"
    "github.com/gorilla/mux"
)

// Application
type Application struct {
    Router      *mux.Router
    Connection  *mgo.Session
    Address     string
}

// Database configuration object
type DatabaseConfiguration struct {
    Username    string      `yaml:"mongo_user"`
    Password    string      `yaml:"mongo_pass"`
    Host        string      `yaml:"mongo_host"`
    Name        string      `yaml:"mongo_name"`
}

// Server configuration object
type ServerConfiguration struct {
    Host        string      `yaml:"api_host"`
    Port        int16       `yaml:"api_port"`
}

// Global configuration object
type GlobalConfiguration struct {
    DbConfig    DatabaseConfiguration   `yaml:"database"`
    SrvConfig   ServerConfiguration     `yaml:"api"`
}

// Initialize application object
func (app *Application) Initialize(config GlobalConfiguration) {
    app.InitializeDBConnection(config.DbConfig)
    app.InitializeServerAddr(config.SrvConfig)
    app.InitializeRouter()
}

// Initialize connection with MongoDB server
func (app *Application) InitializeDBConnection(config DatabaseConfiguration) {
    url := bytes.Buffer{}
    url.WriteString("mongodb://")

    if config.Username != "" && config.Password != "" {
        url.WriteString(config.Username)
        url.WriteString(":")
        url.WriteString(config.Password)
        url.WriteString("@")
    }
    if config.Host != "" {
        url.WriteString(config.Host)
    }
    if config.Name != "" {
        url.WriteString("/")
        url.WriteString(config.Name)
    }

    var err error
    app.Connection, err = mgo.Dial(url.String())
    if err != nil {
        log.Fatal(err)
    }
    log.Println("connection with mongodb established")
}

// Initialize server address
func (app *Application) InitializeServerAddr(config ServerConfiguration) {
    app.Address = fmt.Sprintf("%s:%d", config.Host, config.Port)   
}

// Initialize routes
func (app *Application) InitializeRouter() {
    app.Router = mux.NewRouter()
    app.Router.HandleFunc("/api/register", app.UserRegisterWrapper).Methods("POST")
    app.Router.HandleFunc("/api/login", app.UserLoginWrapper).Methods("POST")

    app.Router.HandleFunc("/api/room/create", app.RoomCreateWrapper).Methods("POST")
    app.Router.HandleFunc("/api/room/join", app.RoomJoinWrapper).Methods("PUT")
    app.Router.HandleFunc("/api/room/leave", app.RoomLeaveWrapper).Methods("PUT")
    app.Router.HandleFunc("/api/room/delete", app.RoomDeleteWrapper).Methods("PUT")
    app.Router.HandleFunc("/api/room/list", app.RoomListWrapper).Methods("GET")
    app.Router.HandleFunc("/api/room/users", app.RoomListUsersWrapper).Methods("POST")

    log.Println("router initialized")
}

// Run application
func (app *Application) Run() {
    log.Println("starting server ...")
    log.Fatal(http.ListenAndServe(app.Address, app.Router))
}

// Wrappers
func (app *Application) UserRegisterWrapper(writer http.ResponseWriter, request *http.Request) {
    controllers.RegisterAction(app.Connection, writer, request)
}

func (app *Application) UserLoginWrapper(writer http.ResponseWriter, request *http.Request) {
    controllers.LoginAction(app.Connection, writer, request)
}

func (app *Application) RoomCreateWrapper(writer http.ResponseWriter, request *http.Request) {
    controllers.RoomCreateAction(app.Connection, writer, request)
}

func (app *Application) RoomJoinWrapper(writer http.ResponseWriter, request *http.Request) {
    controllers.RoomJoinAction(app.Connection, writer, request)
}

func (app *Application) RoomLeaveWrapper(writer http.ResponseWriter, request *http.Request) {
    controllers.RoomLeaveAction(app.Connection, writer, request)
}

func (app *Application) RoomDeleteWrapper(writer http.ResponseWriter, request *http.Request) {
    controllers.RoomDeleteAction(app.Connection, writer, request)
}

func (app *Application) RoomListWrapper(writer http.ResponseWriter, request *http.Request) {
    controllers.RoomListAction(app.Connection, writer, request)
}

func (app *Application) RoomListUsersWrapper(writer http.ResponseWriter, request *http.Request) {
    controllers.RoomListUsersAction(app.Connection, writer, request)
}