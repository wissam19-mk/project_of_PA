package utils

import (
	"encoding/json"
	"math"
	"regexp"
	"slices"
)

var commentRegex = regexp.MustCompile("//.*")

func NewUserConfig(source string) (*UserConfig, error) {
	var m UserConfig
	newSource := commentRegex.ReplaceAllString(source, "")
	err := json.Unmarshal([]byte(newSource), &m)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

func newModuleConfig(source string) (*ModuleConfig, error) {
	var m ModuleConfig
	newSource := commentRegex.ReplaceAllString(source, "")
	err := json.Unmarshal([]byte(newSource), &m)

	if err != nil {
		return nil, err
	}

	return &m, nil
}

type factorization struct {
	A int
	B int
}

func ComputeBestArea(size int) (int, int) {

	/*
		for i := size - 1; i > 1; i-- {
			if size%i == 0 {
				return i, size / i
			}
		}

		// size is prime
		root := int(math.Sqrt(float64(size)))
		return root, root
	*/
	
	var facts []factorization
	// Calculate all the possible factorizations
	for i := size - 1; i > 1; i-- {
		if size%i == 0 {
			facts = append(facts, factorization{A: size / i, B: i})
		}
	}

	// size is prime
	if len(facts) == 0 {
		root := int(math.Sqrt(float64(size)))
		return root, root
	}

	// Sort the factorizations
	slices.SortStableFunc(facts, func(a, b factorization) int {
		distA := int(math.Abs(float64(a.A - a.B)))
		distB := int(math.Abs(float64(b.A - b.B)))
		if distA != distB {
			return distA - distB
		}

		bDiff := a.B - b.B
		if bDiff != 0 {
			return bDiff
		}

		return a.A - b.A
	})

	/*
		for _, fact := range facts {
			Log(fmt.Sprintf("(%d, %d)", fact.A, fact.B))
		}
	*/

	return facts[0].A, facts[0].B

}
