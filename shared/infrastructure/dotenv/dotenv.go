package dotenv

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func Load() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for {
		_ = godotenv.Load()
		if isEmptyPath(dir) {
			break
		}

		dir = filepath.Dir(dir)
	}
}

func isEmptyPath(path string) bool {
	return path == "." || path == "/" || path == `\`
}
