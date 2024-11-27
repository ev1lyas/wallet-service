package api

import "net/http"

// ResponseError отправляет стандартные ошибки с сообщениями
func ResponseError(w http.ResponseWriter, statusCode int, message string) {
	http.Error(w, message, statusCode)
}
