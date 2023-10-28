package dotenv

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Priority from left (highest) to right (lowest)
var defaultFilenames = []string{".env.local", ".env.production", ".env"}

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
