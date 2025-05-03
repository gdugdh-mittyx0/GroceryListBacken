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

type TagsCreateAction struct {
	uc  usecase.TagsCreateUsecase
	log logging.Logger
}

func NewTagsCreateAction(uc usecase.TagsCreateUsecase, log logging.Logger) *TagsCreateAction {
	return &TagsCreateAction{
		uc:  uc,
		log: log,
	}
}

func (a *TagsCreateAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "tag_create"

	var input usecase.TagsCreateInput
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

	tag, err := a.uc.Execute(r.Context(), input)
	if err != nil {
		statusCode := errorsStatus.StatusCode(err)
		logging.NewError(
			a.log,
			err,
			logKey,
			statusCode,
		).Log(utils.LogWithUser(r.Context(), "error when create tag"))

		response.NewError(w, err, statusCode)
		return
	}

	response.NewSuccess(w, tag, http.StatusOK)
}
