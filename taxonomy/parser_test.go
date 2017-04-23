package taxonomy

import (
	"testing"

	"fmt"

	"github.com/xreception/google-taxonomy-go/data"
)

func TestTaxonomyDataParse(t *testing.T) {
	type testCase struct {
		lang string
		cnt  int
	}
	cases := []testCase{
		testCase{"cs-CZ", 5427},
		testCase{"da-DK", 5427},
		testCase{"de-CH", 5427},
		testCase{"de-DE", 5427},
		testCase{"en-US", 5427},
		testCase{"es-ES", 5427},
		testCase{"fr-FR", 5427},
		testCase{"it-IT", 5427},
		testCase{"ja-JP", 5442},
		testCase{"pl-PL", 5427},
		testCase{"pt-BR", 5427},
		testCase{"sv-SE", 5427},
		testCase{"zh-CN", 4586},
	}
	for _, cas := range cases {
		t.Run(fmt.Sprintf("Load taxonomy data file in %s", cas.lang), func(t *testing.T) {
			td := taxonomyData{
				Language: cas.lang,
				LoadFunc: data.Asset,
			}
			if td.Filename() != fmt.Sprintf("taxonomy-with-ids.%s.txt", cas.lang) {
				t.Fatalf("%s is not valid taxonomy data filename", td.Filename())
			}
			if err := td.Parse(); err != nil {
				t.Fatal(err)
			}
			if len(td.data) != cas.cnt {
				t.Fatalf("%s contains %d records, not eq. to %d", td.Filename(), len(td.data), cas.cnt)
			}
		})
	}
}
