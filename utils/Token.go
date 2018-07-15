package utils

import "math/rand"

const characters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateToken(size int) string {
    token := make([]byte, size)
    for i := range token {
        token[i] = characters[rand.Intn(len(characters))]
    }
    return string(token)
}
