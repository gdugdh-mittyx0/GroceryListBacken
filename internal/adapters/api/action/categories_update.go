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

type CategoriesUpdateAction struct {
	uc  usecase.CategoriesUpdateUsecase
	log logging.Logger
}

func NewCategoriesUpdateAction(uc usecase.CategoriesUpdateUsecase, log logging.Logger) *CategoriesUpdateAction {
	return &CategoriesUpdateAction{
		uc:  uc,
		log: log,
	}
}

func (a *CategoriesUpdateAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "category_update"

	idStr := r.URL.Query().Get("id")
	id, err := utils.StringToUint(idStr)
	if err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log(utils.LogWithUser(r.Context(), "error when parse id"))
		response.NewError(w, err, http.StatusBadRequest)
		return
	}

	var input usecase.CategoriesUpdateInput
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
	input.ID = id
	fmt.Printf("input: %+v\n", input)

	category, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		statusCode := errorsStatus.StatusCode(err)
		logging.NewError(
			a.log,
			err,
			logKey,
			statusCode,
		).Log(utils.LogWithUser(r.Context(), "error when update category"))

		response.NewError(w, err, statusCode)
		return
	}

	response.NewSuccess(w, category, http.StatusOK)
}
