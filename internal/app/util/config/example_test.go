package config_test

import (
	"os"
	"api/internal/app/util/config"
	config2 "api/internal/app/driver/config"
	"fmt"
)

func ExampleGet() {
	os.Setenv("EXAMPLE_KEY", "EXAMPLE_VALUE")

	c, _ := config2.NewEnvConfig()
	config.SetConfig(c)

	configValue, _ := config.Get("EXAMPLE_KEY")

	fmt.Println(configValue)

	// Output:
	// EXAMPLE_VALUE
}
