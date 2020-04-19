package enhancers

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUTMStripper(t *testing.T) {
	link := "https://www.reddit.com/r/conspiracy/comments/g3gc0v/the_mandela_effect_and_cern_in_depth_analysis/?utm_medium=android_app&utm_source=share"

	url, err := url.Parse(link)
	require.NoError(t, err)

	removeUTMParameters(url)
	require.Equal(
		t,
		"https://www.reddit.com/r/conspiracy/comments/g3gc0v/the_mandela_effect_and_cern_in_depth_analysis/",
		url.String(),
	)
}
