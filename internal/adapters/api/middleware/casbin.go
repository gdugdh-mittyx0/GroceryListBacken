package middleware

import (
	"fmt"
	"log"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
)

func NewCasbin() *casbin.Enforcer {
	modelFile, err := model.NewModelFromFile("configs/model/model.conf")
	if err != nil {
		log.Fatalf("Error loading model: %s", err)
	}

	adapterFile := fileadapter.NewAdapter("configs/model/policy.csv")
	enforcer, err := casbin.NewEnforcer(modelFile, adapterFile)
	if err != nil {
		panic(fmt.Sprintf("Error creating Enforcer: %s", err))
	}

	return enforcer
}
