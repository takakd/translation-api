package local

import (
	"api/internal/app/controller/translator"
	"api/internal/app/driver/aws"
	"api/internal/app/driver/config"
	"api/internal/app/driver/google"
	"api/internal/app/driver/log"
	"api/internal/app/util/di"
	"context"
	"fmt"
)

// Container implements DI.
type Container struct {
}

var _ di.DI = (*Container)(nil)

// Get returns a concrete struct identified by name.
func (d *Container) Get(name string, args ...interface{}) (interface{}, error) {
	var (
		c   interface{}
		err error
	)

	if name == "driver.aws.TranslationAPI" {
		c, err = aws.NewTranslationAPI()

	} else if name == "driver.google.TranslationAPI" {
		c, err = google.NewTranslationAPI()

	} else if name == "util.config.Config" {
		files := make([]string, len(args))
		for i := range args {
			var ok bool
			files[i], ok = args[i].(string)
			if !ok {
				return nil, fmt.Errorf("argument error: %v", args)
			}
		}
		c, err = config.NewEnvConfig(files...)

	} else if name == "util.log.Logger" {
		c = log.NewStdoutLogger()

	} else if name == "driver.google.Client" {
		if len(args) < 2 {
			return nil, fmt.Errorf("argument error: %v", args)
		}

		ctx, ok := args[0].(context.Context)
		if !ok {
			return nil, fmt.Errorf("argument error: 1:ctx")
		}
		apiKey := args[1].(string)
		if !ok {
			return nil, fmt.Errorf("argument error: 2:apiKey")
		}

		c, err = google.NewClient(ctx, apiKey)
		if err != nil {
			err = fmt.Errorf("google.NewClient creation error: %w", err)
		}

	} else if name == "driver.aws.Translate" {
		c, err = aws.NewTranslate()

	} else if name == "controller.translator.Controller" {
		c, err = translator.NewController()

	}

	return c, err
}
