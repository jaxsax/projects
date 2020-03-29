package internal_test

import (
	"os"
	"strings"
	"testing"

	"github.com/jaxsax/projects/tapeworm/botv2/internal"
	"github.com/stretchr/testify/require"
)

func TestNoEmptyConfig(t *testing.T) {
	_, err := internal.ReadConfig(strings.NewReader(``))

	require.EqualError(t, err, internal.ErrEmptyConfig.Error())
}

func TestNoEmptyToken(t *testing.T) {
	_, err := internal.ReadConfig(strings.NewReader(`token:`))

	require.EqualError(t, err, internal.ErrEmptyToken.Error())
}

func TestPopulatesConfig(t *testing.T) {
	os.Setenv("WEB_PORT", "9999")
	config, err := internal.ReadConfig(strings.NewReader(`
token: aaa
`))

	require.NoError(t, err)
	require.Equal(t, &internal.Config{
		Token: "aaa",
		Port:  9999,
	}, config)
}
