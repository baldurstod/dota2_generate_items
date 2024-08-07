package main

import (
	"encoding/json"
	_ "strconv"
	_ "strings"

	"github.com/baldurstod/vdf"
)

type color struct {
	ig   *itemsGame
	Id   string
	Name string
	Hex  string
	kv   *vdf.KeyValue
}

func (c *color) init(ig *itemsGame, kv *vdf.KeyValue) bool {
	c.ig = ig
	c.Id = kv.Key
	c.kv = kv

	return true
}

func (c *color) getStringAttribute(attributeName string) (string, bool) {
	if s, ok := c.kv.GetString(attributeName); ok {
		return s, true
	}
	return "", false
}

func (c *color) MarshalJSON() ([]byte, error) {
	ret := make(map[string]interface{})

	ret["id"] = c.Id

	if s, ok := c.getStringAttribute("color_name"); ok {
		ret["color_name"] = getStringToken(s)
	}

	if s, ok := c.getStringAttribute("hex_color"); ok {
		ret["hex_color"] = s
	}

	return json.Marshal(ret)
}
