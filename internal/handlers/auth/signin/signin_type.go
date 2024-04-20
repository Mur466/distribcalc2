package signin

type request struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Структура HTTP-ответа на вход в аккаунт
// В ответе содержится JWT-токен авторизованного пользователя
type LoginResponse struct {
    AccessToken string `json:"access_token"`
}