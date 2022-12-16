package oauth2service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"rederinghub.io/pkg"
)

type Auth2 struct {
	SecretKey string 
    signingKey string
    iss        string
}

func NewAuth2() *Auth2 {
	auth2 := &Auth2{}
	auth2.SecretKey = os.Getenv("AUTH_SECRET_KEY")
	auth2.signingKey = os.Getenv("AUTH_SIGNING_KEY")
	auth2.iss = os.Getenv("AUTH_TOKEN_ISS")
	return auth2
}

type SignedDetails struct {
    WalletAddress      string
    Email      string
    First_name string
    Last_name  string
    ID        string
    jwt.StandardClaims
}

func (a Auth2)  GenerateAllTokens(WalletAddress string, email string, firstName string, lastName string, uid string) (signedToken string, signedRefreshToken string, err error) {
    mySigningKey := []byte(a.signingKey)
    claims := &SignedDetails{
        WalletAddress: WalletAddress,
        Email:      email,
        First_name: firstName,
        Last_name:  lastName,
        ID:        uid,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(pkg.TOKEN_CACHE_EXPIRED_TIME)).Unix(),
        },
    }

    refreshClaims := &SignedDetails{
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Local().Add(time.Second * time.Duration(pkg.REFRESH_TOKEN_CACHE_EXPIRED_TIME)).Unix(),
        },
    }

    token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(mySigningKey))
    refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(mySigningKey))

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

func (a Auth2) ClaimToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, status.New(codes.InvalidArgument, fmt.Sprintf("Unexpected signing method: %v", token.Header["alg"])).Err()
		}

		return []byte(a.signingKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, status.New(codes.InvalidArgument, "invalid token").Err()
	}
	if err := claims.Valid(); err != nil {
		return nil, status.New(codes.InvalidArgument, "invalid token").Err()
	}
	// if !claims.VerifyIssuer(a.iss, true) {
	// 	return nil, status.New(codes.Unauthenticated, "unauthorized").Err()
	// }

	return claims, nil
}
