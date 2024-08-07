package main

import (
	"encoding/json"
)

type hero struct {
	npc   string
	items []*item
}

func (h *hero) MarshalJSON() ([]byte, error) {
	ret := []interface{}{}

	for _, val := range h.items {
		ret = append(ret, val)
	}

	return json.Marshal(ret)
}

func (h *hero) addItem(i *item) {
	h.items = append(h.items, i)
}
