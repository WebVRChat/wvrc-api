package controllers

import (
    "net/http"
    "strconv"
    "gopkg.in/mgo.v2"
    "../managers"
    "../utils"
)

func RoomCreateAction(connection *mgo.Session, writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm()
    if err != nil {
        panic(err)
    }

    vars := request.Form
    name := vars.Get("name")
    description := vars.Get("description")
    token := vars.Get("token")
    is_private, _ := strconv.ParseBool(vars.Get("is_private"))
    peer_id := vars.Get("peer_id")

    owner, err := managers.GetUserByToken(connection, token)
    if err != nil {
        utils.TriggerError(writer, 400, "unknown token")
        return
    }

    room, err := managers.RoomCreate(connection, owner, name, description, is_private)
    if err != nil {
        utils.TriggerError(writer, 400, err.Error())
        return
    }
    managers.RoomJoin(connection, owner.ToTmp(peer_id), &room)

    response := utils.ApplicationResponse{200, map[string]string{
        "status": "ok", "reason": "", "invite": room.Token}}

    response.SendResponse(writer)
}

func RoomJoinAction(connection *mgo.Session, writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm()
    if err != nil {
        panic(err)
    }

    vars := request.Form
    token := vars.Get("token")
    invite := vars.Get("invite")
    peer_id := vars.Get("peer_id")

    user, err := managers.GetUserByToken(connection, token)
    if err != nil {
        utils.TriggerError(writer, 400, "unknown token")
        return
    }

    room, err := managers.GetRoomByToken(connection, invite)
    if err != nil {
        utils.TriggerError(writer, 400, "unknown invite")
        return        
    }

    if (!managers.IsRoomUser(connection, user, room)) {
        managers.RoomJoin(connection, user.ToTmp(peer_id), &room)
    } else {
        utils.TriggerError(writer, 400, "you can't join this room")
        return
    }

    response := utils.ApplicationResponse{200, map[string]string{
        "status": "ok", "reason": "", "invite": room.Token}}
    response.SendResponse(writer)
}

func RoomLeaveAction(connection *mgo.Session, writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm()
    if err != nil {
        panic(err)
    }

    vars := request.Form
    token := vars.Get("token")
    invite := vars.Get("invite")

    user, err := managers.GetUserByToken(connection, token)
    if err != nil {
        utils.TriggerError(writer, 400, "unknown token")
        return
    }

    room, err := managers.GetRoomByToken(connection, invite)
    if err != nil {
        utils.TriggerError(writer, 400, "unknown invite")
        return        
    }

    if (managers.IsRoomUser(connection, user, room) && !managers.IsRoomOwner(connection, user, room)) {
        managers.RoomLeave(connection, user.ToTmp(""), &room)  
    } else {
        utils.TriggerError(writer, 400, "you can't leave this room")
        return
    }

    response := utils.ApplicationResponse{200, map[string]string{
        "status": "ok", "reason": "", "invite": room.Token}}
    response.SendResponse(writer)
}

func RoomDeleteAction(connection *mgo.Session, writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm()
    if err != nil {
        panic(err)
    }

    vars := request.Form
    token := vars.Get("token")
    invite := vars.Get("invite")

    user, err := managers.GetUserByToken(connection, token)
    if err != nil {
        utils.TriggerError(writer, 400, "unkown token")
        return
    }

    room, err := managers.GetRoomByToken(connection, invite)
    if err != nil {
        utils.TriggerError(writer, 400, "unknown invite")
        return
    }

    if (managers.IsRoomOwner(connection, user, room)) {
        managers.RoomDelete(connection, room)
    } else {
        utils.TriggerError(writer, 400, "specified user is not owner")
        return
    }

    response := utils.ApplicationResponse{200, map[string]string{
        "status": "ok", "reason": ""}}
    response.SendResponse(writer)
}

func RoomListAction(connection *mgo.Session, writer http.ResponseWriter, request *http.Request) {
    list, err := managers.GetRoomList(connection)
    if err != nil {
        utils.TriggerError(writer, 400, "can't list rooms")
        return
    }

    payload := make(map[string]map[string]string)
    for _, room := range list {
        if !room.IsPrivate {
            payload[room.Name] = map[string]string{
                "description": room.Description,
                "invite": room.Token,
            }
        }
    }

    response := utils.ApplicationResponse{200, map[string]interface{}{
        "status": "ok", "reason": "", "rooms": payload}}
    response.SendResponse(writer)
}

func RoomListUsersAction(connection *mgo.Session, writer http.ResponseWriter, request *http.Request) {
    err := request.ParseForm()
    if err != nil {
        panic(err)
    }

    vars := request.Form
    token := vars.Get("token")
    invite := vars.Get("invite")

    user, err := managers.GetUserByToken(connection, token)
    if err != nil {
        utils.TriggerError(writer, 400, "unkown token")
        return
    }

    room, err := managers.GetRoomByToken(connection, invite)
    if err != nil {
        utils.TriggerError(writer, 400, "unknown invite")
        return
    }

    if !managers.IsRoomUser(connection, user, room) {
        utils.TriggerError(writer, 400, "can't list users")
        return
    }

    payload := make(map[string]string)
    for _, user := range room.Users {
        payload[user.Username] = user.PeerId
    }

    response := utils.ApplicationResponse{200, map[string]interface{}{
        "status": "ok", "reason": "", "users": payload}}
    response.SendResponse(writer)
}