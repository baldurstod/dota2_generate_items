package main

import (
	_ "os"
	"fmt"
	"strings"
	"strconv"
	"encoding/json"
	//"github.com/baldurstod/vdf"
)

type itemMap map[string]*item
type itemStyleMap map[string]*itemStyle
type itemGameMap map[string]interface{}
type itemStringMap map[string]string
type stringPair struct {
	s1 string
	s2 string
}
type collectionMap map[string]stringPair


type itemsGame struct {
	medals bool `default:false`
	itemsVDF *KeyValue
	Prefabs itemMap
	Items itemMap
	itemCollection collectionMap
}

func (this *itemsGame) MarshalJSON() ([]byte, error) {
	ret := make(itemGameMap)

	ret["items"] = *this.MarshalItems()
	ret["systems"] = *this.MarshalSystems()

	return json.Marshal(ret)
}

func (this *itemsGame) MarshalItems() *itemStyleMap {
	items := make(itemStyleMap)
	for itemId, item := range *this.getItems() {
		item.getUsedByHeroes()
		styles := item.getStyles()
		if len(styles) > 1 {
			for _, styleId := range styles {
				items[itemId + "~" + styleId] = &itemStyle{it: item, styleId: styleId}
			}
		} else {
			items[itemId] = &itemStyle{it: item, styleId: "0"}
		}
	}

	/*for itemId, itemData := range this.staticData {
		it := item{}
		if it.init(this, itemId, getMap(itemData)) {
			items[itemId] = &itemStyle{it: &it, styleId: "0"}
		}
	}*/

	return &items
}

func (this *itemsGame) MarshalSystems() *itemGameMap {
	systems := make(itemGameMap)

	if particles, ok := this.itemsVDF.Get("attribute_controlled_attached_particles"); ok {//getMap(getMap(this.itemsVDF)["attribute_controlled_attached_particles"])
		//fmt.Println(particles, ok)
		for _, val := range particles.GetChilds() {
			subParticles := val.GetChilds()
			for _, val := range subParticles {
				particleId := val.key
				systems[particleId] = val.GetChilds()

				if s, ok := getStringTokenRaw("Attrib_Particle" + particleId); ok {
					panic(s)//getMap(systems[particleId])["name"] = s
				} else {
					if s, ok := getStringTokenRaw("Attrib_KillStreakEffect" + particleId); ok {
						panic(s)//getMap(systems[particleId])["name"] = s
					}
				}
			}
		}
	}

	return &systems
}

func (this *itemsGame) getItems() (*itemMap) {
	items := make(itemMap)

	for itemId, item := range this.Items {
		if !this.filterOut(item, !this.medals) {
			items[itemId] = item
		}
	}
	return &items
}

func (this *itemsGame) filterOut(it *item, filterMedals bool) (bool) {
	return false
	it.initPrefabs()
	//fmt.Println(filterMedals)

	itemId, _ := strconv.Atoi(it.Id)

	if itemId == 5838 { //Winter 2015 Mystery Box
		return true
	}

	// Filter medals
	/*if s, ok := it.getStringAttribute("item_type_name"); ok {
		if filterMedals {
			if s == "#TF_Wearable_TournamentMedal" {
				return true
			}
		} else {
			if s != "#TF_Wearable_TournamentMedal" {
				return true
			}
		}
	}*/

	if s, ok := it.getStringAttribute("item_name"); ok {
		if s == "#TF_Item_Zombie_Armory" {
			return true
		}
	}

	if s, ok := it.getStringAttribute("name"); ok {
		if strings.Contains(s, "Autogrant") {
			return true
		}
	}

	if s, ok := it.getStringAttribute("item_class"); ok {
		if s == "tf_weapon_invis" {
			return true
		}
	}

	if s, ok := it.getStringAttribute("baseitem"); ok {
		if s == "1" {
			if itemId != 26 && itemId != 27 {//destruction PDA and disguise kit
				return true;
			}
		}
	}

	// Filter show_in_armory
	if s, ok := it.getStringAttribute("show_in_armory"); ok {
		if s == "0" {
			if s, ok := it.getStringAttribute("name"); ok && s != "Duck Badge" {
				if itemId == 294 || (itemId >= 831 && itemId <= 838) || (itemId >= 30739 && itemId <= 30741) {
				} else {
					return true
				}
			}
		}
	}

	// Filter items whitout models excepts a few (hatless...)
	if s, ok := it.getStringAttribute("model_player"); !ok || s == "" {
		if _, ok := it.getStringAttribute("model_player_per_class"); !ok {
			if _, ok := it.getStringAttribute("extra_wearable"); !ok {
				if itemSlot, ok := it.getStringAttribute("item_slot"); ok {
					if strings.Contains(itemSlot, "action") {
						return true
					}
				} else {
					return true
				}
			}
		}
	}

	usedByClasses := make(itemStringMap)
	it.getStringMapAttribute("used_by_classes", &usedByClasses)
	if len(usedByClasses) == 0 {
		return true
	}



	return false
}

func (this *itemsGame) init(dat []byte) {
	//this.staticData = make(itemGameMap)
	//dat, _ := os.ReadFile(path)
	/*if !this.medals {
		staticContent, _ := os.ReadFile(staticPath)
		_ = json.Unmarshal(staticContent, &this.staticData)
	}*/

	vdf := VDF{}
	root := vdf.Parse(dat)
	this.itemsVDF, _ = root.Get("items_game")// getMap(itemsVdf["items_game"]);
	//this.itemsVDF.Print()
	fmt.Println("")
	//fmt.Println(this.itemsVDF)
	//this.itemsVDF.Print()
	this.Prefabs = make(itemMap)
	this.Items = make(itemMap)
	this.itemCollection = make(collectionMap)

	if prefabs, ok := this.itemsVDF.Get("prefabs"); ok {//getMap(getMap(this.itemsVDF)["prefabs"])
	//prefabs.Print()
		for _, val := range prefabs.GetChilds() {
			var it = item{}
			if it.init(this, val) {
				this.Prefabs[it.Id] = &it
			}
		}
	}

	if items, ok := this.itemsVDF.Get("items"); ok {//getMap(getMap(this.itemsVDF)["items"])
		for _, val := range items.GetChilds() {
			var it = item{}
			//val.Print()
			if it.init(this, val) {
				this.Items[it.Id] = &it
			}
		}
	}

	for _, it := range this.Items {
		it.initPrefabs()
	}
/*
	itemCollections, itemCollectionsOk := this.itemsVDF.Get("item_collections")//getMap(getMap(this.itemsVDF)["item_collections"])

	fmt.Println(itemCollections, itemCollectionsOk)


	for _, val := range itemCollections.GetChilds() {
		collection := val//getMap(val)
		collectionItems, _ := collection.Get("items")//getMap(collection["items"])
		collectionName, _ := collection.GetString("name")//collection["name"].(string)

		for _, val := range collectionItems.GetChilds() {
			gradeName := val.key
			switch val.value.(type) {
			case map[string]interface{}:
				gradeItems := itemGameMap((val.value).(map[string]interface{}))

				for itemName, _ := range gradeItems {
					this.itemCollection[itemName] = stringPair{s1: collectionName, s2: gradeName}
				}

			}
		}
	}*/
}

func (this *itemsGame) getPrefab(prefabName string) *item {
	return this.Prefabs[prefabName]
}

func (this *itemsGame) getItemCollection(it *item) (string, string, bool) {
	if s, ok := it.getStringAttribute("name"); ok {
		if itemCollection, exists := this.itemCollection[s]; exists {
			return itemCollection.s1, itemCollection.s2, true
		}
	}
	return "", "", false
}

func (this *itemsGame) getItemsPerHero() map[string]*hero {
	heroes := make(map[string]*hero)

	for _, item := range this.Items {
		for _, npc := range item.getUsedByHeroes() {
			h, ok := heroes[npc]
			if !ok {
				h = &hero{npc: npc}
				heroes[npc] = h
			}

			h.addItem(item)
		}
	}
	return heroes
}
