package config

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// reads environment variables directly
func LoadConfig() (err error) {

	_, b, _, _ := runtime.Caller(0)
	basePath := filepath.Dir(b)
	projectRootPath := filepath.Join(basePath, "../../..") // Adjust this as needed

	viper.AddConfigPath(filepath.Join(projectRootPath, "pkg/common/envs"))

	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	// Configuration found, but might be empty
	if err == nil {
		return
	}
	return fmt.Errorf("could not read the config file: %v", err)
}
