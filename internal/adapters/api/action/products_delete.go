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

type ProductsDeleteAction struct {
	uc  usecase.ProductsDeleteUsecase
	log logging.Logger
}

func NewProductsDeleteAction(uc usecase.ProductsDeleteUsecase, log logging.Logger) *ProductsDeleteAction {
	return &ProductsDeleteAction{
		uc:  uc,
		log: log,
	}
}

func (a *ProductsDeleteAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "products_delete"

	var input usecase.ProductsDeleteParams
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

	if err := a.uc.Execute(r.Context(), input); err != nil {
		statusCode := errorsStatus.StatusCode(err)
		logging.NewError(
			a.log,
			err,
			logKey,
			statusCode,
		).Log(utils.LogWithUser(r.Context(), "error when delete product"))

		response.NewError(w, err, statusCode)
		return
	}

	response.NewSuccess(w, true, http.StatusOK)
}
