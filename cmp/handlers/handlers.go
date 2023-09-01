package hand

import (
	"encoding/json"
	
	repos "jwt/cmp/dataBaseRepos"
	tokens "jwt/cmp/tokenizer"
	"net/http"
)
type tokenJs struct {
	AccToken string `json:"AccToken"`
	RefToken string `json:"RefToken"`
}
//В случае успешного выполнения функия записывает в response body
//json с Access и Refresh ключами
//а также изменяет в базе пользователя с userId = id  RefToken на новый
func makeTokensById(wp *http.ResponseWriter, r *http.Request, id string)  {
	refToken := tokens.NewRefRefresh()
	accToken, er := tokens.CreateJWT(id, refToken)
	
	w := *wp
	if er != nil {
		//Эта ошибка возникает только в случае если секретный имеет невалидный тип - это ошибка на стороне сервера
		w.WriteHeader(http.StatusInternalServerError) 
		w.Write([]byte(er.Error()))
		return
	}
	
	tokensJs := tokenJs{AccToken: accToken, RefToken: refToken}
	err := repos.LinkIdAndRefToken(id, refToken)
	resp, er := json.Marshal(tokensJs)
	if er != nil {
		//Эта ошибка может возникнуть только в случае неуспешного выполнения json.Marshal - ошибка на стороне сервера
		w.WriteHeader(http.StatusInternalServerError) 
		w.Write([]byte(er.Error()))
		return
	}
	if err != nil {
		//Эта ошибка возникает только если userId содержит невалидные символы либо не был найден в базе - ошибка запроса
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.WriteHeader(http.StatusAccepted)
	w.Header().Set("Content-Type", "application/json")
    w.Write(resp)
}


// @Summary get acc and ref tokens
// @Description sdf
// @Produce json
// @ID getFirst
// @Param id query string true "id"
// @Success 200 {object} error
// @Failure 400 {string} error
// @Router /first [post]
func HandleGetTokens (w http.ResponseWriter, r *http.Request) {
    id := r.URL.Query().Get("id")
	makeTokensById(&w, r, id)
}


// @Summary refresh acc and ref tokens
// @Description sfdgddf
// @Accept json
// @Produce json
// @ID Second
// @Param token body tokenJs true "token"
// @Success 200 {object} error
// @Failure 400 {string} error
// @Router /second [post]
func HandleRefresh(w http.ResponseWriter, r *http.Request) {
	
	var tokensFromBody tokenJs
	er := json.NewDecoder(r.Body).Decode(&tokensFromBody)
	
	
	if er != nil {
		//Эта ошибка возникает только если один из токенов содержит невалидные символы - ошибка запроса
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(er.Error()))
		return
	}
	
	userId, er := tokens.CompareRefAndAccTokens(tokensFromBody.RefToken, tokensFromBody.AccToken)
	if er != nil {
		//Эта ошибка возникает только если Access Token и Refresh Token НЕ были выданы вместе - ошибка запроса
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(er.Error()))
		return
	}
	er = repos.CheckRefTokenRelevance(userId, tokensFromBody.RefToken)
    if er != nil {
		//Эта ошибка возникает только если Refresh Token уже был использован - ошибка запроса
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(er.Error()))
		return
	}
	//Насколько я понял этот эндпоинт должен возращать такой же по структуре json как и первый эндпоинт см localhost:5000/swaggerd
	makeTokensById(&w, r, userId)
	
	
}
