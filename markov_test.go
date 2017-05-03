package markov

import (
	"fmt"
	"os"
	"testing"
	"time"
	"unicode"
)

func WordJoin(words []string) string {
	text := ""
	for i, _ := range words {
		text += words[i]
		isLast := i == (len(words) - 1)
		if !isLast {
			next := words[i+1]
			fc := []rune(next)[0]
			if unicode.IsLetter(fc) || unicode.IsDigit(fc) {
				text += " "
			}
		}
	}
	return text
}

func TestGenerator(t *testing.T) {
	var tg TextGenerator
	tg.Init(time.Now().Unix())
	{
		fd, err := os.Open("alice.txt")
		if err != nil {
			panic(err)
		}
		tg.Feed(fd)
		fd.Close()
	}
	fmt.Println(WordJoin(tg.Generate(10)))
}
