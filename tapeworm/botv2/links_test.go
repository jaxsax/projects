package botv2

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDBObjectToViewObject(t *testing.T) {
	l := link{
		ID:        1,
		CreatedBy: 1,
		CreatedTS: 1,
		Link:      "qweqw",
		Title:     "qweqw",
		ExtraData: "{\"a\": 1}",
	}

	want := Link{
		ID:        1,
		CreatedBy: 1,
		CreatedTS: 1,
		Link:      "qweqw",
		Title:     "qweqw",
		ExtraData: map[string]interface{}{"a": float64(1)},
	}

	require.Equal(t, want, fromDBObject(l))
}
