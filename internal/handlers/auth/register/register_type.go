package register

type request struct {
	Username string `json:"username"`
	Password1 string `json:"password1"`
	Password2 string `json:"password2"`
}

