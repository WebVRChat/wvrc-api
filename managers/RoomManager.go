package managers

import (
    "errors"
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/mgo.v2"
    "../models"
    "../utils"
    "log"
)

func RoomCreate(connection *mgo.Session, owner models.User, name, description string, is_private bool) (models.Room, error){
    column := connection.DB("").C("rooms")

    room := models.Room{
        OwnerId: owner.Id,
        Token: utils.GenerateToken(12),
        Name: name,
        Description: description,
        IsPrivate: is_private,
    }
    room.Users = make([]models.TmpUser, 0)

    err := column.Insert(&room)
    if err != nil {
        log.Println(err.Error())
        return models.Room{}, errors.New("database error")
    }

    return room, nil
}

func RoomJoin(connection *mgo.Session, user models.TmpUser, room *models.Room) (*models.Room, error) {
    column := connection.DB("").C("rooms")

    room.Users = append(room.Users, user)
    err := column.Update(bson.M{
        "token": room.Token,
    }, bson.M{"$set": bson.M{
        "users": room.Users,
    }})
    if err != nil {
        log.Println(err.Error())
        return &models.Room{}, errors.New("database error")
    }

    return room, nil
}

func RoomLeave(connection *mgo.Session, user models.TmpUser, room *models.Room) (*models.Room, error) {
    column := connection.DB("").C("rooms")

    for i := 0; i < len(room.Users); i++ {
        tmpUser := room.Users[i]
        if tmpUser.Username == user.Username {
            room.Users = append(room.Users[:i], room.Users[i+1:]...)
            i--
            break
        }
    }
    err := column.Update(bson.M{
        "token": room.Token,
    }, bson.M{"$set": bson.M{
        "users": room.Users,
    }})
    if err != nil {
        return &models.Room{}, errors.New("database error")
    }

    return room, nil
}

func RoomDelete(connection *mgo.Session, room models.Room) (models.Room, error) {
    column := connection.DB("").C("rooms")

    err := column.Remove(bson.M{
        "token": room.Token,
    })
    if err != nil {
        return models.Room{}, errors.New("database error")
    }

    return room, nil
}

func IsRoomOwner(connection *mgo.Session, user models.User, room models.Room) bool {
    if room.OwnerId == user.Id {
        return true
    }
    return false
}

func IsRoomUser(connection *mgo.Session, user models.User, room models.Room) bool {
    for _, cpy := range room.Users {
        if cpy.Username == user.Username {
            return true
        }
    }
    return false
}

func GetRoomByToken(connection *mgo.Session, token string) (models.Room, error) {
    column := connection.DB("").C("rooms")

    var tmpRoom models.Room
    column.Find(bson.M{
        "token": token,
    }).One(&tmpRoom)

    if tmpRoom.Token == "" {
        return models.Room{}, errors.New("room not found")
    }

    return tmpRoom, nil
}

func GetRoomList(connection *mgo.Session) ([]models.Room, error) {
    column := connection.DB("").C("rooms")

    var rooms []models.Room
    err := column.Find(nil).All(&rooms)
    if err != nil {
        return nil, errors.New("database error")
    }

    return rooms, nil
}