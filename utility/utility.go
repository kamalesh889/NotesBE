package utility

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"github.com/juju/ratelimit"
	"github.com/spf13/viper"
)

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

	tokenString, err := token.SignedString([]byte(viper.GetString("token.secretkey")))
	if err != nil {
		log.Println("Error in creating JWT Token for User:", err)
		return nil, err
	}

	r.Token = tokenString

	return r, nil

}

func VerifyToken(endpoint http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userId, err := ParseUserId(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		accessToken := r.Header.Get("Authtoken")

		if accessToken == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(viper.GetString("token.secretkey")), nil
		})

		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if claims["id"].(float64) != float64(userId) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		} else {
			endpoint(w, r)
		}

	})
}

func RateLimitMiddleware(l *limiter.Limiter) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			httpError := tollbooth.LimitByRequest(l, w, r)
			if httpError != nil {
				http.Error(w, httpError.Message, httpError.StatusCode)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func RequestThrottleMiddleware(limiter *ratelimit.Bucket) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if limiter.TakeAvailable(1) < 1 {
				http.Error(w, "Too many requests", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func ParseUserId(r *http.Request) (uint64, error) {
	query := r.URL.Query()
	user := query.Get("userid")

	userId, err := strconv.ParseUint(user, 10, 64)
	if err != nil {
		return 0, err
	}
	return userId, err
}

func ParseNoteId(r *http.Request) (uint64, error) {
	vars := mux.Vars(r)

	note, found := vars["id"]
	if !found {
		return 0, errors.New("noteId is not present")
	}
	noteId, err := strconv.ParseUint(note, 10, 64)
	if err != nil {
		return 0, err
	}
	return noteId, err
}
