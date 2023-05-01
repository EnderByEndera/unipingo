package tests

import (
	"melodie-site/server/config"
	"os"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestReadConfig(t *testing.T) {
	rootPath := config.GetRootFilepath()
	assert.Equal(t, rootPath, "/root/backend")
	os.Chdir("../../../..")
	rootPath = config.GetRootFilepath()
	assert.Equal(t, rootPath, "/root/backend")
}
