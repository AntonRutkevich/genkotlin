package kotlin

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/AntonRutkevich/astparser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	cfg := astparser.Config{
		InputDir:      "fixtures_test",
		IncludeRegexp: "test.go$",
	}
	sources, err := astparser.Load(cfg)
	require.NoError(t, err)

	generator := NewGenerator(Config{})
	files := generator.Generate(sources)
	for name, got := range files {
		want, err := ioutil.ReadFile(
			fmt.Sprintf("fixtures_test/%s.kt", name))
		require.NoError(t, err)
		assert.Equal(t, string(want), string(got))
	}
}
