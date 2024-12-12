package api

type GetParamResponse struct {
	// The parameter value, stringified for simplicity
	Value string `json:"value"`
}

type PostParamRequest struct {
	// ID of the parameter to set
	ID string `json:"id"`

	// The new value of the parameter
	Value string `json:"value"`
}
