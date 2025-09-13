package stivenking

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

type ItemJson struct {
	Body   string   `json:"body"`
	Source []string `json:"source"`
	Tags   []string `json:"tags"`
}

func LoadJsonItems(filePath string) ([]ItemJson, error) {
	// read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	var obj []ItemJson

	// unmarshall it
	err = json.Unmarshal(data, &obj)
	if err != nil {
		fmt.Println("error:", err)
		return nil, err
	}

	return obj, nil
}

func SetReplaces(text string) string {
	for _, name := range []string{"Показать цитату целиком"} {
		text = strings.Replace(text, name, "", -1)
	}

	return text
}

func GetQuote() string {
	quote, err := LoadJsonItems("./stivenking/stiven-king_09-08-2025_12.json")
	if err != nil {
		fmt.Print(err)
		return "not found ;("
	}

	return GetQuoteRandom(quote)
}

// return random string
func GetQuoteRandom(quote []ItemJson) string {
	rand.Seed(time.Now().UnixNano())
	index := rand.Intn(len(quote))

	msg := []string{
		SetReplaces(quote[index].Body),
		strings.Repeat("~", 39),
	}

	if len(quote[index].Source) > 0 {
		msg = append(msg, "👾 → "+strings.Join(quote[index].Source, ", "))
	}

	if len(quote[index].Tags) > 0 {
		msg = append(msg, "👽 → "+strings.Join(quote[index].Tags, ", "))
	}

	return strings.Join(msg, "\n")
}
