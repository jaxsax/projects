package internal_test

import (
	"strings"
	"testing"

	"github.com/jaxsax/projects/tapeworm/botv2"
	"github.com/stretchr/testify/require"
)

func TestNoEmptyConfig(t *testing.T) {
	_, err := botv2.ReadConfig(strings.NewReader(``))

	require.EqualError(t, err, botv2.ErrEmptyConfig.Error())
}

func TestNoEmptyToken(t *testing.T) {
	_, err := botv2.ReadConfig(strings.NewReader(`token:`))

	require.EqualError(t, err, botv2.ErrEmptyToken.Error())
}

func TestPopulatesConfig(t *testing.T) {
	config, err := botv2.ReadConfig(strings.NewReader(`
token: aaa
`))

	require.NoError(t, err)
	require.Equal(t, &botv2.Config{
		Token: "aaa",
	}, config)
}
