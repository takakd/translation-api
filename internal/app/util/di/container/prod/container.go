package prod

import (
	"api/internal/app/driver/aws"
	"api/internal/app/driver/google"
	"api/internal/app/util/di"
)

// Container implements DI on production env.
type Container struct {
}

var _ di.DI = (*Container)(nil)

// Get returns interfaces corresponding name.
func (d *Container) Get(name string) (interface{}, error) {
	var (
		c   interface{}
		err error
	)

	if name == "translator.awsTextTranslator" {
		c, err = aws.NewTranslationAPI()
	} else if name == "translator.googleTextTranslator" {
		c, err = google.NewTranslationAPI()
	}

	return c, err
}
