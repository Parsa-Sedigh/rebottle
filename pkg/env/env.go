package env

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

const projectDirName = "rebottle"

func LoadEnv() error {
	re := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	cwd, _ := os.Getwd()
	rootPath := re.Find([]byte(cwd))

	_, b, _, _ := runtime.Caller(0)

	// TODO: Get the root of project without using relative paths here
	basepath := filepath.Dir(filepath.Join(b + "../../.."))

	fmt.Println("rootPath: ", string(rootPath), "cwd: ", cwd, "caller: ", basepath, "hello: ")

	err := godotenv.Load(basepath + `/.env`)
	fmt.Println("err in rootPath: ", err)
	if err != nil {
		return err
	}

	return nil
}
