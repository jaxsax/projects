package botv2_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuneIndexing(t *testing.T) {
	str := []rune("—11")

	require.Equal(t, []rune("—")[0], str[0])
}
