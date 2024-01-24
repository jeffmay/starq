package starq_test

import (
	"fmt"
	"os"
	"starq/internal/starq"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/valyala/fastjson"
)

const petstoreSpecJSONFile = "./fake/input/petstore-openapi.json"

var petstoreSpecJSON string

func loadJSONPetstoreOpenAPI() string {
	if len(petstoreSpecJSON) == 0 {
		bytes, err := os.ReadFile(petstoreSpecJSONFile)
		if err != nil {
			panic(fmt.Errorf("could not load '%s': %w", petstoreSpecJSONFile, err))
		}
		petstoreSpecJSON = string(bytes)
	}
	return petstoreSpecJSON
}

func TestSimpleConfig(t *testing.T) {
	opts := MakeTestOpts(
		WithConfigFile("./fake/config/simple.yaml"),
	)
	err := starq.Run(opts.Opts())
	require.NoError(t, err)
}

func TestPetstore(t *testing.T) {
	opts := MakeTestOpts(
		PrependRules(".info.title"),
		WithInputString(loadJSONPetstoreOpenAPI()),
	)
	err := starq.Run(opts.Opts())
	require.NoError(t, err)
	title := opts.MustOutputValue()
	require.Equal(t, `"Swagger Petstore"`, title.String())
}

func TestPetstoreReadonly(t *testing.T) {
	opts := MakeTestOpts(
		WithConfigFile("./fake/config/petstore-readonly.yaml"),
		WithInputString(loadJSONPetstoreOpenAPI()),
	)
	err := starq.Run(opts.Opts())
	require.NoError(t, err, opts.Errors())
	doc := opts.MustOutputValue()
	paths := doc.GetObject("paths")
	require.NotNilf(t, paths, ".paths should be defined in:\n%s", opts.Pretty())
	require.Equalf(t, paths.Len(), 2, ".paths should have length 2 in:\n%s", opts.Pretty())
	paths.Visit(func(key []byte, path *fastjson.Value) {
		getRoute := path.GetObject("get")
		require.NotNilf(t, getRoute, ".paths[\"%s\"].get should be defined in:\n%s", string(key), opts.Pretty())
		postRoute := path.GetObject("post")
		require.Nilf(t, postRoute, ".paths[\"%s\"].post should NOT be defined in:\n%s", string(key), opts.Pretty())
	})
	components := doc.GetObject("components")
	require.NotNilf(t, components, ".components should be defined in:\n%s", opts.Pretty())
}
