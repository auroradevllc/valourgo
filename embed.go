package valourgo

import "encoding/json"

type EmbedItemType int

const (
	EmbedTypeText EmbedItemType = 1 + iota
	EmbedTypeButton
	EmbedTypeInputBox
	EmbedTypeTextArea
	EmbedTypeProgressBar
	EmbedTypeForm
	EmbedTypeGoTo
	EmbedTypeDropDownItem
	EmbedTypeDropDownMenu
	EmbedTypeEmbedRow
	EmbedTypeEmbedPage
	EmbedTypeProgress
	EmbedTypeMedia
)

type Embed struct {
	Id                   string `json:"id"`
	Name                 string `json:"name"`
	StartPage            int    `json:"startPage"`
	HideChangePageArrows bool   `json:"hideChangePageArrows"
`
	Pages []EmbedPage `json:"pages"`
}

type EmbedItem interface {
	Type() EmbedItemType
}

type EmbedPage struct {
	Children []EmbedItem `json:"children"`
	Title    string      `json:"title"`
	Footer   string      `json:"footer"`
}

func (e EmbedPage) Type() EmbedItemType {
	return EmbedTypeEmbedPage
}

func (e EmbedPage) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":     EmbedTypeEmbedPage,
		"children": e.Children,
		"title":    e.Title,
		"footer":   e.Footer,
	})
}

type EmbedRow struct {
	Children []EmbedItem
}

func (e EmbedRow) Type() EmbedItemType {
	return EmbedTypeEmbedRow
}

func (e EmbedRow) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type":     EmbedTypeEmbedRow,
		"children": e.Children,
	})
}

type EmbedText struct {
	Text string `json:"text"`
}

func (e EmbedText) Type() EmbedItemType {
	return EmbedTypeText
}

func (e EmbedText) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]any{
		"type": EmbedTypeText,
		"text": e.Text,
	})
}
