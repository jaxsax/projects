package enhancers_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/enhancers"
	"github.com/stretchr/testify/require"
)

func TestTitleReader(t *testing.T) {
	body := `
		<html>
			<title>Hey</title>
		</html>
	`

	got, err := enhancers.ReadTitle(strings.NewReader(body))

	require.NoError(t, err)
	require.Equal(t, "Hey", got)
}

func TestNoTitleReader(t *testing.T) {
	body := `
		<html>
		</html>
	`

	got, err := enhancers.ReadTitle(strings.NewReader(body))

	require.Equal(t, "", got)
	require.Equal(t, fmt.Errorf("not found"), err)
}
