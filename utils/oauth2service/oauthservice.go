package oauth2service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"time"

	"rederinghub.io/utils"

	"github.com/golang-jwt/jwt"
)

type Auth2 struct {
	SecretKey string
}

func NewAuth2() *Auth2 {
	auth2 := &Auth2{}
	auth2.SecretKey = os.Getenv("AUTH_SECRET_KEY")
	return auth2
}

type SignedDetails struct {
	WalletAddress string // Replace eth address => btc segwit address
	Email         string
	First_name    string
	Last_name     string
	Uid           string
	jwt.StandardClaims
}

type VerifyResponse struct {
	IsVerified   bool   `json:"is_verified"`
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (a Auth2) GenerateAllTokens(WalletAddress string, email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		WalletAddress: WalletAddress,
		Email:         email,
		First_name:    firstName,
		Last_name:     lastName,
		Uid:           uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(utils.TOKEN_CACHE_EXPIRED_TIME)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(utils.REFRESH_TOKEN_CACHE_EXPIRED_TIME)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.SecretKey))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(a.SecretKey))

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

func (a Auth2) GenerateAllTokensUpdate(segwitBTCAddress string, email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &SignedDetails{
		WalletAddress: segwitBTCAddress,
		Email:         email,
		First_name:    firstName,
		Last_name:     lastName,
		Uid:           uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(utils.TOKEN_CACHE_EXPIRED_TIME)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(utils.REFRESH_TOKEN_CACHE_EXPIRED_TIME)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(a.SecretKey))
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(a.SecretKey))

	if err != nil {
		return "", "", err
	}

	return token, refreshToken, err
}

func (a Auth2) ValidateToken(signedToken string) (*SignedDetails, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(a.SecretKey), nil
		},
	)
	claims := &SignedDetails{}

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg := fmt.Sprintf("the token is invalid")
		return nil, errors.New(msg)
	}

	timeNow := time.Now().Local().Unix()
	if claims.ExpiresAt < timeNow {
		msg := fmt.Sprintf("token is expired")
		return nil, errors.New(msg)
	}

	return claims, nil
}

func (a Auth2) GenerateMd5String(input string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(input)))
}
