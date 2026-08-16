package helpers

import (
	"database/sql"
	"gamebiller/configs"
	"gamebiller/models"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// Helper to get claims from context or Authorization header
func GetClaims(c echo.Context) (*models.JwtCustomClaims, bool) {
	if c == nil {
		return nil, false
	}
	tokenVal := c.Get("user")
	if tokenVal != nil {
		if token, ok := tokenVal.(*jwt.Token); ok {
			if claims, ok := token.Claims.(*models.JwtCustomClaims); ok {
				return claims, true
			}
		}
	}

	// Fallback for public routes where echojwt middleware is not registered:
	// Parse token from Authorization header if present
	if c.Request() != nil {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader != "" {
			tokenStr := strings.TrimSpace(authHeader)
			if strings.HasPrefix(strings.ToLower(tokenStr), "bearer ") {
				tokenStr = strings.TrimSpace(tokenStr[7:])
			}
			if tokenStr != "" {
				claims := new(models.JwtCustomClaims)
				token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
					return []byte(configs.APP_KEY), nil
				})
				if err == nil && token != nil && token.Valid {
					return claims, true
				}
			}
		}
	}

	return nil, false
}

// QuerySupport replaces '?' with postgres '$1', '$2', etc. placeholders
func QuerySupport(query string) string {
	count := strings.Count(query, "?")
	for i := 0; i < count; i++ {
		query = strings.Replace(query, "?", "$"+strconv.Itoa(i+1), 1)
	}
	return query
}

// HashPassword hashes a password string using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash compares a password with its hash
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// DBTransaction handles sql transaction with automatic rollback on error or panic
func DBTransaction(db *sql.DB, txFunc func(*sql.Tx) error) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			_ = tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()
	err = txFunc(tx)
	return err
}

// RandomDigits generates a random string of numbers with length n
func RandomDigits(n int) string {
	rand.Seed(time.Now().UnixNano())
	digits := "0123456789"
	result := make([]byte, n)
	for i := range result {
		result[i] = digits[rand.Intn(len(digits))]
	}
	return string(result)
}

// BuildWhereClause constructs WHERE query fragment and appends arguments.
func BuildWhereClause(baseWhere string, baseArgs []any, filters interface{}) (string, []any) {
	where := baseWhere
	args := baseArgs

	filterMap, ok := filters.(map[string]any)
	if !ok {
		return where, args
	}

	for col, val := range filterMap {
		if val == nil || val == "" {
			continue
		}
		cleanCol := sanitizeColumnName(col)
		if cleanCol == "" {
			continue
		}

		if where == "" {
			where = " WHERE " + cleanCol + " = ?"
		} else {
			where += " AND " + cleanCol + " = ?"
		}
		args = append(args, val)
	}

	return where, args
}

func sanitizeColumnName(col string) string {
	var result strings.Builder
	for _, r := range col {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_' {
			result.WriteRune(r)
		}
	}
	return result.String()
}
