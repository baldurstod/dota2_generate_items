package main

import (
	"strconv"
	"encoding/json"
	"github.com/baldurstod/vdf"
)

type particle struct {
	ig *itemsGame
	Id string
	kv *vdf.KeyValue
}

func (this *particle) init(ig *itemsGame, kv *vdf.KeyValue) bool {
	this.ig = ig
	this.Id = kv.Key
	this.kv = kv

	return true
}

func (this *particle) getStringAttribute(attributeName string) (string, bool) {
	if s, ok := this.kv.GetString(attributeName); ok {
		return s, true
	}
	return "", false
}

func (this *particle) MarshalJSON() ([]byte, error) {
	ret := make(map[string]interface{})

	ret["id"], _ = strconv.Atoi(this.Id)

	for _, val := range this.kv.GetChilds() {
		ret[val.Key] = val.Value
	}

	return json.Marshal(ret)
}
