package controllers

import (
    "net/http"
    "gopkg.in/mgo.v2"
    "../managers"
    "../utils"
)

func RegisterAction(connection *mgo.Session, writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm()
    if err != nil {
        panic(err)
    }

    vars := request.Form
    username := vars.Get("username")
    password := vars.Get("password")

    err = managers.UserRegister(connection, username, password)
    response := utils.ApplicationResponse{}

    if err != nil {
        response.Code = 400
        response.Data = map[string]string{"status":"fail", "reason": err.Error()}
    } else {
        response.Code = 200
        response.Data = map[string]string{"status":"ok", "reason": ""}
    }
    
    response.SendResponse(writer)
}

func LoginAction(connection *mgo.Session, writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm()
    if err != nil {
        panic(err)
    }

    vars := request.Form
    username := vars.Get("username")
    password := vars.Get("password")

    user, err := managers.UserLogin(connection, username, password)
    response := utils.ApplicationResponse{}

    if err != nil {
        response.Code = 400
        response.Data = map[string]string{"status":"fail", "reason": err.Error()}
    } else {
        response.Code = 200
        response.Data = map[string]string{"status":"ok", "reason":"", "token": user.Token}
    }

    response.SendResponse(writer)
}