package app

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv(envVarKeys []string, fileBasePath string) error {
	for _, k := range envVarKeys {
		str := os.Getenv(k)

		if str == "" {
			continue
		}

		r := strings.NewReader(str)
		envMap, err := godotenv.Parse(r)
		if err != nil {
			return err
		}

		for ik, iv := range envMap {
			if _, exists := os.LookupEnv(ik); exists {
				_ = os.Setenv(ik, iv)
			}
		}
	}

	for _, k := range envVarKeys {
		envPath := filepath.Join(fileBasePath, k+".env")
		_, err := os.Stat(envPath)
		if os.IsNotExist(err) {
			path, err := os.Getwd()
			if err != nil {
				return err
			}
			envPath = filepath.Join(path, envPath)
		}
		err = godotenv.Load(envPath)
		if err != nil {
			pathErr, _ := err.(*os.PathError)
			if pathErr == nil {
				return err
			}
		}
	}

	return nil
}
