package hook

import (
	"fmt"
	"log"

	"github.com/achu-1612/glcm"
)

func Hook(args ...interface{}) error {
	if len(args) != 3 {
		return fmt.Errorf("invalid number of arguments")
	}

	name, ok := args[0].(string)
	if !ok {
		return fmt.Errorf("invalid argument type")
	}

	hType, ok := args[1].(string)
	if !ok {
		return fmt.Errorf("invalid argument type")
	}

	serviceName, ok := args[2].(string)
	if !ok {
		return fmt.Errorf("invalid argument type")
	}

	log.Println("hook ", name, " ", hType, " ", serviceName)

	return nil
}

func NewHookHandler(hookName, hookType, serviceName string) glcm.Hook {
	return glcm.NewHook(
		hookName,
		Hook,
		[]interface{}{hookName, hookType, serviceName}...,
	)
}
