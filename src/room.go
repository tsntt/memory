package src

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
)

var RoomPathID = make(map[string]uuid.UUID, 0)

type Room struct {
	Path    Path                  `json:"id"`
	ConnLen Players               `json:"-"`
	Clients map[uuid.UUID]*Client `json:"clients"`
	Cards   []Card                `json:"cards"`
	Turn    uuid.UUID             `json:"-"`
}

func NewRoom(p Path, n Players) *Room {
	return &Room{
		Path:    p,
		ConnLen: n,
		Clients: make(map[uuid.UUID]*Client, n.Int()),
	}
}

func (r *Room) NewGame(cards []string) *Room {
	shuffleSlice(cards)

	r.Cards = MakeCardPairs(cards[:10]...)

	return r
}

func (r *Room) ClientTurn() uuid.UUID {
	next := false

	for _, c := range r.Clients {
		if r.Turn == (uuid.UUID{}) {
			r.Turn = c.ID
			break
		}
		if r.Turn == c.ID {
			next = true
			continue
		}
		if next {
			r.Turn = c.ID
			next = false
			break
		}
	}

	return r.Turn
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
