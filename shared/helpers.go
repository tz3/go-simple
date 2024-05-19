package shared

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// HandleJSONError will handle the error if occurs in api handler functions
func HandleJSONError(w http.ResponseWriter, errMsg string, status int, log *zap.Logger, messages ...interface{}) {
	errMsg = fmt.Sprintf("[%d] %s", status, errMsg)

	errorResponse := ErrorResponse{
		Error: errMsg,
	}

	jsonResponse, err := json.Marshal(errorResponse)
	if err != nil {
		log.Error("Failed to marshal error response", zap.Error(err))
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(jsonResponse)
	if err != nil {
		log.Error("Error writing response", zap.Error(err))
	}

	logWithOptions := log.WithOptions(zap.AddCallerSkip(1)) // Adjust the caller skip value to skip the HandleJSONError function

	logWithOptions.Error(errMsg, convertToFields(messages)...)
}

func convertToFields(messages []interface{}) []zap.Field {
	fields := make([]zapcore.Field, len(messages))
	for i, msg := range messages {
		fields[i] = zap.Any("message", msg)
	}
	return fields
}
