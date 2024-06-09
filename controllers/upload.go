package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/musab-olurode/lis_backend/utils"
)

func UploadFile(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		File string `json:"file"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	url, err := utils.UploadFile(params.File)
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, map[string]string{"url": url})
}

func DeleteFile(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		File string `json:"file"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	matched, err := regexp.MatchString(`^https://res.cloudinary.com/dvn5ksfrp/image/upload/v\d+/.+\.\w+$`, params.File)
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}
	if !matched {
		utils.RespondWithErr(w, http.StatusBadRequest, "invalid file URL")
		return
	}

	err = utils.DeleteFileFromUrl(params.File)
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "file deleted"})
}
