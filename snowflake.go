package valour

import (
	"strconv"
	"strings"
	"time"
)

// Epoch timestamp is 01/11/2021 4:37:00 UTC
const Epoch = 1610339820000 * time.Millisecond

const (
	generatorBits = 10
	sequenceBits  = 8

	lowerBits = generatorBits + sequenceBits
)

type Snowflake uint64

func (i Snowflake) String() string {
	return strconv.FormatInt(int64(i), 10)
}

func (i Snowflake) IsValid() bool {
	return i != 0
}

func (i Snowflake) Time() time.Time {
	unixnano := time.Duration(i>>lowerBits)*time.Millisecond + Epoch
	return time.Unix(0, int64(unixnano))
}

func (i Snowflake) Generator() uint16 {
	return uint16((i >> sequenceBits) & 0x3FF)
}

func (i Snowflake) Sequence() uint8 {
	return uint8(i & 0xFF)
}

type PlanetID Snowflake

func (i PlanetID) String() string {
	return Snowflake(i).String()
}

func (i PlanetID) IsValid() bool {
	return Snowflake(i).IsValid()
}

func (i PlanetID) Route(path ...string) string {
	p := []string{
		apiPlanetBase,
		i.String(),
	}

	p = append(p, path...)

	return strings.Join(p, "/")
}

type ChannelID Snowflake

func (i ChannelID) String() string {
	return Snowflake(i).String()
}

func (i ChannelID) IsValid() bool {
	return Snowflake(i).IsValid()
}

type UserID Snowflake

func (i UserID) String() string {
	return Snowflake(i).String()
}

func (i UserID) IsValid() bool {
	return Snowflake(i).IsValid()
}

type MemberID Snowflake

func (i MemberID) String() string {
	return Snowflake(i).String()
}

func (i MemberID) IsValid() bool {
	return Snowflake(i).IsValid()
}

type MessageID Snowflake

func (i MessageID) String() string {
	return Snowflake(i).String()
}

func (i MessageID) IsValid() bool {
	return Snowflake(i).IsValid()
}

func (i MessageID) Route(path ...string) string {
	p := []string{
		apiMessageBase,
		i.String(),
	}

	p = append(p, path...)

	return strings.Join(p, "/")
}

type RoleID Snowflake

func (i RoleID) String() string {
	return Snowflake(i).String()
}

func (i RoleID) IsValid() bool {
	return Snowflake(i).IsValid()
}

type EmojiID Snowflake

func (i EmojiID) String() string {
	return Snowflake(i).String()
}

func (i EmojiID) IsValid() bool {
	return Snowflake(i).IsValid()
}
