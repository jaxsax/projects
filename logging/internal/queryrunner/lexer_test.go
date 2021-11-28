package queryrunner

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestKeywords(t *testing.T) {
	query := `
    status_code:503 proto:https authority:clickhouse.internal.jaxsax.co
    `

	got, err := NewQuery(query)
	require.NoError(t, err)

	want := &RawQuery{
		Float64Fields: map[string]float64{
			"status_code": 503,
		},
		StringFields: map[string]string{
			"proto":     "https",
			"authority": "clickhouse.internal.jaxsax.co",
		},
	}

	require.Equal(t, want, got)
}
