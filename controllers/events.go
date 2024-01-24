package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/musab-olurode/lis_backend/constants"
	"github.com/musab-olurode/lis_backend/database"
	"github.com/musab-olurode/lis_backend/utils"
)

func CreateEvent(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string         `json:"title"`
		Description sql.NullString `json:"description"`
		ImageUrl    string         `json:"image_url"`
		Venue       string         `json:"venue"`
		Date        time.Time      `json:"date"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	event, err := database.DB.CreateEvent(r.Context(), database.CreateEventParams{
		ID:          uuid.New(),
		Title:       params.Title,
		Description: params.Description,
		ImageUrl:    params.ImageUrl,
		Venue:       params.Venue,
		Date:        params.Date,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, event)
}

func GetEvents(w http.ResponseWriter, r *http.Request) {
	perPageQuery := r.URL.Query().Get("per_page")
	pageQuery := r.URL.Query().Get("page")

	perPage, err := strconv.Atoi(perPageQuery)
	if err != nil || perPage < 1 || perPage > constants.ITEMS_PER_PAGE {
		perPage = constants.ITEMS_PER_PAGE
	}
	currentPage, err := strconv.Atoi(pageQuery)
	if err != nil || currentPage < 1 {
		currentPage = 1
	}

	events, err := database.DB.GetPaginatedEvents(r.Context(), database.GetPaginatedEventsParams{
		Limit:  int32(perPage),
		Offset: int32((currentPage - 1) * perPage),
	})
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	eventsCount, err := database.DB.CountEvents(r.Context())
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithPaginatedData(w, events, int32(eventsCount), int32(currentPage), int32(perPage), "events retrieved successfully")
}

func GetEvent(w http.ResponseWriter, r *http.Request) {
	eventIDString := chi.URLParam(r, "eventID")
	eventID, err := uuid.Parse(eventIDString)
	if err != nil {
		utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("event with id %s not found", eventIDString))
		return
	}

	event, err := database.DB.GetEventByID(r.Context(), eventID)
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, event, "event retrieved successfully")
}

func UpdateEvent(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string         `json:"title"`
		Description sql.NullString `json:"description"`
		ImageUrl    string         `json:"image_url"`
		Venue       string         `json:"venue"`
		Date        time.Time      `json:"date"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	eventIDString := chi.URLParam(r, "eventID")
	eventID, err := uuid.Parse(eventIDString)
	if err != nil {
		utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("event with id %s not found", eventIDString))
		return
	}

	event, err := database.DB.UpdateEvent(r.Context(), database.UpdateEventParams{
		ID:          eventID,
		Title:       params.Title,
		Description: params.Description,
		ImageUrl:    params.ImageUrl,
		Venue:       params.Venue,
		Date:        params.Date,
		UpdatedAt:   time.Now().UTC(),
	})
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, event, "event updated successfully")
}

func DeleteEvent(w http.ResponseWriter, r *http.Request) {
	eventIDString := chi.URLParam(r, "eventID")
	eventID, err := uuid.Parse(eventIDString)
	if err != nil {
		utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("event with id %s not found", eventIDString))
		return
	}

	err = database.DB.DeleteEvent(r.Context(), eventID)
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, nil, "event deleted successfully")
}
