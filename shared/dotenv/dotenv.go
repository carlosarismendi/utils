package dotenv

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

var defaultFilenames = []string{".env", ".env.local"}

func Load(filenames ...string) {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if len(filenames) == 0 {
		filenames = defaultFilenames
	}

	for {
		for i := range filenames {
			path := fmt.Sprintf("%s/%s", dir, filenames[i])
			_ = godotenv.Load(path)
		}

		if isEmptyPath(dir) {
			break
		}

		dir = filepath.Dir(dir)
	}
}

func isEmptyPath(path string) bool {
	return path == "." || path == "/" || path == `\`
}
