package managers

import (
    "errors"
    "encoding/hex"
    "crypto/sha256"
    
    "gopkg.in/mgo.v2/bson"
    "gopkg.in/mgo.v2"
    
    "../models"
    "../utils"
)

func UserRegister(connection *mgo.Session, username, password string) error {
    if username == "" || password == "" {
        return errors.New("empty username or password")
    }

    column := connection.DB("").C("users")

    hash := sha256.New()
    hash.Write([]byte(password))
    checksum := hex.EncodeToString(hash.Sum(nil))

    user := models.User{
        Username: username,
        Password: checksum,
        Token: utils.GenerateToken(12),
    }

    var testUser models.User
    column.Find(bson.M{
            "username": username,
    }).One(&testUser)

    if testUser.Username != "" {
        return errors.New("user already exists")
    }

    err := column.Insert(user)
    if err != nil {
        return errors.New("database error")
    }

    return nil
}

func UserLogin(connection *mgo.Session, username, password string) (models.User, error) {
    if username == "" || password == "" {
        return models.User{}, errors.New("empty username or password")
    }

    column := connection.DB("").C("users")
    
    hash := sha256.New()
    hash.Write([]byte(password))
    checksum := hex.EncodeToString(hash.Sum(nil))

    var user models.User
    column.Find(bson.M{
        "username": username,
        "password": checksum,
    }).One(&user)

    if user.Username == "" {
        return models.User{}, errors.New("user failed to login")
    }

    return user, nil
}

func GetUserByToken(connection *mgo.Session, token string) (models.User, error) {
    if token == "" {
        return models.User{}, errors.New("empty token")
    }

    column := connection.DB("").C("users")

    var user models.User
    column.Find(bson.M{
        "token": token,
    }).One(&user)

    if user.Token == "" {
        return models.User{}, errors.New("user not found")
    }

    return user, nil
}
