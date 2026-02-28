package valourgo

import "encoding/json"

const latestEmbedVersion = "1.3"

type EmbedItemType int

const (
	EmbedTypeItem EmbedItemType = 1 + iota
	EmbedTypeText
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

type EmbedOption func(e *Embed)

func WithEmbedID(id string) EmbedOption {
	return func(e *Embed) {
		e.Id = &id
	}
}

func WithEmbedName(name string) EmbedOption {
	return func(e *Embed) {
		e.Name = &name
	}
}

func WithEmbedStartPage(page int) EmbedOption {
	return func(e *Embed) {
		e.StartPage = page
	}
}

func WithEmbedPages(pages ...EmbedPage) EmbedOption {
	return func(e *Embed) {
		e.Pages = append(e.Pages, pages...)
	}
}

// NewEmbed is recommended instead of creating Embeds directly as it sets some defaults
func NewEmbed(opts ...EmbedOption) *Embed {
	e := &Embed{
		Version:              latestEmbedVersion,
		HideChangePageArrows: true,
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

type Embed struct {
	Id                   *string     `json:"Id"`
	Name                 *string     `json:"Name"`
	StartPage            int         `json:"StartPage"`
	HideChangePageArrows bool        `json:"HideChangePageArrows"`
	Version              string      `json:"EmbedVersion"`
	Pages                []EmbedPage `json:"Pages"`
}

type EmbedItem interface {
	Type() EmbedItemType
}

type EmbedPage struct {
	Children []EmbedItem `json:"Children"`
	Title    string      `json:"Title"`
	Footer   string      `json:"Footer"`
}

func (e EmbedPage) Type() EmbedItemType {
	return EmbedTypeEmbedPage
}

func (e EmbedPage) MarshalJSON() ([]byte, error) {
	type Alias EmbedPage
	return json.Marshal(struct {
		Type EmbedItemType `json:"$type"`
		Alias
	}{
		Type:  e.Type(),
		Alias: Alias(e),
	})
}

type EmbedRow struct {
	Children []EmbedItem
}

func (e EmbedRow) Type() EmbedItemType {
	return EmbedTypeEmbedRow
}

func (e EmbedRow) MarshalJSON() ([]byte, error) {
	type Alias EmbedRow
	return json.Marshal(struct {
		Type EmbedItemType `json:"$type"`
		Alias
	}{
		Type:  e.Type(),
		Alias: Alias(e),
	})
}

type Clickable struct {
	ClickTarget any `json:"ClickTarget,omitempty"`
}

type EmbedText struct {
	Clickable
	Text string `json:"Text"`
}

func (e EmbedText) Type() EmbedItemType {
	return EmbedTypeText
}

func (e EmbedText) MarshalJSON() ([]byte, error) {
	type Alias EmbedText
	return json.Marshal(struct {
		Type EmbedItemType `json:"$type"`
		Alias
	}{
		Type:  e.Type(),
		Alias: Alias(e),
	})
}

type EmbedButton struct {
	Clickable
	Children []EmbedItem `json:"Children"`
}

func (e EmbedButton) Type() EmbedItemType {
	return EmbedTypeButton
}

func (e EmbedButton) MarshalJSON() ([]byte, error) {
	type Alias EmbedButton
	return json.Marshal(struct {
		Type EmbedItemType `json:"$type"`
		Alias
	}{
		Type:  e.Type(),
		Alias: Alias(e),
	})
}

type TargetType int

const (
	ClickTargetLink = 1 + iota
	ClickTargetPage
	ClickTargetEvent
	ClickTargetFormSubmit
)

type EmbedLinkTarget struct {
	Href string `json:"h"`
}

func (e EmbedLinkTarget) MarshalJSON() ([]byte, error) {
	type Alias EmbedLinkTarget
	return json.Marshal(struct {
		Type TargetType `json:"$type"`
		Alias
	}{
		Type:  ClickTargetLink,
		Alias: Alias(e),
	})
}

type EmbedPageTarget struct {
	Page int `json:"p"`
}

func (e EmbedPageTarget) MarshalJSON() ([]byte, error) {
	type Alias EmbedPageTarget
	return json.Marshal(struct {
		Type TargetType `json:"$type"`
		Alias
	}{
		Type:  ClickTargetPage,
		Alias: Alias(e),
	})
}

type EmbedEventTarget struct {
	EventElementID string `json:"e"`
}

func (e EmbedEventTarget) MarshalJSON() ([]byte, error) {
	type Alias EmbedEventTarget
	return json.Marshal(struct {
		Type TargetType `json:"$type"`
		Alias
	}{
		Type:  ClickTargetEvent,
		Alias: Alias(e),
	})
}

type EmbedFormSubmitTarget struct {
	EventElementID string `json:"e"`
}

func (e EmbedFormSubmitTarget) MarshalJSON() ([]byte, error) {
	type Alias EmbedFormSubmitTarget
	return json.Marshal(struct {
		Type TargetType `json:"$type"`
		Alias
	}{
		Type:  ClickTargetFormSubmit,
		Alias: Alias(e),
	})
}
