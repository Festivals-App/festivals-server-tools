package servertools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/zerolog/log"
)

func RespondCode(w http.ResponseWriter, code int) {

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
}

func RespondString(w http.ResponseWriter, code int, message string) {

	response := []byte(message)
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}

func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {

	//TODO String comparison is not very elegant!
	if fmt.Sprint(payload) == "[]" {
		payload = []interface{}{}
	}

	resultMap := map[string]interface{}{"data": payload}
	response, err := json.Marshal(resultMap)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal payload")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("failed to write response")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(response)
	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}

func RespondError(w http.ResponseWriter, code int, message string) {

	resultMap := map[string]interface{}{"error": message}
	response, err := json.Marshal(resultMap)
	if err != nil {
		log.Error().Err(err).Msg("failed to marshal payload")
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			log.Error().Err(err).Msg("failed to write response")
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err = w.Write(response)
	if err != nil {
		log.Error().Err(err).Msg("failed to write response")
	}
}

// respondFile makes the response with payload as json format
func RespondFile(w http.ResponseWriter, file *os.File) {

	// calculate content size
	fileInfo, err := file.Stat()
	if err != nil || fileInfo == nil {
		log.Error().Err(err).Msg("Failed to read file stats for file: '" + file.Name() + "'")
		RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	size := fileInfo.Size()

	// calculate content type dynamically
	contentType, err := GetFileContentType(file)
	if err != nil {
		log.Error().Err(err).Msg("Failed to retrieve content type for file: '" + file.Name() + "'")
		RespondError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Length", strconv.FormatInt(size, 10))
	w.Header().Set("Content-Type", contentType)

	_, err = io.Copy(w, file)
	if err != nil {
		log.Error().Err(err).Msg("Failed to send write file to response")
	}
}

func UnauthorizedResponse(w http.ResponseWriter) {

	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	RespondError(w, http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))
}
