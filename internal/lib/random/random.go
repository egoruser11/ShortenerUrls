package random

import (
	"fmt"
	"math/rand"
)

func NewRandomString(aliasMaxLength int, aliases []string) string {
	generatedAlias := GenerateAlias(aliasMaxLength)
	if len(aliases) == 0 {
		return generatedAlias
	}
	for {
		if IsAliasNotInArray(generatedAlias, aliases) {
			return generatedAlias
		}
		generatedAlias = GenerateAlias(aliasMaxLength)
	}
}

func IsAliasNotInArray(generatedAlias string, aliases []string) bool {
	for _, a := range aliases {
		if a == generatedAlias {
			return false
		}
	}
	return true
}

func GenerateAlias(aliasMaxLength int) string {
	symbols := GetAllSymbols()
	result := "https://"
	for i := 0; i < aliasMaxLength; i++ {
		randomIndex := rand.Intn(len(symbols))
		result += symbols[randomIndex]
	}
	return result
}

func GetAllSymbols() []string {
	symbols := []string{}
	for ch := 'A'; ch <= 'Z'; ch++ {
		symbols = append(symbols, string(ch))
	}
	for ch := 'a'; ch <= 'z'; ch++ {
		symbols = append(symbols, string(ch))
	}
	for ch := '0'; ch <= '9'; ch++ {
		symbols = append(symbols, string(ch))
	}
	return symbols
}

func main() {
	arr := []string{"https://2312", "https://2312s"}
	fmt.Println(NewRandomString(6, arr))
}
