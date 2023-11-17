package src

import (
	"encoding/json"
	"io"
	"math/rand"
	"os"
	"reflect"
	"time"
)

func FromJson(filepath string, to any) {
	jsonFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(bytes, to); err != nil {
		panic(err)
	}
}

func shuffleSlice(s interface{}) {
	rv := reflect.ValueOf(s)
	swap := reflect.Swapper(s)
	seed := rand.New(rand.NewSource(time.Now().UnixNano()))
	seed.Shuffle(rv.Len(), swap)
}
