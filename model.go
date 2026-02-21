package valourgo

import "time"

type User struct {
	HasCustomAvatar   bool        `json:"hasCustomAvatar"`
	HasAnimatedAvatar bool        `json:"hasAnimatedAvatar"`
	TimeJoined        time.Time   `json:"timeJoined"`
	Name              string      `json:"name"`
	Tag               string      `json:"tag"`
	Bot               bool        `json:"bot"`
	Disabled          bool        `json:"disabled"`
	ValourStaff       bool        `json:"valourStaff"`
	Status            *string     `json:"status"`
	UserStateCode     int         `json:"userStateCode"`
	TimeLastActive    time.Time   `json:"timeLastActive"`
	IsMobile          bool        `json:"isMobile"`
	Compliance        bool        `json:"compliance"`
	SubscriptionType  string      `json:"subscriptionType"`
	PriorName         string      `json:"priorName"`
	NameChangeTime    time.Time   `json:"nameChangeTime"`
	Version           int         `json:"version"`
	TutorialState     int         `json:"tutorialState"`
	OwnerId           interface{} `json:"ownerId"`
	NameAndTag        string      `json:"nameAndTag"`
	ID                UserID      `json:"id"`
}

type Member struct {
	ID       MemberID `json:"id"`
	User     User     `json:"user"`
	UserID   UserID   `json:"userId"`
	PlanetID PlanetID `json:"planetId"`
	Nickname *string  `json:"nickname"`
	Avatar   *string  `json:"memberAvatar"`
}

type Reaction struct {
	ID             int64     `json:"id"`
	Emoji          string    `json:"emoji"`
	MessageID      MessageID `json:"messageId"`
	AuthorUserID   UserID    `json:"authorUserId"`
	AuthorMemberID MemberID  `json:"authorMemberId"`
	CreatedAt      time.Time `json:"createdAt"`
}

type ChannelType int

type Channel struct {
	ID             ChannelID       `json:"id"`
	PlanetID       PlanetID        `json:"planetId"`
	ParentID       ChannelID       `json:"parentId"`
	ChannelType    ChannelType     `json:"channelType"`
	Name           string          `json:"name"`
	Description    string          `json:"description"`
	InheritsPerms  bool            `json:"inheritsPerms"`
	IsDefault      bool            `json:"isDefault"`
	RawPosition    int64           `json:"rawPosition"`
	Position       ChannelPosition `json:"position"`
	LastUpdateTime time.Time       `json:"lastUpdateTime"`
}

type ChannelPosition struct {
	RawPosition   int64 `json:"rawPosition"`
	Depth         int   `json:"depth"`
	LocalPosition int   `json:"localPosition"`
}
