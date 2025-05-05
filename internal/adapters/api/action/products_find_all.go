package action

import (
	"glbackend/internal/adapters/api/response"
	"glbackend/internal/adapters/logging"
	"glbackend/internal/errorsStatus"
	"glbackend/internal/usecase"
	"glbackend/internal/utils"
	"net/http"
)

type ProductsFindAllAction struct {
	uc  usecase.ProductsFindAllUsecase
	log logging.Logger
}

func NewProductsFindAllAction(uc usecase.ProductsFindAllUsecase, log logging.Logger) *ProductsFindAllAction {
	return &ProductsFindAllAction{
		uc:  uc,
		log: log,
	}
}

func (a *ProductsFindAllAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "products_find_all"

	result, err := a.uc.Execute(r.Context())
	if err != nil {
		statusCode := errorsStatus.StatusCode(err)
		logging.NewError(
			a.log,
			err,
			logKey,
			statusCode,
		).Log(utils.LogWithUser(r.Context(), "error when products find all"))

		response.NewError(w, err, statusCode)
		return
	}

	response.NewSuccessList(w, result, int(len(result)), http.StatusOK)
}
