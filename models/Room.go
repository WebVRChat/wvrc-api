package models

import "gopkg.in/mgo.v2/bson"

type Room struct {
    Id          bson.ObjectId   `json:"id" bson:"_id,omitempty"`
    OwnerId     bson.ObjectId   `json:"owner_id"`
    Token       string          `json:"token"`
    Name        string          `json:"name"`
    Description string          `json:"description"`
    Users       []TmpUser       `json:"users"`
    IsPrivate   bool            `json:"is_private"`
}