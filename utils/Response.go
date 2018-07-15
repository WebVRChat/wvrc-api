package utils

import (
    "net/http"
    "encoding/json"
)

type ApplicationResponse struct {
    Code    int
    Data    interface{}
}

func (ap *ApplicationResponse) SendResponse(writer http.ResponseWriter) {
    response, _ := json.Marshal(ap.Data)
    writer.Header().Set("Content-Type", "application/json")
    writer.WriteHeader(ap.Code)
    writer.Write(response)
}

func TriggerError(writer http.ResponseWriter, status int, reason string) {
    response := ApplicationResponse{status, map[string]string{
        "status": "failed",
        "reason": reason,
    }}
    response.SendResponse(writer)
}