package utility

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
)

var secretKey = []byte("Scanning-SecretKey")

type TokenReq struct {
	Id    uint64
	Token string
}

func (r *TokenReq) CreateJwtToken() (*TokenReq, error) {

	claims := jwt.MapClaims{
		"id":  r.Id,
		"exp": time.Now().Add(time.Hour * 2).Unix(),
	}

	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println("Error in creating JWT Token for User:", err)
		return nil, err
	}

	r.Token = tokenString

	return r, nil

}

func VerifyToken(endpoint http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		params := mux.Vars(r)
		id := params["id"]

		userid, _ := strconv.Atoi(id)

		accessToken := r.Header.Get("Authtoken")

		if accessToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			return secretKey, nil
		})

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims["id"] != uint64(userid) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

	})
}
