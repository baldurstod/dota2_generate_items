package main

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/baldurstod/vdf"
)

var ITEM_FIELDS = map[string]string{
	"image_inventory":   "imageInventory",
	"model_player":      "modelPlayer",
	"model_player1":     "modelPlayer1",
	"model_player2":     "modelPlayer2",
	"model_player3":     "modelPlayer3",
	"item_slot":         "slot",
	"baseitem":          "baseItem",
	"item_rarity":       "rarity",
	"workshop_accepted": "workshopAccepted",
}

type item struct {
	ig                   *itemsGame
	Id                   string
	Prefab               string
	prefabs              []*item
	isPrefabsInitialized bool
	kv                   *vdf.KeyValue
}

func (i *item) init(ig *itemsGame, kv *vdf.KeyValue) bool {
	i.ig = ig
	i.Id = kv.Key
	i.kv = kv

	return true
}

func (i *item) initPrefabs() {
	if !i.isPrefabsInitialized {
		i.isPrefabsInitialized = true
		if s, err := i.kv.GetString("prefab"); err == nil {
			prefabs := strings.Split(s, " ")
			for _, prefabName := range prefabs {
				prefab := i.ig.getPrefab(prefabName)
				prefab.initPrefabs() //Ensure prefab is initialized
				i.prefabs = append(i.prefabs, prefab)
			}
		}
	}
}

func (i *item) getStringAttribute(attributeName string) (string, bool) {
	if s, err := i.kv.GetString(attributeName); err == nil {
		return s, true
	}

	for _, prefab := range i.prefabs {
		if s, ok := prefab.getStringAttribute(attributeName); ok && s != "0" { //TODO: remove s != "0"
			return s, true
		}
	}
	return "", false
}

func (i *item) getUsedByHeroes() []string {
	ret := []string{}

	if usedByHeroes, err := i.kv.GetStringMap("used_by_heroes"); err == nil {
		for key, val := range *usedByHeroes {
			if val == "1" {
				ret = append(ret, key)
			}
		}
	} else {
		if prefab, err := i.kv.GetString("prefab"); err == nil {
			ret = append(ret, prefab)
		}
	}
	return ret
}

func (i *item) MarshalJSON() ([]byte, error) {
	ret := make(map[string]interface{})

	ret["id"], _ = strconv.Atoi(i.Id)

	if s, ok := i.getStringAttribute("item_name"); ok {
		ret["name"] = getStringToken(s)
	}

	for inKey, outKey := range ITEM_FIELDS {
		if s, ok := i.getStringAttribute(inKey); ok {
			s = strings.ReplaceAll(s, "\\", "/")
			if inKey == "image_inventory" {
				s = strings.ToLower(s)
			}

			ret[outKey] = s
		}
	}

	i.MarshalVisuals(&ret)
	i.MarshalStaticAttributes(&ret)

	if bundle, err := i.kv.Get("bundle"); err == nil {
		if sm, err := bundle.ToStringMap(); err == nil {
			items := []string{}
			for key, val := range *sm {
				if val == "1" {
					items = append(items, key)
				}
			}
			ret["bundle"] = items
		}
	}

	return json.Marshal(ret)
}

func (i *item) MarshalVisuals(ret *map[string]interface{}) {

	modifiers := []interface{}{}

	if visuals, err := i.kv.GetAll("visuals"); err == nil {
		for _, kv := range visuals {
			for _, kv := range kv.GetChilds() {
				if strings.HasPrefix(kv.Key, "asset_modifier") {
					modifiers = append(modifiers, kv)
				}
				if strings.HasPrefix(kv.Key, "styles") {
					marshalStyles(kv, ret)
				}
				if strings.HasPrefix(kv.Key, "skin") {
					(*ret)["skin"] = kv
				}
				if strings.HasPrefix(kv.Key, "hide_styles_from_ui") {
					(*ret)["hide_styles_from_ui"] = kv
				}
			}
		}
	}

	if len(modifiers) > 0 {
		(*ret)["assetmodifiers"] = modifiers
	}
}

func (i *item) MarshalStaticAttributes(ret *map[string]interface{}) {
	staticAttributes := map[string]interface{}{}

	if attributes, err := i.kv.GetAll("static_attributes"); err == nil {
		for _, kv := range attributes {
			for _, kv := range kv.GetChilds() {
				staticAttributes[kv.Key] = kv
				//staticAttributes = append(staticAttributes, map[string]interface{}{"key": kv.Key, "value": kv.Value})
			}
		}
	}

	if len(staticAttributes) > 0 {
		(*ret)["static_attributes"] = staticAttributes
	}
}

func marshalStyles(kvStyles *vdf.KeyValue, ret *map[string]interface{}) {
	styles := make(map[string]interface{})

	for _, kv := range kvStyles.GetChilds() {
		marshalStyle(kv, &styles)
	}

	(*ret)["styles"] = styles
}

func marshalStyle(kvStyle *vdf.KeyValue, ret *map[string]interface{}) {
	style := make(map[string]interface{})

	for _, kv := range kvStyle.GetChilds() {
		if kv.Key == "name" {
			if name, err := kv.ToString(); err == nil {
				style["name"] = getStringToken(name)
			}
		} else {
			style[kv.Key] = kv.GetValue()
		}
	}

	(*ret)[kvStyle.Key] = style
}
