package action

import (
	"fmt"
	"glbackend/internal/adapters/api/response"
	"glbackend/internal/adapters/logging"
	"glbackend/internal/errorsStatus"
	"glbackend/internal/usecase"
	"glbackend/internal/utils"
	"net/http"
)

type ProductsFindByIDAction struct {
	uc  usecase.ProductsFindByIDUsecase
	log logging.Logger
}

func NewProductsFindByIDAction(uc usecase.ProductsFindByIDUsecase, log logging.Logger) *ProductsFindByIDAction {
	return &ProductsFindByIDAction{
		uc:  uc,
		log: log,
	}
}

func (a *ProductsFindByIDAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "products_find_by_id"

	var input usecase.ProductsFindByIDParams
	if err := utils.ParseParams(&input, r.URL.Query()); err != nil {
		logging.NewError(
			a.log,
			err,
			logKey,
			http.StatusBadRequest,
		).Log(utils.LogWithUser(r.Context(), "error when parse params"))
		response.NewError(w, err, http.StatusBadRequest)
		return
	}
	fmt.Printf("input: %+v\n", input)

	result, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		statusCode := errorsStatus.StatusCode(err)
		logging.NewError(
			a.log,
			err,
			logKey,
			statusCode,
		).Log(utils.LogWithUser(r.Context(), "error when find products"))

		response.NewError(w, err, statusCode)
		return
	}

	response.NewSuccess(w, result, http.StatusOK)
}
