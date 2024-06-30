package controllers

import (
	"database/sql"
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
	"github.com/musab-olurode/lis_backend/database/models"
	"github.com/musab-olurode/lis_backend/utils"
)

func CreatePost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		Content       string    `json:"content"`
		CoverImageUrl string    `json:"cover_image_url"`
		Date          time.Time `json:"date"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	slug := utils.GetSlugFromTitle(params.Title)
	description := sql.NullString{
		String: params.Description,
		Valid:  params.Description != "",
	}

	if !description.Valid {
		description.String = utils.GetBlogContentDescription(params.Description)
		description.Valid = true
	}

	post, err := database.DB.CreatePost(r.Context(), database.CreatePostParams{
		ID:            uuid.New(),
		CoverImageUrl: params.CoverImageUrl,
		Title:         params.Title,
		Slug:          slug,
		Description:   description,
		Content:       params.Content,
		CreatedAt:     params.Date,
		UpdatedAt:     time.Now().UTC(),
	})
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusCreated, post)
}

func GetPosts(w http.ResponseWriter, r *http.Request) {
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

	posts, err := database.DB.GetPaginatedPosts(r.Context(), database.GetPaginatedPostsParams{
		Limit:  int32(perPage),
		Offset: int32((currentPage - 1) * perPage),
	})
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	postsCount, err := database.DB.CountPosts(r.Context())
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithPaginatedData(w, models.DatabasePostsToPosts(posts), int32(postsCount), int32(currentPage), int32(perPage), "posts retrieved successfully")
}

func GetPost(w http.ResponseWriter, r *http.Request) {
	postIDString := chi.URLParam(r, "postID")
	postID, err := uuid.Parse(postIDString)
	if err != nil {
		utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("post with id %s not found", postIDString))
		return
	}

	post, err := database.DB.GetPostByID(r.Context(), postID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("post with id %s not found", postIDString))
			return
		}

		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabasePostToPost(post), "post retrieved successfully")
}

func GetPostBySlug(w http.ResponseWriter, r *http.Request) {
	postSlug := chi.URLParam(r, "postSlug")

	post, err := database.DB.GetPostBySlug(r.Context(), postSlug)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("post with slug %s not found", postSlug))
			return
		}

		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabasePostToPost(post), "post retrieved successfully")
}

func UpdatePost(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		Content       string    `json:"content"`
		CoverImageUrl string    `json:"cover_image_url"`
		Date          time.Time `json:"date"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	postIDString := chi.URLParam(r, "postID")
	postID, err := uuid.Parse(postIDString)
	if err != nil {
		utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("post with id %s not found", postIDString))
		return
	}

	existingPost, err := database.DB.GetPostByID(r.Context(), postID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("post with id %s not found", postIDString))
			return
		}

		utils.RespondWithInternalServerError(w, err)
		return
	}

	slug := existingPost.Slug
	description := existingPost.Description
	coverImageUrl := existingPost.CoverImageUrl

	if params.Title != existingPost.Title {
		slug = utils.GetSlugFromTitle(params.Title)
	}

	if params.Description != existingPost.Description.String {
		description = sql.NullString{
			String: params.Description,
			Valid:  params.Description != "",
		}
	}

	if params.CoverImageUrl != existingPost.CoverImageUrl {
		utils.DeleteFileFromUrl(existingPost.CoverImageUrl)
		coverImageUrl = params.CoverImageUrl
	}

	if !description.Valid {
		description.String = utils.GetBlogContentDescription(params.Description)
		description.Valid = true
	}

	post, err := database.DB.UpdatePost(r.Context(), database.UpdatePostParams{
		ID:            postID,
		Title:         params.Title,
		Slug:          slug,
		CoverImageUrl: coverImageUrl,
		Description:   description,
		Content:       params.Content,
		CreatedAt:     params.Date,
		UpdatedAt:     time.Now().UTC(),
	})
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("post with id %s not found", postIDString))
			return
		}

		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, models.DatabasePostToPost(post), "post updated successfully")
}

func DeletePost(w http.ResponseWriter, r *http.Request) {
	postIDString := chi.URLParam(r, "postID")
	postID, err := uuid.Parse(postIDString)
	if err != nil {
		utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("post with id %s not found", postIDString))
		return
	}

	err = database.DB.DeletePost(r.Context(), postID)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			utils.RespondWithErr(w, http.StatusNotFound, fmt.Sprintf("post with id %s not found", postIDString))
			return
		}

		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, nil, "post deleted successfully")
}
