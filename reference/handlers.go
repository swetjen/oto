package main

import "net/http"

type State struct {
	ID   int32  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type StatesResponse struct {
	Data  []State `json:"data"`
	Error string  `json:"error,omitempty"`
}

func StatesGetMany(w http.ResponseWriter, r *http.Request) {
	var response StatesResponse
	for _, state := range mockData {
		response.Data = append(response.Data, State{
			ID:   state.ID,
			Code: state.Code,
			Name: state.Name,
		})
	}

	Encode(w, r, http.StatusOK, response)
}

type StateResponse struct {
	State State  `json:"state"`
	Error string `json:"error,omitempty"`
}

func StateByCode(w http.ResponseWriter, r *http.Request) {
	var response StateResponse
	code := r.PathValue("code")
	if code == "" {
		response.Error = "code is required"
		Encode(w, r, http.StatusBadRequest, response)
		return
	}

	for _, state := range mockData {
		if state.Code == code {
			response.State = state
			Encode(w, r, http.StatusOK, response)
			return
		}
	}

	response.Error = "code not found"
	Encode(w, r, http.StatusBadRequest, response)
}

var mockData = []State{
	{
		ID:   1,
		Code: "mn",
		Name: "Minnesota",
	},
	{
		ID:   2,
		Code: "tx",
		Name: "Texas",
	},
}
