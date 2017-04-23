package taxonomy

import (
	"fmt"
	"testing"
)

const KeyLanguage = "en-US"

func TestTaxonomyGetRootsCategoryInfo(t *testing.T) {
	tx, err := NewTaxonomy(KeyLanguage, []string{})
	if err != nil {
		t.Fatal(err)
	}
	infs, err := tx.GetRootsCategoryInfo(KeyLanguage)
	if err != nil {
		t.Fatal(err)
	}
	if len(infs) != 21 {
		t.Fatalf("mismatch: Expeacted %d root categories, got %d", 21, len(infs))
	}
}
func TestTaxonomyLanguageCompatibility(t *testing.T) {
	langs := []string{
		"cs-CZ",
		"da-DK",
		"de-CH",
		"de-DE",
		"es-ES",
		"fr-FR",
		"it-IT",
		"ja-JP",
		"pl-PL",
		"pt-BR",
		"sv-SE",
	}
	tx, err := NewTaxonomy(KeyLanguage, langs)
	if err != nil {
		t.Fatal(err)
	}
	for id := range tx.idIndex {
		t.Run(fmt.Sprintf("Test language compatibility of %d", id), func(t *testing.T) {
			ref, err := tx.GetCategoryInfo(id, KeyLanguage)
			if err != nil {
				t.Fatal(err)
			}
			for _, lang := range langs {
				inf, err := tx.GetCategoryInfo(id, lang)
				if err != nil {
					t.Error(err)
					continue
				}
				if ref == nil {
					ref = inf
				} else if len(ref.Fullpath) != len(inf.Fullpath) || len(ref.Children) != len(inf.Children) {
					t.Fatalf("language_definition_mismatch: Category(%d:%s) mismatch with language %s", id, inf.String(), lang)
				}
			}
		})
	}
}
func TestNewTaxonomy(t *testing.T) {
	_, err := NewTaxonomy(KeyLanguage, []string{})
	if err != nil {
		t.Fatal(err)
	}
}
func TestTaxonomyGetCategoryInfo(t *testing.T) {
	type testCase struct {
		lang     string
		id       int64
		expected string
	}

	cases := []testCase{
		testCase{"en-US", 6536, "Mature > Erotic > Pole Dancing Kits"},
		testCase{"cs-CZ", 6536, "Pro dospělé > Erotika > Sady na tanec u tyče"},
		testCase{"da-DK", 6536, "Voksne > Erotik > Stripperstangsæt"},
		testCase{"de-CH", 6536, "Für Erwachsene > Erotik > Pole Dance-Tanzstangenkits"},
		testCase{"de-DE", 6536, "Für Erwachsene > Erotik > Pole Dance-Tanzstangenkits"},
		testCase{"es-ES", 6536, "Productos para adultos > Sexo > Kits de baile de barra americana"},
		testCase{"fr-FR", 6536, "Adulte > Érotisme > Kits de pole dance"},
		testCase{"it-IT", 6536, "Articoli per adulti > Erotismo > Kit per pole dance"},
		testCase{"ja-JP", 6536, "成人向け > アダルト > ポールダンスキット"},
		testCase{"pl-PL", 6536, "Artykuły dla dorosłych > Artykuły erotyczne > Zestawy do tańca na rurze"},
		testCase{"pt-BR", 6536, "Adultos > Erótico > Kits para pole dancing"},
		testCase{"sv-SE", 6536, "Avsett för vuxna > Erotik > Poledancekit"},
		testCase{"zh-CN", 1594, "服饰 > 服装 > 西服/套装 > {燕尾服/男士晚礼服 | 西服套装 | 西裙套装}"},
	}
	for _, cas := range cases {
		t.Run(fmt.Sprintf("New taxonomy in %s", cas.lang), func(t *testing.T) {
			tx, err := NewTaxonomy(KeyLanguage, []string{cas.lang})
			if err != nil {
				t.Fatal(err)
			}
			inf, err := tx.GetCategoryInfo(cas.id, cas.lang)
			if err != nil {
				t.Fatal(err)
			}
			if inf == nil {
				t.Fatal(fmt.Errorf("not_found: Category info of %d in %s", cas.id, cas.lang))
			}
			if inf.String() != cas.expected {
				t.Fatal(fmt.Sprintf("mismatch: Expected %s, got %s", cas.expected, inf.String()))
			}
		})
	}
}
