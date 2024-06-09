package utils

import (
	"math/rand"
	"strings"
	"time"

	"github.com/musab-olurode/lis_backend/database"
	"golang.org/x/net/html"
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

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func GetSlugFromTitle(title string) string {
	title = strings.ReplaceAll(title, " ", "-")
	title = strings.ToLower(title)
	title = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, title)
	if len(title) > 50 {
		title = title[:50]
	}
	title = title + "-" + GenerateRandomString(5)
	return title
}

func GetBlogContentDescription(content string) string {
	doc, err := html.Parse(strings.NewReader(content))
	if err != nil {
		return ""
	}

	var f func(*html.Node)
	var text strings.Builder

	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			text.WriteString(n.Data)
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	result := text.String()
	if len(result) > 100 {
		result = result[:100]
	}

	return result
}
