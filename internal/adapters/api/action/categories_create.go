package action

import (
	"encoding/json"
	"fmt"
	"glbackend/internal/adapters/api/response"
	"glbackend/internal/adapters/logging"
	"glbackend/internal/errorsStatus"
	"glbackend/internal/usecase"
	"glbackend/internal/utils"
	"net/http"
)

type CategoriesCreateAction struct {
	uc  usecase.CategoriesCreateUsecase
	log logging.Logger
}

func NewCategoriesCreateAction(uc usecase.CategoriesCreateUsecase, log logging.Logger) *CategoriesCreateAction {
	return &CategoriesCreateAction{
		uc:  uc,
		log: log,
	}
}

func (a *CategoriesCreateAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "category_create"

	var input usecase.CategoriesCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log(utils.LogWithUser(r.Context(), "error when read json"))
		response.NewError(w, err, http.StatusBadRequest)
		return
	}
	fmt.Printf("input: %+v\n", input)

	category, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		statusCode := errorsStatus.StatusCode(err)
		logging.NewError(
			a.log,
			err,
			logKey,
			statusCode,
		).Log(utils.LogWithUser(r.Context(), "error when create category"))

		response.NewError(w, err, statusCode)
		return
	}

	response.NewSuccess(w, category, http.StatusOK)
}
