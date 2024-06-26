package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/musab-olurode/lis_backend/constants"
	"github.com/musab-olurode/lis_backend/database"
	"github.com/musab-olurode/lis_backend/utils"
)

func CreateMaterial(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		FileUrl     string `json:"file_url"`
		Category    string `json:"category"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	categoryEnum := database.Category(params.Category)

	material, err := database.DB.CreateMaterial(r.Context(), database.CreateMaterialParams{
		ID:        uuid.New(),
		Title:     params.Title,
		FileUrl:   params.FileUrl,
		Category:  categoryEnum,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, material)
}

func GetMaterials(w http.ResponseWriter, r *http.Request) {
	perPageQuery := r.URL.Query().Get("per_page")
	pageQuery := r.URL.Query().Get("page")
	category := r.URL.Query().Get("category")

	perPage, err := strconv.Atoi(perPageQuery)
	if err != nil || perPage < 1 || perPage > constants.ITEMS_PER_PAGE {
		perPage = constants.ITEMS_PER_PAGE
	}
	currentPage, err := strconv.Atoi(pageQuery)
	if err != nil || currentPage < 1 {
		currentPage = 1
	}

	materials := []database.Material{}
	var materialsCount int64 = 0

	if category != "" || category == "All" || category == "undefined" {
		categoryEnum := database.Category(category)
		materials, err = database.DB.GetPaginatedMaterialsByCategory(r.Context(), database.GetPaginatedMaterialsByCategoryParams{
			Category: categoryEnum,
			Limit:    int32(perPage),
			Offset:   int32((currentPage - 1) * perPage),
		})
		if err != nil {
			utils.RespondWithInternalServerError(w, err)
			return
		}

		materialsCount, err = database.DB.CountMaterialsByCategory(r.Context(), categoryEnum)
		if err != nil {
			utils.RespondWithInternalServerError(w, err)
			return
		}
	} else {
		materials, err = database.DB.GetPaginatedMaterials(r.Context(), database.GetPaginatedMaterialsParams{
			Limit:  int32(perPage),
			Offset: int32((currentPage - 1) * perPage),
		})
		if err != nil {
			utils.RespondWithInternalServerError(w, err)
			return
		}

		materialsCount, err = database.DB.CountMaterials(r.Context())
		if err != nil {
			utils.RespondWithInternalServerError(w, err)
			return
		}
	}

	utils.RespondWithPaginatedData(w, materials, int32(materialsCount), int32(currentPage), int32(perPage), "materials retrieved successfully")
}

func GetMaterial(w http.ResponseWriter, r *http.Request) {
	materialIDString := chi.URLParam(r, "materialID")
	materialID, err := uuid.Parse(materialIDString)
	if err != nil {
		utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("material with id %s not found", materialIDString))
		return
	}

	material, err := database.DB.GetMaterialByID(r.Context(), materialID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("material with id %s not found", materialIDString))
			return
		}

		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, material, "material retrieved successfully")
}

func UpdateMaterial(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		FileUrl     string `json:"file_url"`
		Category    string `json:"category"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	categoryEnum := database.Category(params.Category)

	materialIDString := chi.URLParam(r, "materialID")
	materialID, err := uuid.Parse(materialIDString)
	if err != nil {
		utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("material with id %s not found", materialIDString))
		return
	}

	material, err := database.DB.UpdateMaterial(r.Context(), database.UpdateMaterialParams{
		ID:        materialID,
		Title:     params.Title,
		FileUrl:   params.FileUrl,
		Category:  categoryEnum,
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("material with id %s not found", materialIDString))
			return
		}

		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, material, "material updated successfully")
}

func DeleteMaterial(w http.ResponseWriter, r *http.Request) {
	materialIDString := chi.URLParam(r, "materialID")
	materialID, err := uuid.Parse(materialIDString)
	if err != nil {
		utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("material with id %s not found", materialIDString))
		return
	}

	err = database.DB.DeleteMaterial(r.Context(), materialID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("material with id %s not found", materialIDString))
			return
		}

		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, nil, "material deleted successfully")
}
