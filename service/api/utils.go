package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/nodeset-org/hyperdrive-example/shared/api"
)

// Logs the request and returns the query args and path args
func ProcessApiRequest(w http.ResponseWriter, r *http.Request, logger *slog.Logger, requestBody any) (url.Values, map[string]string) {
	args := r.URL.Query()
	logger.Info("New request",
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	)
	logger.Debug("Request params:",
		slog.String("query", r.URL.RawQuery),
	)

	if requestBody != nil {
		// Read the body
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			HandleInputError(w, logger, fmt.Errorf("error reading request body: %w", err))
			return nil, nil
		}
		logger.Debug("Request body:",
			slog.String("body", string(bodyBytes)),
		)

		// Deserialize the body
		err = json.Unmarshal(bodyBytes, &requestBody)
		if err != nil {
			HandleInputError(w, logger, fmt.Errorf("error deserializing request body: %w", err))
			return nil, nil
		}
	}

	return args, mux.Vars(r)
}

// Handle routes called with an invalid method
func HandleInvalidMethod(w http.ResponseWriter, logger *slog.Logger) {
	writeResponse(w, logger, http.StatusMethodNotAllowed, []byte{})
}

// Handles an error related to parsing the input parameters of a request
func HandleInputError(w http.ResponseWriter, logger *slog.Logger, err error) {
	msg := err.Error()
	bytes := formatError(msg)
	writeResponse(w, logger, http.StatusBadRequest, bytes)
}

// Write an error if something went wrong server-side
func HandleServerError(w http.ResponseWriter, logger *slog.Logger, err error) {
	msg := err.Error()
	bytes := formatError(msg)
	writeResponse(w, logger, http.StatusInternalServerError, bytes)
}

// The request completed successfully
func HandleSuccess[DataType any](w http.ResponseWriter, logger *slog.Logger, data DataType) {
	response := api.ApiResponse[DataType]{
		Error: "",
		Data:  data,
	}

	// Serialize the response
	bytes, err := json.Marshal(response)
	if err != nil {
		HandleServerError(w, logger, fmt.Errorf("error serializing response: %w", err))
	}
	// Write it
	logger.Debug("Response body",
		slog.String("body", string(bytes)),
	)
	writeResponse(w, logger, http.StatusOK, bytes)
}

// Writes a response to an HTTP request back to the client and logs it
func writeResponse(w http.ResponseWriter, logger *slog.Logger, statusCode int, message []byte) {
	// Prep the log attributes
	codeMsg := fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode))
	attrs := []any{
		slog.String("code", codeMsg),
	}

	// Log the response
	logMsg := "Responded with:"
	switch statusCode {
	case http.StatusOK:
		logger.Info(logMsg, attrs...)
	case http.StatusInternalServerError:
		logger.Error(logMsg, attrs...)
	default:
		logger.Warn(logMsg, attrs...)
	}

	// Write it to the client
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, writeErr := w.Write(message)
	if writeErr != nil {
		logger.Error("Error writing response", "error", writeErr)
	}
}

// JSONifies an error for responding to requests
func formatError(message string) []byte {
	msg := api.ApiResponse[struct{}]{
		Error: message,
		Data:  struct{}{},
	}

	bytes, _ := json.Marshal(msg)
	return bytes
}
