package markov

import (
	"fmt"
	"math/rand"
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
			word := []rune(words[i])
			lc := word[len(word)-1]
			if lc == '.' || lc == ',' || lc == '?' ||
				lc == '!' || lc == ';' || lc == ':' ||
				(unicode.IsLetter(lc) || unicode.IsDigit(lc)) &&
					(unicode.IsLetter(fc) || unicode.IsDigit(fc)) {
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
	}
	const max = 10
	const min = 3
	mr := rand.New(rand.NewSource(time.Now().Unix()))
	x := mr.Int31n((max-min)+1) + min
	fmt.Println(x, "words for you:", WordJoin(tg.Generate(uint(x))))
}
