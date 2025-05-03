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

type ProductsCreateAction struct {
	uc  usecase.ProductsCreateUsecase
	log logging.Logger
}

func NewProductsCreateAction(uc usecase.ProductsCreateUsecase, log logging.Logger) *ProductsCreateAction {
	return &ProductsCreateAction{
		uc:  uc,
		log: log,
	}
}

func (a *ProductsCreateAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "product_create"

	var input usecase.ProductsCreateInput
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

	product, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		statusCode := errorsStatus.StatusCode(err)
		logging.NewError(
			a.log,
			err,
			logKey,
			statusCode,
		).Log(utils.LogWithUser(r.Context(), "error when create product"))

		response.NewError(w, err, statusCode)
		return
	}

	response.NewSuccess(w, product, http.StatusOK)
}
