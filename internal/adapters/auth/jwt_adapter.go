package auth

import (
	"time"

	"github.com/RiosHectorM/iso-stack/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
)

type JWTAdapter struct {
	Secret string
}

type CustomClaims struct {
	UserID string `json:"user_id"`
	OrgID  string `json:"org_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (j *JWTAdapter) GenerateToken(user *domain.User) (string, error) {
	claims := CustomClaims{
		UserID: user.ID,
		OrgID:  user.OrganizationID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.Secret))
}
