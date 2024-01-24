package utils

import (
	"time"

	"github.com/musab-olurode/lis_backend/database"
)

type UserWithoutPassword struct {
	ID           string            `json:"id"`
	FirstName    string            `json:"first_name"`
	LastName     string            `json:"last_name"`
	MatricNumber string            `json:"matric_number"`
	Role         database.UserRole `json:"role"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

func StripPassWordFromUser(user database.User) UserWithoutPassword {
	return UserWithoutPassword{
		ID:           user.ID.String(),
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		MatricNumber: user.MatricNumber,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
	}
}
