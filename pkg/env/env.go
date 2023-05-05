package env

import (
	"github.com/joho/godotenv"
	"os"
	"regexp"
)

const projectDirName = "rebottle"

func LoadEnv() error {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	err := godotenv.Load(string(rootPath) + `/.env`)
	if err != nil {
		return err
	}

	return nil
}
