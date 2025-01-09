package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/nodeset-org/hyperdrive-example/shared/api"
	"github.com/nodeset-org/hyperdrive-example/shared/config"
)

func (s *ApiServer) HandleParam(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleParamGet(w, r)
	case http.MethodPost:
		s.handleParamPost(w, r)
	default:
		HandleInvalidMethod(w, s.logger)
	}
}

// Handle a GET request to get a parameter
func (s *ApiServer) handleParamGet(w http.ResponseWriter, r *http.Request) {
	// Get the args
	args, _ := ProcessApiRequest(w, r, s.logger, nil)

	// Get the param
	if !args.Has("id") {
		HandleInputError(w, s.logger, fmt.Errorf("missing required parameter [id]"))
		return
	}

	var response api.GetParamResponse
	param := args.Get("id")
	switch param {
	case "exampleBool":
		response.Value = strconv.FormatBool(s.cfgMgr.Config.ExampleBool)
	case "exampleInt":
		response.Value = strconv.FormatInt(s.cfgMgr.Config.ExampleInt, 10)
	case "exampleUint":
		response.Value = strconv.FormatUint(s.cfgMgr.Config.ExampleUint, 10)
	case "exampleFloat":
		response.Value = strconv.FormatFloat(s.cfgMgr.Config.ExampleFloat, 'f', -1, 64)
	case "exampleString":
		response.Value = s.cfgMgr.Config.ExampleString
	case "exampleChoice":
		response.Value = string(s.cfgMgr.Config.ExampleChoice)
	case "subBool":
		response.Value = strconv.FormatBool(s.cfgMgr.Config.SubConfig.SubExampleBool)
	case "subChoice":
		response.Value = string(s.cfgMgr.Config.SubConfig.SubExampleChoice)
	default:
		HandleInputError(w, s.logger, fmt.Errorf("invalid parameter [%s]", param))
	}

	// Send the response
	HandleSuccess(w, s.logger, response)
}

// Handle a POST request to set a parameter
func (s *ApiServer) handleParamPost(w http.ResponseWriter, r *http.Request) {
	// Get the body
	var request api.PostParamRequest
	queryArgs, pathArgs := ProcessApiRequest(w, r, s.logger, &request)
	if queryArgs == nil && pathArgs == nil {
		return
	}

	// Set the param
	var err error
	switch request.ID {
	case "exampleBool":
		s.cfgMgr.Config.ExampleBool, err = strconv.ParseBool(request.Value)
	case "exampleInt":
		s.cfgMgr.Config.ExampleInt, err = strconv.ParseInt(request.Value, 10, 64)
	case "exampleUint":
		s.cfgMgr.Config.ExampleUint, err = strconv.ParseUint(request.Value, 10, 64)
	case "exampleFloat":
		s.cfgMgr.Config.ExampleFloat, err = strconv.ParseFloat(request.Value, 64)
	case "exampleString":
		s.cfgMgr.Config.ExampleString = request.Value
	case "exampleChoice":
		s.cfgMgr.Config.ExampleChoice = config.ExampleOption(request.Value)
	case "subBool":
		s.cfgMgr.Config.SubConfig.SubExampleBool, err = strconv.ParseBool(request.Value)
	case "subChoice":
		s.cfgMgr.Config.SubConfig.SubExampleChoice = config.ExampleOption(request.Value)
	default:
		HandleInputError(w, s.logger, fmt.Errorf("invalid parameter [%s]", request.ID))
	}
	if err != nil {
		HandleInputError(w, s.logger, fmt.Errorf("error setting parameter [%s]: %w", request.ID, err))
		return
	}

	// Send the response
	HandleSuccess[struct{}](w, s.logger, struct{}{})
}
