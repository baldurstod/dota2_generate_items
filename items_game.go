package main

import (
	"github.com/baldurstod/vdf"
)

type itemMap map[string]*item
type colorMap map[string]*color
type particleMap map[string]*particle

type itemsGame struct {
	medals    bool
	itemsVDF  *vdf.KeyValue
	Prefabs   itemMap
	Items     itemMap
	Colors    colorMap
	Particles particleMap
}

func (ig *itemsGame) getItems() *itemMap {
	items := make(itemMap)

	for itemId, item := range ig.Items {
		items[itemId] = item
	}
	return &items
}

func (ig *itemsGame) init(dat []byte) {
	vdf := vdf.VDF{}
	root := vdf.Parse(dat)
	ig.itemsVDF, _ = root.Get("items_game")
	ig.Prefabs = make(itemMap)
	ig.Items = make(itemMap)
	ig.Colors = make(colorMap)
	ig.Particles = make(particleMap)

	if prefabs, ok := ig.itemsVDF.Get("prefabs"); ok {
		for _, val := range prefabs.GetChilds() {
			var it = item{}
			if it.init(ig, val) {
				ig.Prefabs[it.Id] = &it
			}
		}
	}

	if items, ok := ig.itemsVDF.Get("items"); ok {
		for _, val := range items.GetChilds() {
			var it = item{}
			if it.init(ig, val) {
				ig.Items[it.Id] = &it
			}
		}
	}

	if colors, ok := ig.itemsVDF.Get("colors"); ok {
		for _, val := range colors.GetChilds() {
			var c = color{}
			if c.init(ig, val) {
				ig.Colors[c.Id] = &c
			}
		}
	}

	if particles, ok := ig.itemsVDF.Get("attribute_controlled_attached_particles"); ok {
		for _, val := range particles.GetChilds() {
			var p = particle{}
			if p.init(ig, val) {
				ig.Particles[p.Id] = &p
			}
		}
	}

	for _, it := range ig.Items {
		it.initPrefabs()
	}
}

func (ig *itemsGame) getPrefab(prefabName string) *item {
	return ig.Prefabs[prefabName]
}

func (ig *itemsGame) getItemsPerHero() map[string]*hero {
	heroes := make(map[string]*hero)

	for _, item := range ig.Items {
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

func (ig *itemsGame) getColors() []*color {
	ret := []*color{}

	for _, color := range ig.Colors {
		ret = append(ret, color)
	}

	return ret
}

func (ig *itemsGame) getParticles() []*particle {
	ret := []*particle{}

	for _, particle := range ig.Particles {
		ret = append(ret, particle)
	}

	return ret
}
