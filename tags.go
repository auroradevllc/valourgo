package valourgo

type Tag struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Created Time   `json:"created"`
	Slug    string `json:"slug"`
}
