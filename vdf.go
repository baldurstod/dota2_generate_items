package main

import (
	"github.com/golang-collections/collections/stack"
	"unicode/utf8"
	_ "log"
	"fmt"
	_ "reflect"
	"encoding/json"
)

var depth = 0

type Token int

const (
	openingBrace Token = iota
	closingBrace
	newLine
	stringValue
	endToken
)

type VDF struct {
	s []byte
	i int
	len int
}

type KeyValue struct {
	key string
	value interface{}
	isRoot bool `default:false`
}

func PrintTabs(tabs int) {
	for i := 0; i < tabs; i++ {
		fmt.Print("\t")
	}
}

func (this KeyValue) GetString(key string) (string, bool) {
	a, ok := this.Get(key)
	if ok {
		switch a.value.(type) {
		case string:
			return a.value.(string), true
		}
	}
	return "", false
}

func (this KeyValue) Get(key string) (*KeyValue, bool) {
	//fmt.Println(reflect.TypeOf(this.value))
	switch this.value.(type) {
	case string:
		return nil, false
	case []*KeyValue:
		arr := this.value.([]*KeyValue)
		for _, item := range arr {
			if key == item.key {
				return item, true
			}
		}
	}
	return nil, false
}

func (this KeyValue) GetSubElement(path []string/*, params ...string*/) (*KeyValue, bool) {
	/*separator := "."
	if len(params) == 1 {
		separator = params[0]
	}
	paths := strings.Split(path, separator)*/

	fmt.Println(path[0]/*, path[1:]*/)
	if subElement, ok := this.Get(path[0]); ok {
		if len(path) == 1 {
			return subElement, true
		} else {
			return subElement.GetSubElement(path[1:])
		}
	}
	return nil, false
}

func (this KeyValue) GetChilds() ([]*KeyValue) {
	switch this.value.(type) {
	case []*KeyValue:
		return this.value.([]*KeyValue)
	}
	return []*KeyValue{}
}

func (this KeyValue) ToStringMap() (*map[string]string, bool) {
	switch this.value.(type) {
	case []*KeyValue:
		ret := make(map[string]string)
		arr := this.value.([]*KeyValue)
		for _, item := range arr {
			switch item.value.(type) {
			case string:
				//fmt.Println(item)
				//panic("remove me")
				ret[item.key] = item.value.(string)
			}
		}
		return &ret, true
	}
	return nil, false
}

func (this KeyValue) GetStringMap(key string) (*map[string]string, bool) {
	if sub, ok := this.Get(key); ok {
		return sub.ToStringMap()
	}
	return nil, false
}

func (this KeyValue) RemoveDuplicates() {
	switch this.value.(type) {
	case []*KeyValue:
		allKeys := make(map[string]bool)
		list := []*KeyValue{}

		arr := this.value.([]*KeyValue)
		for _, item := range arr {
			key := item.key
			if _, value := allKeys[key]; !value {
				allKeys[key] = true
				list = append(list, item)
				item.RemoveDuplicates()
			}
		}
		this.value = list
	}
}

func (this KeyValue) Print(optional ...int) {
	tabs := 0
	if len(optional) > 0 {
		tabs = optional[0]
	}

	if this.isRoot {
		tabs = -1
	}

	switch this.value.(type) {
	case []*KeyValue:
		if !this.isRoot {
			PrintTabs(tabs)
			fmt.Println("\"" + this.key + "\"")
			PrintTabs(tabs)
			fmt.Println("{")
		}
		arr := this.value.([]*KeyValue)
		for _, val := range arr {
			val.Print(tabs + 1);
		}

		if !this.isRoot {
			PrintTabs(tabs)
			fmt.Println("}")
		}
	case string:
		PrintTabs(tabs)
		fmt.Println("\"" + this.key + "\"		\"" + this.value.(string) + "\"")
	default:
		fmt.Println(this)
		panic("unknown type")
	}
}

func (this *KeyValue) toJSON() interface{} {
	ret := make(map[string]interface{})

	switch this.value.(type) {
	case string:
		//ret[this.key] = this.value.(string)
		return this.value.(string)
	case []*KeyValue:
		arr := this.value.([]*KeyValue)
		for _, item := range arr {
			ret[item.key] = item.toJSON()
		}
	}

	return ret
}

func (this KeyValue) MarshalJSON() ([]byte, error) {
	//ret := make(map[string]interface{})
	var ret interface{}

	//fmt.Println(this)

/*	switch this.value.(type) {
	case string:
		ret[this.key] = this.value.(string)
	case []*KeyValue:
		arr := this.value.([]*KeyValue)
		for _, item := range arr {
			ret[item.key] = item.toJSON()
			fmt.Println(item.key, item.value)
		}
	}
	*/

	ret = this.toJSON()

	return json.Marshal(ret)
}


func (this *VDF) Parse(s []byte) KeyValue {
	this.s = s
	this.i = 0
	this.len = len(s)
	//result := KeyValue{value: []*KeyValue{}}

	stringStack := stack.New()
	levelStack := stack.New()
	//currentLevel := &result
	var currentLevel *KeyValue = &KeyValue{key: "root", value: []*KeyValue{}, isRoot: true}
	var result KeyValue

TokenLoop:
	for {
		token, s := this.getNextToken()
		switch token {
		case openingBrace:
			key := stringStack.Pop().(string)
			subLevel := KeyValue{key: key, value: []*KeyValue{}}

			if currentLevel != nil {
				currentLevel.value = append(currentLevel.value.([]*KeyValue), &subLevel)
			}

			levelStack.Push(currentLevel)
			currentLevel = &subLevel
		case closingBrace:
			currentLevel = levelStack.Pop().(*KeyValue)
			if currentLevel != nil {
				result = *currentLevel
			}
		case newLine:
			if stringStack.Len() > 1 {
				value := stringStack.Pop().(string)
				key := stringStack.Pop().(string)
				currentLevel.value = append(currentLevel.value.([]*KeyValue), &KeyValue{key: key, value: value})
			}
		case stringValue: stringStack.Push(s)
		case endToken: break TokenLoop
		}

	}

	//result.Print()
	return result
}

func (this *VDF) getNextRune() (rune, int) {
	c, size := utf8.DecodeRune(this.s)
	this.s = this.s[size:]

	return c, size
}

func (this *VDF) getNextToken() (Token, string) {
	for this.i < this.len {
		c, size := this.getNextRune()
		this.i += size
		switch c {
		case '{': return openingBrace, ""
		case '}': return closingBrace, ""
		case '\r', '\n': return newLine, ""
		case ' ', '\t'://just eat a char
		case '"':
			s := ""
			for this.i < this.len {
				c, size := this.getNextRune()
				this.i += size
				switch c {
				case '\\':
					if this.i < this.len {
						c, size := this.getNextRune()
						this.i += size
						if c == '"' {
							s += "\\\""
						} else {
							s += `\` + string(c)
						}
					}
				case '"':
					return stringValue, s
				default:
					s += string(c)
				}
			}
		case '/':
			for this.i < this.len {
				c, size := this.getNextRune()
				this.i += size
				if c == '\r' || c == '\n' {
					break
				}
			}
		}
	}
	return endToken, ""
}