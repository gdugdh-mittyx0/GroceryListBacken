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

type ProductsStatusesUpdateAction struct {
	uc  usecase.ProductsStatusesUpdateUsecase
	log logging.Logger
}

func NewProductsStatusesUpdateAction(uc usecase.ProductsStatusesUpdateUsecase, log logging.Logger) *ProductsStatusesUpdateAction {
	return &ProductsStatusesUpdateAction{
		uc:  uc,
		log: log,
	}
}

func (a *ProductsStatusesUpdateAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "product_statuses_update"

	var input usecase.ProductsStatusesUpdateInput
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

	products, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		statusCode := errorsStatus.StatusCode(err)
		logging.NewError(
			a.log,
			err,
			logKey,
			statusCode,
		).Log(utils.LogWithUser(r.Context(), "error when update product"))

		response.NewError(w, err, statusCode)
		return
	}

	response.NewSuccessList(w, products, len(products), http.StatusOK)
}
