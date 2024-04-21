package calculateexpression

import (
	"github.com/Mur466/distribcalc2/internal/errors"

	"github.com/Mur466/distribcalc2/internal/task"
)

type Service struct {
	repo Repo
}

func NewSvc(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Do(extexpr *ExtExpr) (ExprResult, error) {

	if extexpr.User == nil {
		return ExprResult{}, errors.ErrInvalidCreds
	}
	if extexpr.Ext_id != "" {
		// если есть ext_id, попробуем найти уже такое ранее созданное задание
		t := task.CheckUnique(extexpr.Ext_id, extexpr.User.Id)
		if t != nil {
			return ToExprResult(t), nil
		}
	}
	// у нас новое задание, требуем чтобы было выражение
	if extexpr.Expr == "" {
		return ExprResult{}, errors.ErrInvalidInput
	}
	// создаем задание
	t := task.NewTask(extexpr.Expr, extexpr.Ext_id, extexpr.User.Id)
	// добавляем в мапу
	task.Tasks[t.Task_id] = t
	//возвращаем результат
	return ToExprResult(t), nil

}

func ToExprResult(t *task.Task) ExprResult {
	return ExprResult{
		Ext_id:  t.Ext_id,
		Expr:    t.Expr,
		Task_id: t.Task_id,
		Status:  t.Status,
		Result:  t.Result,
		Message: t.Message,
	}
}
