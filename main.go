package main

import (
	"encoding/json"
	"os"
	_ "log"
	"flag"
	"fmt"
	"path"
)
var lg language

func main() {
	var lang string
	var outputFolder string
	var itemsFolder string
	var resourceFolder string
	var staticFile string
	var medals bool

	flag.StringVar(&lang, "l", "english", "Language")
	flag.BoolVar(&medals, "m", false, "Tournament medals")
	flag.StringVar(&outputFolder, "o", "", "Output folder")
	flag.StringVar(&itemsFolder, "i", "", "Items folder")
	flag.StringVar(&resourceFolder, "r", "", "Resource folder")
	flag.StringVar(&staticFile, "s", "", "Static file")
	flag.Parse()

	if itemsFolder == "" {
		fmt.Println("No items folder provided. Use the flag -i")
		os.Exit(1)
	}
	if resourceFolder == "" {
		fmt.Println("No resource folder provided. Use the flag -r")
		os.Exit(1)
	}
	if outputFolder == "" {
		fmt.Println("No output folder provided. Use the flag -o")
		os.Exit(1)
	}

	//file, _ := os.OpenFile("var\\dota2_generate_items.log", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	//os.Stdout = file

	lg = language{}
	lg.init(path.Join(resourceFolder, "items_" + lang + ".txt"))
	//fmt.Println(lg/*.getToken("#DOTA_Style_Spirit_Breaker_TI10_Cavern_Crawl_Tail_1")*/);

	ig := itemsGame{}
	itemsGameDatas, _ := os.ReadFile(path.Join(itemsFolder, "items_game.txt"))
	ig.init(itemsGameDatas)

	/*sub, _ := ig.itemsVDF.GetSubElement([]string{"items", "12534", "visuals", "styles"})
	j2, _ := json.MarshalIndent(sub, "", "\t")
	os.WriteFile(path.Join(outputFolder, "test", "test.json"), j2, 0666)
	return*/




	heroes := ig.getItemsPerHero()


	for npc, hero := range heroes {
		j, _ := json.MarshalIndent(hero, "", "\t")
		os.WriteFile(path.Join(outputFolder, "test", npc + ".json"), j, 0666)
	}

	return

	//heroes := []*hero{}
	//fmt.Println(ig.Prefabs);
	//return

	//dat, _ := os.ReadFile(path.Join(itemsFolder, "items_game.txt"))
	//dat, _ := os.ReadFile(path.Join(itemsFolder, "items_game_test.txt"))
	//dat, _ := os.ReadFile(path.Join(itemsFolder, "items_game_2025.txt"))
	//dat, _ := os.ReadFile(path.Join(itemsFolder, "tf2_items_game.txt"))
	//vdf := VDF{}
	//kv := vdf.Parse(dat)
	//kv.Print()
	//return

	//fmt.Println(kv.GetChild("items"))
	//kv.Get("items").Print()
	//kv.key = "a"

	//j, _ := json.MarshalIndent(&ig, "", "\t")*/
	j, _ := json.MarshalIndent(&ig, "", "\t")
	os.WriteFile(path.Join(outputFolder, "items_" + lang + ".json"), j, 0666)

	/*h := hero{items: []*item{&item{}}}
	j, _ = json.MarshalIndent(&h, "", "\t")
	os.WriteFile(path.Join(outputFolder, "test.json"), j, 0666)*/
}

/*
func getMap(i interface{}) itemGameMap {
	switch i.(type) {
	case itemGameMap: return i.(itemGameMap)
	case map[string]interface{}: return itemGameMap((i).(map[string]interface{}))
	default:
		fmt.Println(i);
		panic("Unknown type")
	}
}*/

func getStringToken(token string) string {
	s, exist := lg.getToken(token)

	if (exist) {
		return s
	} else {
		return token
	}
}

func getStringTokenRaw(token string) (string, bool) {
	return lg.getToken(token)
}
