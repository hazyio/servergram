package internal

import (
	"encoding/json"
	"net/http"
)

type JsonResponse struct {
	writer http.ResponseWriter
}

func NewResponse(w http.ResponseWriter) *JsonResponse {
	return &JsonResponse{
		writer: w,
	}
}
func (response *JsonResponse) WriteSuccess(message map[string]string) {
	response.writer.Header().Set("Content-Type", "application/json")
	response.writer.WriteHeader(http.StatusOK)
	json.NewEncoder(response.writer).Encode(message)
}
func (response *JsonResponse) WriteError(message string) {
	response.writer.Header().Set("Content-Type", "application/json")
	response.writer.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(response.writer).Encode(map[string]string{"error": message})
}
func (response *JsonResponse) WriteBadRequest(message string) {
	response.writer.Header().Set("Content-Type", "application/json")
	response.writer.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(response.writer).Encode(map[string]string{"error": message})
}
func (response *JsonResponse) Write(status int, message interface{}) {
	response.writer.Header().Set("Content-Type", "application/json")
	response.writer.WriteHeader(status)
	json.NewEncoder(response.writer).Encode(message)
}

func (response *JsonResponse) WriteWithCookie(status int, message interface{}, cookie *http.Cookie) {
	http.SetCookie(response.writer, cookie)
	response.writer.Header().Set("Content-Type", "application/json")
	response.writer.WriteHeader(status)
	json.NewEncoder(response.writer).Encode(message)
}
