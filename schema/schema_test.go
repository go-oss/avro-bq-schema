package schema

import (
	"embed"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDataDir       = "testdata"
	avroSchemaFileExt = ".avsc"
	bqSchemaFileExt   = ".json"
)

var (
	//go:embed testdata
	testFiles embed.FS
)

func TestConvert(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		indent int
	}{
		{
			name:   "basic",
			indent: 2,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			avroSchema, err := testFiles.ReadFile(filepath.Join(testDataDir, tt.name+avroSchemaFileExt))
			require.NoError(t, err)
			bqSchema, err := testFiles.ReadFile(filepath.Join(testDataDir, tt.name+bqSchemaFileExt))
			require.NoError(t, err)
			bq, err := Convert(avroSchema)
			require.NoError(t, err)
			got, err := ToJSON(bq, tt.indent)
			require.NoError(t, err)
			assert.Equal(t, string(bqSchema), string(got))
		})
	}
}
