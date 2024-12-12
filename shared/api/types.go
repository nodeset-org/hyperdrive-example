package api

const (
	// API Route
	ApiRoute string = "api"
)

type ApiResponse[DataType any] struct {
	Data  DataType `json:"data"`
	Error string   `json:"error"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}
