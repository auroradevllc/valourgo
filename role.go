package valourgo

type Role struct {
	ID                  RoleID   `json:"id"`
	PlanetID            PlanetID `json:"planetId"`
	Name                string   `json:"name"`
	Position            int      `json:"position"`
	IsDefault           bool     `json:"isDefault"`
	Permissions         int64    `json:"permissions"`
	ChatPermissions     int      `json:"chatPermissions"`
	CategoryPermissions int      `json:"categoryPermissions"`
	VoicePermissions    int      `json:"voicePermissions"`
	Color               string   `json:"color"`
	Bold                bool     `json:"bold"`
	Italics             bool     `json:"italics"`
	FlagBitIndex        int      `json:"flagBitIndex"`
	AnyoneCanMention    bool     `json:"anyoneCanMention"`
	IsAdmin             bool     `json:"isAdmin"`
}
