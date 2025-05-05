package action

import (
	"glbackend/internal/adapters/api/response"
	"glbackend/internal/adapters/logging"
	"glbackend/internal/errorsStatus"
	"glbackend/internal/usecase"
	"glbackend/internal/utils"
	"net/http"
)

type TagsFindAllAction struct {
	uc  usecase.TagsFindAllUsecase
	log logging.Logger
}

func NewTagsFindAllAction(uc usecase.TagsFindAllUsecase, log logging.Logger) *TagsFindAllAction {
	return &TagsFindAllAction{
		uc:  uc,
		log: log,
	}
}

func (a *TagsFindAllAction) Execute(w http.ResponseWriter, r *http.Request) {
	const logKey = "tags_find_all"

	result, err := a.uc.Execute(r.Context())
	if err != nil {
		statusCode := errorsStatus.StatusCode(err)
		logging.NewError(
			a.log,
			err,
			logKey,
			statusCode,
		).Log(utils.LogWithUser(r.Context(), "error when tags find all"))

		response.NewError(w, err, statusCode)
		return
	}

	response.NewSuccessList(w, result, int(len(result)), http.StatusOK)
}
