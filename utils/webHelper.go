package utils

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"google.golang.org/grpc"
)

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func WithGRPCConnection(serverType string, endpointHandler func(conn *grpc.ClientConn, writer http.ResponseWriter, request *http.Request)) http.HandlerFunc {
	servers := map[string]string{
		"authServer":     "localhost:4040",
		"activityServer": "localhost:4041",
	}
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		conn, err := grpc.Dial(servers[serverType], grpc.WithInsecure())
		if err != nil {
			fmt.Fprintf(writer, "Internal Server Error")
			return
		}
		defer conn.Close()
		endpointHandler(conn, writer, request)
	})
}

func VerifyJWT(endpointHandler func(writer http.ResponseWriter, request *http.Request, claims *Claims)) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tokenObject, err := request.Cookie("token")
		claims := &Claims{}
		if err == nil {
			tkn, err := jwt.ParseWithClaims(tokenObject.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(UtilsConfig.SecretKey), nil
			})
			if err != nil {
				if err == jwt.ErrSignatureInvalid {
					fmt.Fprintf(writer, "Status Unauthorized")
					return
				} else {
					fmt.Fprintf(writer, err.Error())
					return
				}
			}
			if !tkn.Valid {
				fmt.Fprintf(writer, "Status Unauthorized")
				return
			}
			// #============================================================================================
			// refreshing the token here
			// #============================================================================================
			expirationTime := time.Now().Add(5 * time.Minute)
			claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			tokenString, err := token.SignedString([]byte(UtilsConfig.SecretKey))
			if err != nil {
				fmt.Fprintf(writer, "Internal Server Error")
				return
			}
			http.SetCookie(writer, &http.Cookie{
				Name:    "token",
				Value:   tokenString,
				Expires: expirationTime,
			})
			// #================================================================================================
			endpointHandler(writer, request, claims)
			return
		} else {
			if err == http.ErrNoCookie {
				fmt.Fprintf(writer, "Status Unauthorized")
				return
			}
			fmt.Fprintf(writer, "Status Bad Request")
			return

		}
	})
}

func VerifyToken(tokenString string) (bool, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodECDSA)
		if !ok {
			return nil, errors.New("invalid token")
		}
		return "", nil
	})
	if err != nil {
		return false, err
	}
	if token.Valid {
		return true, nil
	}
	return false, errors.New("invalid token")
}

func GenerateJWT(username string) (string, time.Time, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(UtilsConfig.SecretKey))
	if err != nil {
		return "", time.Time{}, err
	}
	return tokenString, expirationTime, nil
}
