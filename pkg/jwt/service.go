package jwt

import (
	"errors"
	"math/rand"
	"net"
	"net/http"
	"netradio/internal/repository"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const (
	AccessTokenTTL  = time.Minute * 30
	RefreshTokenTTL = time.Hour * 24 * 10
)

var config Config

func SetConfig(cfg Config) {
	config = cfg
}

func NewAccessToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		ID:        strconv.Itoa(userID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenTTL)),
	})

	jwtSignedValue, err := token.SignedString([]byte(config.SecretJWTKey))
	if err != nil {
		return "", err
	}

	return string(jwtSignedValue), nil
}

func NewRefreshToken() string {
	symbols := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, 32)
	for i := range b {
		b[i] = symbols[rand.Intn(len(symbols))]
	}
	return string(b)
}

func AddRefreshTokenToCookie(w http.ResponseWriter, r *http.Request, userID int) {
	refreshToken := NewRefreshToken()
	expires := time.Now().Add(RefreshTokenTTL)
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	repository.NewUserDB().CreateSession(userID, refreshToken, expires, ip)
	http.SetCookie(w, &http.Cookie{Name: "refreshToken", Value: refreshToken, Expires: expires, Path: "/", SameSite: http.SameSiteNoneMode, HttpOnly: true, Secure: true})
}

func GetUserIDFromHeader(header http.Header) (int, error) {
	headerValue := strings.Split(header.Get("Authorization"), "Bearer ")
	if len(headerValue) != 2 {
		return 0, errors.New("Access token not found")
	}

	token, err := jwt.ParseWithClaims(headerValue[1], &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenUnverifiable
		}
		return []byte(config.SecretJWTKey), nil
	})
	if err != nil {
		return 0, errors.New("Invalid access token")
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return 0, errors.New("Invalid access token")
	}

	id, err := strconv.Atoi(claims.ID)
	if err != nil {
		return 0, err
	}
	return id, nil
}
