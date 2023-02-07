package dimension

import (
	"encoding/json"
	"fmt"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/services/algolia"
	"github.com/jaxsax/projects/tapeworm/botv2/internal/types"
)

type Collector interface {
	Collect(link *types.Link) ([]*types.Dimension, error)
}

type HNCollector struct {
	caller *algolia.HN
}

func NewHNCollector(api *algolia.HN) *HNCollector {
	return &HNCollector{
		caller: api,
	}
}

func (c *HNCollector) Collect(link *types.Link) ([]*types.Dimension, error) {
	response, err := c.caller.Search(&algolia.SearchRequest{
		Query: link.Link,
	})
	if err != nil {
		return nil, err
	}

	if len(response.Hits) == 0 {
		return nil, nil
	}

	hit := response.Hits[0]
	hnURL := fmt.Sprintf("https://news.ycombinator.com/item?id=%s", hit.ObjectID)

	data := map[string]string{
		"url": hnURL,
	}

	dataJson, err := json.Marshal(data)
	dim := &types.Dimension{
		Kind: types.DimensionHackernewsThread,
		Data: json.RawMessage(dataJson),
	}

	return []*types.Dimension{dim}, nil
}
