package valourgo

type Embed struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	StartPage            int    `json:"startPage"`
	HideChangePageArrows bool   `json:"hideChangePageArrows"
`
	Pages []EmbedPage `json:"pages"`
}

type EmbedItem struct {
	Children []EmbedItem `json:"children"`
}

type EmbedPage struct {
	EmbedItem
	Title  string `json:"title"`
	Footer string `json:"footer"`
}
