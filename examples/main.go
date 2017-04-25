package main

import (
	"fmt"

	"github.com/go-google-taxonomy/google-taxonomy/taxonomy"
)

const keyLanguage = "en-US"

func main() {
	tx, err := taxonomy.NewTaxonomy(keyLanguage, []string{})
	if err != nil {
		panic(err)
	}
	infs, err := tx.GetRootsCategoryInfo(keyLanguage)
	if err != nil {
		panic(err)
	}
	for _, inf := range infs {
		fmt.Printf("%d - %s\n", inf.ID, inf.String())
	}
	inf, err := tx.GetCategoryInfo(16, keyLanguage)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d - %s\n", inf.ID, inf.String())
}
