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

type TagsDeleteAction struct {
	uc  usecase.TagsDeleteUsecase
	log logging.Logger
}

func NewTagsDeleteAction(uc usecase.TagsDeleteUsecase, log logging.Logger) *TagsDeleteAction {
	return &TagsDeleteAction{
		uc:  uc,
		log: log,
	}
}

func (a *TagsDeleteAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "tags_delete"

	var input usecase.TagsDeleteParams
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
		).Log(utils.LogWithUser(r.Context(), "error when delete tag"))

		response.NewError(w, err, statusCode)
		return
	}

	response.NewSuccess(w, true, http.StatusOK)
}
