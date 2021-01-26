package prod

import (
	"api/internal/app/driver/aws"
	"api/internal/app/driver/google"
	"api/internal/app/util/di"
	"api/internal/app/driver/config"
	"fmt"
	"api/internal/app/driver/log"
	"cloud.google.com/go/translate/apiv3"
	"google.golang.org/api/option"
	"context"
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

	if name == "translator.awsTextTranslator" {
		c, err = aws.NewTranslationAPI()

	} else if name == "translator.googleTextTranslator" {
		c, err = google.NewTranslationAPI()

	} else if name == "config.config" {
		files := make([]string, len(args))
		for i := range args {
			var ok bool
			files[i], ok = args[i].(string)
			if !ok {
				return nil, fmt.Errorf("argument error: %v", args)
			}
		}
		c, err = config.NewEnvConfig(files...)

	} else if name == "log.logger" {
		c = log.NewStdoutLogger()

	} else if name == "translate.NewTranslationClient" {
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

		// Ref: https://cloud.google.com/translate/docs/reference/rpc/google.cloud.translation.v3#google.cloud.translation.v3.TranslationService
		c, err = translate.NewTranslationClient(ctx, option.WithCredentialsJSON([]byte(apiKey)))
		if err != nil {
			err = fmt.Errorf("api initialize error: %w", err)
		}
	}

	return c, err
}
