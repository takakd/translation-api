package config_test

import (
	config2 "api/internal/app/driver/config"
	"api/internal/app/util/config"
	"fmt"
	"os"
)

func ExampleGet() {
	os.Setenv("EXAMPLE_KEY", "EXAMPLE_VALUE")

	c, _ := config2.NewEnvConfig()
	config.SetConfig(c)

	configValue := config.Get("EXAMPLE_KEY")

	fmt.Println(configValue)

	// Output:
	// EXAMPLE_VALUE
}
