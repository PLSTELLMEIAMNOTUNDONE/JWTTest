package tokens

import (
	"encoding/base64"
	"errors"
	"fmt"

	"math/rand"
	"time"

	"github.com/golang-jwt/jwt"
)

//см Readme
const KEY = "SECRET"

//Создает JWT токен
//в payload JWT токена вшиваем созданный с ним Refresh token
//Таким образом каждый JWT такое становиться уникальным а также всегда можно проверить что 
//Access и Refresh токены свяязаны
func CreateJWT (id string, refToken string) (string, error) {
	var (
	t   *jwt.Token
	s   string
	)
	claims := jwt.StandardClaims{
		Id: id,
		Subject: refToken,
	}
	t = jwt.NewWithClaims(jwt.SigningMethodHS256, claims) 
	s, e := t.SignedString([]byte(KEY)) 
	return s, e
}

//Создает случайный набор из 16-х байтов и превращает их в строку с кодировкой base64
//Так гарантируется уникальность каждого Refresh токена
func NewRefRefresh() string {
	tokenBts := make([]byte, 16)
    str := rand.NewSource(time.Now().Unix())
	randoms := rand.New(str)
	randoms.Read(tokenBts)
	token := base64.StdEncoding.EncodeToString(tokenBts)
	return token
}

//Функция проверяет что Access и Refresh токены были выданы вместе 
//и, если это так, выдает userId пользователя которому они были выданы
func CompareRefAndAccTokens(refToken string, accToken string) (string, error) {
	claims := jwt.StandardClaims{}
	
	_, er := jwt.ParseWithClaims(accToken, &claims, func(t *jwt.Token) (interface{}, error) {
	    if t.Method == jwt.SigningMethodHS256 {
			return []byte(KEY), nil
		}
		return nil, errors.New(fmt.Sprintf("invalid token %s", t.Raw))
	})
	if er != nil {
		return "", er
	}

	if claims.Subject == refToken {
		return claims.Id, nil
	}
	return "", errors.New(fmt.Sprintf("%s and %s are not match", accToken, refToken))


}