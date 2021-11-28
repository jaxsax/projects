package queryrunner

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mustParse(layout, t string) time.Time {
	rt, err := time.Parse(layout, t)
	if err != nil {
		panic(err)
	}

	return rt
}

func TestBasicQuery(t *testing.T) {
	query := &Query{
		Raw:  "hello",
		From: mustParse(time.RFC3339, "2021-11-27T18:48:04.320Z"),
		To:   mustParse(time.RFC3339, "2021-11-27T18:58:04.320Z"),
		Table: &TableInfo{
			Database: "logging",
			Table:    "logs_v1",
		},
		ReturnFields: []string{"_timestamp", "message"},
	}

	got, err := query.GetMessageStatement()
	want := "SELECT _timestamp, message FROM logging.logs_v1 WHERE _timestamp >= 1638038884 AND _timestamp <= 1638039484 AND multiSearchAny(message, ['hello'])"

	require.NoError(t, err)
	assert.Equal(t, want, got)
}
