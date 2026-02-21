package valourgo

import (
	"strconv"
	"strings"
)

type ValourID int64

func (i ValourID) String() string {
	return strconv.FormatInt(int64(i), 10)
}

type PlanetID ValourID

func (i PlanetID) String() string {
	return ValourID(i).String()
}

func (i PlanetID) Route(path ...string) string {
	p := []string{
		apiPlanetBase,
		i.String(),
	}

	p = append(p, path...)

	return strings.Join(p, "/")
}

type ChannelID ValourID

func (i ChannelID) String() string {
	return ValourID(i).String()
}

type UserID ValourID

func (i UserID) String() string {
	return ValourID(i).String()
}

type MemberID ValourID

func (i MemberID) String() string {
	return ValourID(i).String()
}

type MessageID ValourID

func (i MessageID) String() string {
	return ValourID(i).String()
}

func (i MessageID) Route(path ...string) string {
	p := []string{
		apiMessageBase,
		i.String(),
	}

	p = append(p, path...)

	return strings.Join(p, "/")
}

type RoleID ValourID

func (i RoleID) String() string {
	return ValourID(i).String()
}

type EmojiID ValourID

func (i EmojiID) String() string {
	return ValourID(i).String()
}
