package main

import "net/http"

type State struct {
	ID   int32  `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}

type StatesResponse struct {
	Data  []State
	Error string
}

func StatesGetMany(w http.ResponseWriter, r *http.Request) {
	var response = StatesResponse{}
	for _, state := range mockData {
		response.Data = append(response.Data, State{
			ID:   state.ID,
			Code: state.Code,
			Name: state.Name,
		})
	}

	// @Codex: What's our return func?
	Encode[StatesResponse](w, r, http.StatusOK, response)
}

type StateResponse struct {
	State
	Error string `json:"error"`
}

func StateByCode(w http.ResponseWriter, r *http.Request) {
	var response = StateResponse{}
	code := r.PathValue("id")
	if code == "" {
		response.Error = "code is required"
		Encode[StateResponse](w, r, http.StatusBadRequest, response)
	}

	for _, state := range mockData {
		if state.Code == code {
			response.State = state
			Encode[StateResponse](w, r, http.StatusOK, response)
			return
		}
	}

	response.Error = "code not found"
	Encode[StateResponse](w, r, http.StatusBadRequest, response)
}

var mockData = []State{{
	ID:   1,
	Code: "mn",
	Name: "Minnesota",
}, {
	ID:   2,
	Code: "tx",
	Name: "Texas",
},
}
