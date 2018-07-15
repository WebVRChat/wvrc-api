package models

import "gopkg.in/mgo.v2/bson"

type User struct {
    Id          bson.ObjectId           `json:"id" bson:"_id,omitempty"`
    Username    string                  `json:"username"`
    Password    string                  `json:"password"`
    Token       string                  `json:"token"`
}

func (u User) ToTmp(peer_id string) TmpUser {
    return TmpUser{
        Username: u.Username,
        PeerId: peer_id,
    }
}

type TmpUser struct {
    Username    string  `json:"username"`
    PeerId      string  `json:"peer_id"`
}