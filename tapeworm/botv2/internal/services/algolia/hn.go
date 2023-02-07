package algolia

type HN struct {
	*Algolia
}

func NewForHN() *HN {
	return &HN{
		Algolia: New("https://hn.algolia.com"),
	}
}
