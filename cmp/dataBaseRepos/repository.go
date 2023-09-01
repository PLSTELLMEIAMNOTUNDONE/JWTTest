package rep

import (
	"errors"
	"fmt"
	

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"context"
	"log"
)

const (
	//Правильнее было бы хранить эти константы в file.env 
	//см ReadMe
	MONGOADRRESS = "mongodb://localhost:27017/" 
	DBNAME= "test"
	COLLNAME = "users-tokens"
	

)
var collection *mongo.Collection
var ctx = context.TODO()
//Подключаемся к mongoDB
func Connect() {
	clientOptions := options.Client().ApplyURI(MONGOADRRESS)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
	collection = client.Database(DBNAME).Collection(COLLNAME)
}

//Функция находит в базе пользователя с id = 'userId'
//и, если такой пользователь был найден, устанавливает ему RefToken = 'hashRefToken'
//иначе кидает ошибку
func LinkIdAndRefToken (userId string, refToken string) error {
	hashRefToken, err := bcrypt.GenerateFromPassword([]byte(refToken), 10) // 10 - дефолтое значние для bcrypt
	
	if err != nil {
		return err
	}
	
	filter := bson.M{"_id": userId}
	change := bson.M{"$set": bson.M{"RefToken": string(hashRefToken)}}

	//Так как в задании речь шла об аутентификации, то пользователь с данными userId уже должен быть авторизирован 
	//и его Id уже должен быть в базе данных
	//соотвественно если мы не находим в базе данного Id то мы кидаем ошибку и НЕ создаем нового пользователя
	opts := options.Update().SetUpsert(false) 


	result, err := collection.UpdateOne(ctx, filter, change, opts)

	if err != nil {
		return err
	}
    //Этот случае соотвествует случаю когда пользователен с данным Id не был авторизирован 
	if result.MatchedCount == 0 {
		return errors.New(fmt.Sprintf("no user with id = %s", userId))
	}

	return nil
}
//В базе данных мы храним id пользователей и актуальный на данных момент для них Refresh Token
//Функция находит в базе пользователя с Id = 'userId' и проверяет что 'refToken' является актуальнным для него
func CheckRefTokenRelevance(userId string, refToken string) error {
    
	filter := bson.M{"_id": userId}
	var user bson.M
	result := collection.FindOne(ctx, filter)
	er := result.Decode(&user)
	if er != nil {
		return errors.New(fmt.Sprintf("no user with id = %s", userId))
	}
	hashRefToken := fmt.Sprintf("%v", user["RefToken"])
	if bcrypt.CompareHashAndPassword([]byte(hashRefToken), []byte(refToken)) != nil {
		return errors.New("password are not relevant anymore")
	}
	
	return nil
	
}
