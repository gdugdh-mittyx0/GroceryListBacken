package action

import (
	"glbackend/internal/adapters/api/response"
	"glbackend/internal/adapters/logging"
	"glbackend/internal/errorsStatus"
	"glbackend/internal/usecase"
	"glbackend/internal/utils"
	"net/http"
)

type CategoriesFindAllAction struct {
	uc  usecase.CategoriesFindAllUsecase
	log logging.Logger
}

func NewCategoriesFindAllAction(uc usecase.CategoriesFindAllUsecase, log logging.Logger) *CategoriesFindAllAction {
	return &CategoriesFindAllAction{
		uc:  uc,
		log: log,
	}
}

func (a *CategoriesFindAllAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "categories_find_all"

	result, err := a.uc.Execute(r.Context())
	if err != nil {
		statusCode := errorsStatus.StatusCode(err)
		logging.NewError(
			a.log,
			err,
			logKey,
			statusCode,
		).Log(utils.LogWithUser(r.Context(), "error when categories find all"))

		response.NewError(w, err, statusCode)
		return
	}

	response.NewSuccessList(w, result, int(len(result)), http.StatusOK)
}
