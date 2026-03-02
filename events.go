package valourgo

import "time"

type ChannelStateEvent struct {
	ChannelID ChannelID `json:"channelId"`
	PlanetID  PlanetID  `json:"planetId"`
	Time      time.Time `json:"time"`
}

type ChannelWatchingUpdate struct {
	PlanetID  PlanetID  `json:"planetId"`
	ChannelID ChannelID `json:"channelId"`
	UserIDs   []UserID  `json:"userIds"`
}

type ChannelCurrentlyTypingUpdate struct {
	PlanetID  PlanetID  `json:"planetId"`
	ChannelID ChannelID `json:"channelId"`
	UserID    UserID    `json:"userId"`
}

type MessageCreateEvent struct {
	Message
}

type MessageEditEvent struct {
	Message
}

type MessageDeleteEvent struct {
	Message
}

type UserUpdateEvent struct {
	User
}

type PlanetMemberUpdate struct {
	Member
}

type MessageReactionEvent struct {
	MessageID MessageID `json:"messageId"`
	UserID    UserID    `json:"authorUserId"`
	MemberID  MemberID  `json:"authorMemberId"`
	Emoji     string    `json:"emoji"`
}

type MessageReactionAddedEvent struct {
	MessageReactionEvent
}

type MessageReactionRemovedEvent struct {
	MessageReactionEvent
}
