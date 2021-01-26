package prod

import (
	"api/internal/app/driver/aws"
	"api/internal/app/driver/google"
	"api/internal/app/util/di"
	"api/internal/app/driver/config"
	"fmt"
	"api/internal/app/driver/log"
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
	}

	return c, err
}
