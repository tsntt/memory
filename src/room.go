package src

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
)

var RoomIDToPath = make(map[uuid.UUID]string, 0)

type Room struct {
	ID      uuid.UUID          `json:"id"`
	Path    Path               `json:"-"`
	ConnLen Players            `json:"-"`
	Clients map[string]*Client `json:"clients"`
	Cards   []Card             `json:"cards"`
}

func NewRoom(p Path, n Players) *Room {
	return &Room{
		ID:      uuid.New(),
		Path:    p,
		ConnLen: n,
		Clients: make(map[string]*Client, n.Int()),
	}
}

func (r *Room) NewGame(cards []string) *Room {
	shuffleSlice(cards)

	r.Cards = MakeCardPairs(cards[:10]...)

	return r
}

type Path string

func NewPath(p string) (Path, error) {
	// validate p before returning
	// accept only [a-Z,0-9,-,_]
	return Path(p), nil
}

func (p Path) String() string { return string(p) }

type Players int

func NewPlayers(n string) (Players, error) {
	num, err := strconv.Atoi(n)
	if err != nil {
		return Players(0), errors.New("number of players is not a number")
	}

	if num < 2 {
		return Players(0), errors.New("you need aleast two players")
	}

	if num > 4 {
		return Players(0), errors.New("max of 4 players per room")
	}

	return Players(num), nil
}

func (p Players) Int() int { return int(p) }
