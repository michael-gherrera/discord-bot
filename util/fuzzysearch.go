package util

import (
	"strings"

	iex "github.com/jonwho/go-iex"
)

// compute the edit-distance aka cost between the two strings
func Levenshtein(s, t string) int {
	sLen := len(s)
	tLen := len(t)
	cost := 0

	if sLen == 0 {
		return tLen
	}
	if tLen == 0 {
		return sLen
	}

	matrix := make([][]int, sLen+1)
	for i := range matrix {
		matrix[i] = make([]int, tLen+1)
	}

	for i := 0; i <= sLen; i++ {
		matrix[i][0] = i
	}
	for j := 0; j <= tLen; j++ {
		matrix[0][j] = j
	}

	for i := 1; i <= sLen; i++ {
		sRune := s[i-1]

		for j := 1; j <= tLen; j++ {

			tRune := t[j-1]

			if sRune == tRune {
				cost = 0
			} else {
				cost = 1
			}

			matrix[i][j] = Min3(matrix[i-1][j]+1, matrix[i][j-1]+1, matrix[i-1][j-1]+cost)
		}
	}

	return matrix[sLen][tLen]
}

// computes Soundex value for the string
func Soundex(s string) string {
	m := map[byte]string{
		'B': "1", 'P': "1", 'F': "1", 'V': "1",
		'C': "2", 'S': "2", 'K': "2", 'G': "2", 'J': "2", 'Q': "2", 'X': "2", 'Z': "2",
		'D': "3", 'T': "3",
		'L': "4",
		'M': "5", 'N': "5",
		'R': "6",
	}

	s = strings.ToUpper(s)

	r := string(s[0])
	p := s[0]

	for i := 1; i < len(s) && len(r) < 4; i++ {
		c := s[i]

		if (c < 'A' || c > 'Z') || (c == p) {
			continue
		}

		p = c

		if n, ok := m[c]; ok {
			r += n
		}
	}

	for i := len(r); i < 4; i++ {
		r += "0"
	}

	return r
}

func FuzzySearch(ticker string, symbols []iex.SymbolDTO) []iex.SymbolDTO {
	var fuzzySymbols []iex.SymbolDTO

	tickerSoundex := Soundex(ticker)
	for _, symbolDTO := range symbols {
		if tickerSoundex == Soundex(symbolDTO.Symbol) {
			fuzzySymbols = append(fuzzySymbols, symbolDTO)
		}
	}

	return fuzzySymbols
}

// TODO: tweak Levenshtein algo with other constraints to get better recommendations
// func FuzzySearch(ticker string, symbols []iex.SymbolDTO) []string {
//   minHeap := &MinStockHeap{}
//   heap.Init(minHeap)
//
//   levenshteinTimeStart := time.Now()
//   fmt.Println("Running fuzzy search... ", levenshteinTimeStart)
//   for _, symbolDTO := range symbols {
//     cost := Levenshtein(ticker, symbolDTO.Symbol)
//     heap.Push(minHeap, &struct {
//       LDCost int
//       Symbol string
//     }{cost, symbolDTO.Symbol})
//   }
//   fmt.Println("Finished running fuzzy search... ", time.Now().Sub(levenshteinTimeStart))
//
//   fuzzySymbols := make([]string, 10)
//   for i := 0; i < 10; i++ {
//     fuzzySymbols[i] = heap.Pop(minHeap).(*struct {
//       LDCost int
//       Symbol string
//     }).Symbol
//   }
//
//   return fuzzySymbols
// }
