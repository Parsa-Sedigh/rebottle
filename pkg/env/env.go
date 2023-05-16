package env

import (
	"github.com/joho/godotenv"
	"path/filepath"
	"runtime"
)

const projectDirName = "rebottle"

func LoadEnv() error {
	_, b, _, _ := runtime.Caller(0)

	// TODO: Get the root of project without using relative paths here
	basepath := filepath.Dir(filepath.Join(b + "../../.."))

	err := godotenv.Load(basepath + `/.env`)
	if err != nil {
		return err
	}

	return nil
}
