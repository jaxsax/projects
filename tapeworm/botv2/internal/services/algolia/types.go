package algolia

import "time"

type SearchRequest struct {
	Query string
}

type SearchResponse struct {
	Hits []struct {
		CreatedAt       time.Time   `json:"created_at"`
		Title           string      `json:"title"`
		URL             string      `json:"url"`
		Author          string      `json:"author"`
		Points          int         `json:"points"`
		StoryText       interface{} `json:"story_text"`
		CommentText     interface{} `json:"comment_text"`
		NumComments     int         `json:"num_comments"`
		StoryID         interface{} `json:"story_id"`
		StoryTitle      interface{} `json:"story_title"`
		StoryURL        interface{} `json:"story_url"`
		ParentID        interface{} `json:"parent_id"`
		CreatedAtI      int         `json:"created_at_i"`
		Tags            []string    `json:"_tags"`
		ObjectID        string      `json:"objectID"`
		HighlightResult struct {
			Title struct {
				Value        string        `json:"value"`
				MatchLevel   string        `json:"matchLevel"`
				MatchedWords []interface{} `json:"matchedWords"`
			} `json:"title"`
			URL struct {
				Value            string   `json:"value"`
				MatchLevel       string   `json:"matchLevel"`
				FullyHighlighted bool     `json:"fullyHighlighted"`
				MatchedWords     []string `json:"matchedWords"`
			} `json:"url"`
			Author struct {
				Value        string        `json:"value"`
				MatchLevel   string        `json:"matchLevel"`
				MatchedWords []interface{} `json:"matchedWords"`
			} `json:"author"`
		} `json:"_highlightResult"`
	} `json:"hits"`
	NbHits           int  `json:"nbHits"`
	Page             int  `json:"page"`
	NbPages          int  `json:"nbPages"`
	HitsPerPage      int  `json:"hitsPerPage"`
	ExhaustiveNbHits bool `json:"exhaustiveNbHits"`
	ExhaustiveTypo   bool `json:"exhaustiveTypo"`
	Exhaustive       struct {
		NbHits bool `json:"nbHits"`
		Typo   bool `json:"typo"`
	} `json:"exhaustive"`
	Query               string `json:"query"`
	Params              string `json:"params"`
	ProcessingTimeMS    int    `json:"processingTimeMS"`
	ProcessingTimingsMS struct {
		Fetch struct {
			Total int `json:"total"`
		} `json:"fetch"`
		Request struct {
			RoundTrip int `json:"roundTrip"`
		} `json:"request"`
		Total int `json:"total"`
	} `json:"processingTimingsMS"`
	ServerTimeMS int `json:"serverTimeMS"`
}
