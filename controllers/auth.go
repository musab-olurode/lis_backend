package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/musab-olurode/lis_backend/constants"
	"github.com/musab-olurode/lis_backend/database"
	"github.com/musab-olurode/lis_backend/utils"
	"golang.org/x/crypto/bcrypt"
)

type AuthResponse struct {
	User         utils.UserWithoutPassword `json:"user"`
	AuthToken    string                    `json:"auth_token"`
	RefreshToken string                    `json:"refresh_token"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
		MatricNumber string `json:"matric_number"`
		Password     string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Invalid request payload: %v", err))
		return
	}

	params.MatricNumber = strings.ToUpper(params.MatricNumber)

	userExists := true
	_, err = database.DB.GetUserByMatricNumber(r.Context(), params.MatricNumber)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			userExists = false
		} else {
			utils.RespondWithInternalServerError(w, err)
			return
		}
	}

	if userExists {
		utils.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("Account with matric number %s already exists", params.MatricNumber))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	user, err := database.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		FirstName:    params.FirstName,
		LastName:     params.LastName,
		MatricNumber: params.MatricNumber,
		Role:         database.UserRoleUSER,
		Password:     string(hashedPassword),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	authToken, err := utils.GetAuthToken(user.ID.String())
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	refreshToken, err := utils.GetRefreshToken(user.ID.String())
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, AuthResponse{
		User:         utils.StripPassWordFromUser(user),
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, "Account created successfully")
}

func Login(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		MatricNumber string `json:"matric_number"`
		Password     string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		utils.RespondWithErr(w, http.StatusBadRequest, fmt.Sprintf("invalid request payload: %v", err))
		return
	}

	params.MatricNumber = strings.ToUpper(params.MatricNumber)

	user, err := database.DB.GetUserByMatricNumber(r.Context(), params.MatricNumber)
	if err != nil {
		if strings.Contains(err.Error(), "no rows") {
			utils.RespondWithErr(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		utils.RespondWithInternalServerError(w, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password))
	if err != nil {
		utils.RespondWithErr(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	authToken, err := utils.GetAuthToken(user.ID.String())
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	refreshToken, err := utils.GetRefreshToken(user.ID.String())
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, AuthResponse{
		User:         utils.StripPassWordFromUser(user),
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, "Login successful")
}

func GetLoggedInUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(constants.USER_CONTEXT_KEY).(database.User)
	utils.RespondWithJSON(w, http.StatusOK, utils.StripPassWordFromUser(user), "user retried successfully")
}

func RefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshTokenString := r.Header.Get("x-refresh-token")
	if refreshTokenString == "" {
		utils.RespondWithErr(w, http.StatusUnauthorized, "a refresh token is required")
		return
	}

	jwtClaims, err := utils.GetJWTClaims(refreshTokenString)
	if err != nil {
		utils.RespondWithErr(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	claims := *jwtClaims

	userID, err := uuid.Parse(claims["jti"].(string))
	if err != nil {
		utils.RespondWithErr(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	user, err := database.DB.GetUserByID(r.Context(), userID)
	if err != nil {
		utils.RespondWithErr(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	authToken, err := utils.GetAuthToken(user.ID.String())
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	refreshToken, err := utils.GetRefreshToken(user.ID.String())
	if err != nil {
		utils.RespondWithInternalServerError(w, err)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, AuthResponse{
		User:         utils.StripPassWordFromUser(user),
		AuthToken:    authToken,
		RefreshToken: refreshToken,
	}, "token refreshed successfully")
}
