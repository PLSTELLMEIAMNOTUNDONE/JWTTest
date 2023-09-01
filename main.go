package main

import (
	repos "jwt/cmp/dataBaseRepos"
	hand "jwt/cmp/handlers"
	"net/http"
	_ "jwt/docs"
	"fmt"  
	//Это импорт нужен исключительно для работы swagger 
	//Программа может функционировать и без него
	//см ReadMe
	"github.com/swaggo/http-swagger" 
)



func init () {
	repos.Connect()
	
}
// @title Swagger Example
// @version 1.1
// @description jwt token test

// @host localhost:5000
// @BasePath /
// @name jwt
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/first", hand.HandleGetTokens)
	mux.HandleFunc("/second", hand.HandleRefresh)
	//Это эндпоинт для swagger, см Readme
	mux.HandleFunc("/swaggerd/", httpSwagger.WrapHandler) 
	fmt.Println("start")
    fmt.Println(http.ListenAndServe(":5000", mux).Error())
	
}