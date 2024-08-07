package main

import (
	"encoding/json"
	"strconv"

	"github.com/baldurstod/vdf"
)

type particle struct {
	ig *itemsGame
	Id string
	kv *vdf.KeyValue
}

func (p *particle) init(ig *itemsGame, kv *vdf.KeyValue) bool {
	p.ig = ig
	p.Id = kv.Key
	p.kv = kv

	return true
}

func (p *particle) getStringAttribute(attributeName string) (string, bool) {
	if s, ok := p.kv.GetString(attributeName); ok {
		return s, true
	}
	return "", false
}

func (p *particle) MarshalJSON() ([]byte, error) {
	ret := make(map[string]interface{})

	ret["id"], _ = strconv.Atoi(p.Id)

	for _, val := range p.kv.GetChilds() {
		ret[val.Key] = val.Value
	}

	return json.Marshal(ret)
}
