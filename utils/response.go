package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type successResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type errResponse struct {
	Error string `json:"error"`
}

type pagination struct {
	PerPage int32 `json:"per_page"`
	Page    int32 `json:"page"`
	Total   int32 `json:"total"`
}

type paginatedResponse struct {
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
	Pagination pagination  `json:"pagination"`
}

func RespondWithInternalServerError(w http.ResponseWriter, err error) {
	log.Printf("internal server error: %v", err)
	RespondWithErr(w, http.StatusInternalServerError, "internal server error")
}

func RespondWithErr(w http.ResponseWriter, code int, msg string) {
	if code > 499 && msg != "internal server error" {
		log.Printf("Responding with 5xx error: %s", msg)
	}

	RespondWithJSON(w, code, errResponse{Error: msg})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}, message ...string) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errResponse, _ := json.Marshal(errResponse{Error: "internal server error"})
		w.Write(errResponse)
		return
	}
	if code >= 200 && code < 300 {
		msg := "success"
		if len(message) > 0 {
			msg = message[0]
		}
		data, err = json.Marshal(successResponse{Message: msg, Data: payload})
		if err != nil {
			log.Printf("Error marshalling JSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			errResponse, _ := json.Marshal(errResponse{Error: "internal server error"})
			w.Write(errResponse)
			return
		}
	}
	w.WriteHeader(code)
	w.Write(data)
}

func RespondWithPaginatedData(w http.ResponseWriter, payload interface{}, total int32, current_page int32, per_page int32, message ...string) {
	w.Header().Set("Content-Type", "application/json")
	_, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errResponse, _ := json.Marshal(errResponse{Error: "internal server error"})
		w.Write(errResponse)
		return
	}

	msg := "success"
	if len(message) > 0 {
		msg = message[0]
	}

	paginatedResponse := paginatedResponse{
		Message: msg,
		Data:    payload,
		Pagination: pagination{
			PerPage: per_page,
			Page:    current_page,
			Total:   total,
		},
	}

	data, err := json.Marshal(paginatedResponse)
	if err != nil {
		log.Printf("Error marshalling JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		errResponse, _ := json.Marshal(errResponse{Error: "internal server error"})
		w.Write(errResponse)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
