package main

import (
	_ "strconv"
	_ "strings"
	"encoding/json"
	"github.com/baldurstod/vdf"
)

type color struct {
	ig *itemsGame
	Id string
	Name string
	Hex string
	kv *vdf.KeyValue
}

func (this *color) init(ig *itemsGame, kv *vdf.KeyValue) bool {
	this.ig = ig
	this.Id = kv.Key
	this.kv = kv

	return true
}

func (this *color) getStringAttribute(attributeName string) (string, bool) {
	if s, ok := this.kv.GetString(attributeName); ok {
		return s, true
	}
	return "", false
}

func (this *color) MarshalJSON() ([]byte, error) {
	ret := make(map[string]interface{})

	ret["id"] = this.Id

	if s, ok := this.getStringAttribute("color_name"); ok {
		ret["color_name"] = getStringToken(s)
	}

	if s, ok := this.getStringAttribute("hex_color"); ok {
		ret["hex_color"] = s
	}

	return json.Marshal(ret)
}
