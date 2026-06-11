package main

import (
	"os"
	"strings"

	"github.com/baldurstod/vdf"
)

type language struct {
	lang   string
	tokens map[string]string
}

func (l *language) init(path string) error {
	dat, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	v := vdf.VDF{}
	languageVdf := v.Parse(dat)

	lang, err := languageVdf.Get("lang")
	if err != nil {
		panic("lang key not found")
	}
	language, err := lang.GetString("Language")
	if err != nil {
		panic("Language key not found")
	}

	tokens, err := lang.Get("Tokens")
	if err != nil {
		panic("Tokens key not found")
	}

	l.lang = language
	l.tokens = make(map[string]string)
	for _, val := range tokens.GetChilds() {
		if s, err := val.ToString(); err == nil {
			l.tokens[strings.ToLower(val.Key)] = s
		}

	}
	return nil
}

func (l *language) getToken(token string) (string, bool) {
	token = strings.TrimPrefix(token, "#")
	s, ok := l.tokens[strings.ToLower(token)]
	return s, ok
}
