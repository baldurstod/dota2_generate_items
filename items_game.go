package main

import (
	_ "os"
	"fmt"
	_ "strings"
	_ "strconv"
	_ "encoding/json"
	//"github.com/baldurstod/vdf"
)

type itemMap map[string]*item
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

func (this *itemsGame) getItems() (*itemMap) {
	items := make(itemMap)

	for itemId, item := range this.Items {
		items[itemId] = item
	}
	return &items
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
