package main

import (
	"encoding/json"
)

type hero struct {
	npc string
	items []*item
}

func (this *hero) MarshalJSON() ([]byte, error) {
	ret := []interface{}{}

	for _, val := range this.items {
		ret = append(ret, val)
	}

	return json.Marshal(ret)
}

func (this *hero) addItem(i *item) {
	this.items = append(this.items, i)
}
