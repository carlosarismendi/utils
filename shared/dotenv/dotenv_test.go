package dotenv

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAlgo_shit(t *testing.T) {
	// ARRANGE
	currentDir, err := os.Getwd()
	require.NoError(t, err)

	envFileName := ".env.test"
	file, err := ioutil.TempFile(currentDir, envFileName)
	require.NoError(t, err)
	defer os.Remove(file.Name())

	_, err = file.Write([]byte("TESTING_ENV_LOAD=ok"))
	require.NoError(t, err)

	Load(file.Name())
	value := os.Getenv("TESTING_ENV_LOAD")
	require.NotEmpty(t, value)
	require.Equal(t, "ok", value)
}
