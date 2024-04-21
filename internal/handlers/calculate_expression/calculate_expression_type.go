package calculateexpression

import "github.com/Mur466/distribcalc2/internal/entities"

type ExtExpr struct {
	Ext_id string `json:"ext_id"`
	Expr   string `json:"expr"`
	User   *entities.User
}

type ExprResult struct {
	Ext_id  string `json:"ext_id"`
	Expr    string `json:"expr"`
	Task_id int    `json:"task_id"`
	Status  string `json:"status"`
	Result  int64  `json:"result"`
	Message string `json:"message"`
}
