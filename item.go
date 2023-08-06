package main

import (
	"fmt"
	"strconv"
	"strings"
	"reflect"
	"encoding/json"
)

//var ITEM_FIELDS = [...]string{"image_inventory", /*"item_class",*/ "item_slot"/*, /*"item_type_name"/*, "item_quality"*/, "holiday_restriction", "anim_slot", "particle_suffix", "extra_wearable"}
//var ITEM_FIELDS = [...]string{"name", "image_inventory", "slot", "model_player", }

var ITEM_FIELDS = map[string]string{
	"image_inventory": "imageInventory",
	"model_player": "modelPlayer",
	"model_player1": "modelPlayer1",
	"model_player2": "modelPlayer2",
	"model_player3": "modelPlayer3",
	"item_slot": "slot",
	"baseitem": "baseItem",
	"item_rarity": "rarity",
	"workshop_accepted": "workshopAccepted",
}

type item struct {
	ig *itemsGame
	Id string
	Prefab string
	prefabs []*item
	isPrefabsInitialized bool `default:false`
	kv *KeyValue
}

func (this *item) init(ig *itemsGame, kv *KeyValue) bool {
	this.ig = ig
	this.Id = kv.key
	this.kv = kv

	return true
}

func (this *item) initPrefabs() {
	if !this.isPrefabsInitialized {
		this.isPrefabsInitialized = true
		if s, ok := this.kv.GetString("prefab"); ok {
			prefabs := strings.Split(s, " ")
			for _, prefabName := range prefabs {
				prefab := this.ig.getPrefab(prefabName)
				prefab.initPrefabs() //Ensure prefab is initialized
				this.prefabs = append(this.prefabs, prefab)
			}
		}
	}
}

func (this *item) getStringMapAttribute(attributeName string, i *itemStringMap) {
	if v, ok := this.kv.Get(attributeName); ok {
		switch v.value.(type) {
		case map[string]interface{}:
			mapValue := itemGameMap((v.value).(map[string]interface{}))
			for key, val := range mapValue {
				switch val.(type) {
				case string:
					(*i)[key] = val.(string)
				}
			}
		}
	}

	for _, prefab := range this.prefabs {
		prefab.getStringMapAttribute(attributeName, i)
	}
}

func (this *item) getStringMapSubAttribute(attributePath string, i *itemStringMap) {
	path := strings.Split(attributePath, ".")

	current := this.kv

ForLoop:
	for _, p := range path {
		next, ok := current.Get(p)
		if !ok {
			current = nil
			break ForLoop
		}
		switch next.value.(type) {
			case nil:
				current = nil
				break ForLoop
			case []*KeyValue:
				//current = (next.value).(*KeyValue)
				current = next
			case string:
				panic("Found a string " + attributePath + this.Id)
			default:
				fmt.Println(next)
				panic("Unknown type")
		}
		//current = (next.value).(*KeyValue)//itemGameMap((next).(map[string]interface{}))
	}

	for _, prefab := range this.prefabs {
		prefab.getStringMapSubAttribute(attributePath, i)
	}

	if current != nil {
		for key, val := range current.GetChilds() {
			switch val.value.(type) {
			case string:
				fmt.Println(key)
				panic("TODO")//(*i)[key] = val.value.(string)
			}
		}
	}

	return
}

func (this *item) getSubAttribute(attributePath string, i *itemGameMap) {
	path := strings.Split(attributePath, ".")

	current := this.kv

ForLoop:
	for _, p := range path {
		next, ok := current.Get(p)
		if !ok {
			current = nil
			break ForLoop
		}
		switch next.value.(type) {
			case nil:
				current = nil
				break ForLoop
			case []*KeyValue:
				//current = itemGameMap((next).(map[string]interface{}))
				current = next//(next.value).([]*KeyValue)
			case string:
				//return next.(string), true
				panic("Found a string " + attributePath + this.Id)
			default:
				fmt.Println(next.value, reflect.TypeOf(next.value))
				panic("Unknown type")
		}
		//current = (next.value).(*KeyValue)//current = itemGameMap((next).(map[string]interface{}))
	}
	for _, prefab := range this.prefabs {
		prefab.getSubAttribute(attributePath, i)
	}

	if current != nil {
		for key, val := range current.GetChilds() {
			switch val.value.(type) {
			case map[string]interface{}: //(*i)[key] = itemGameMap((val).(map[string]interface{}))
				fmt.Println(key)
				panic("TODO")
			case string: //(*i)[key] = val.(string)
				fmt.Println(key)
				panic("TODO")
			}
		}
	}

	return
}


func (this *item) getStringAttribute(attributeName string) (string, bool) {
	if s, ok := this.kv.GetString(attributeName); ok {
		return s, true
	}

	for _, prefab := range this.prefabs {
		if s, ok := prefab.getStringAttribute(attributeName); ok && s != "0" {//TODO: remove s != "0"
			return s, true
		}
	}
	return "", false
}

func (this *item) getStringSubAttribute(attributePath string) (string, bool) {
	path := strings.Split(attributePath, ".")

	current := this.kv
ForLoop:
	for _, p := range path {
		next, ok := current.Get(p)
		if !ok {
			break ForLoop
		}
		switch next.value.(type) {
			case nil:
				break ForLoop
			case []*KeyValue:
				//current = itemGameMap((next).(map[string]interface{}))
				current = next//(next.value).(*KeyValue)
			case string:
				return next.value.(string), true
			default:
				fmt.Println(next)
				panic("Unknown type")
		}
		//current = (next.value).(*KeyValue)//current = itemGameMap((next).(map[string]interface{}))
	}

	for _, prefab := range this.prefabs {
		if s, ok := prefab.getStringSubAttribute(attributePath); ok {
			return s, true
		}
	}

	return "", false
}

func (this *item) getUsedByHeroes() []string {
	ret := []string{}

	if usedByHeroes, ok := this.kv.GetStringMap("used_by_heroes"); ok {
		for key, val := range *usedByHeroes {
			if val == "1" {
				ret = append(ret, key)
			}
		}
		//fmt.Println(this.Id, ret)
	}
	return ret
}

func (this *item) MarshalJSON() ([]byte, error) {
	ret := make(map[string]interface{})

	ret["id"], _ = strconv.Atoi(this.Id)

	if s, ok := this.getStringAttribute("item_name"); ok {
		ret["name"] = getStringToken(s)
	}

	for inKey, outKey := range ITEM_FIELDS {
		if s, ok := this.getStringAttribute(inKey); ok {
			s = strings.ReplaceAll(s, "\\", "/")
			if inKey == "image_inventory" {
				s = strings.ToLower(s)
			}

			ret[outKey] = s
		}
	}

	this.MarshalVisuals(&ret)


	if bundle, ok := this.kv.Get("bundle"); ok {
		if sm, ok := bundle.ToStringMap(); ok {
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

func (this *item) MarshalVisuals(ret *map[string]interface{}) {

	modifiers := []interface{}{}
	//styles := []interface{}{}

	if visuals, ok := this.kv.Get("visuals"); ok {
		for _, kv := range visuals.value.([]*KeyValue) {
			if strings.HasPrefix(kv.key, "asset_modifier") {
				modifiers = append(modifiers, kv)
			}
			if strings.HasPrefix(kv.key, "styles") {
				//fmt.Println(kv.key)
				//styles = append(styles, kv)

				(*ret)["styles"] = kv
			}

			/*if val == "1" {
				ret = append(ret, key)
			}*/
		}
	}

	if len(modifiers) > 0 {
		(*ret)["assetmodifiers"] = modifiers
	}
	/*if len(styles) > 0 {
		(*ret)["styles"] = styles
	}*/

}

func (this itemGameMap) getMapStringValue(key string) (string, bool) {
	if mapValue := this[key]; mapValue != nil {
		switch mapValue.(type) {
			case string:
				return mapValue.(string), true
			default:
				return "", false
		}
	}
	return "", false
}

func (this *itemGameMap) getMapIntValue(key string) (int, bool) {
	if mapValue := (*this)[key]; mapValue != nil {
		if i, err := strconv.Atoi(mapValue.(string)); err == nil {
			return i, true
		}
	}
	return 0, false
}
