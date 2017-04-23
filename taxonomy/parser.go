package taxonomy

import (
	"bytes"
	"fmt"
	"strings"
)

type languageDictionary map[string]map[int64]string
type categoryIDIndex map[int64]*category
type categoryNameIndex map[string]*category

type categoryRecord struct {
	ID      int64
	Parents []string
	Name    string
}

func parseLine(line string) (*categoryRecord, error) {
	if strings.HasPrefix(line, "#") {
		return nil, nil
	}
	hds := strings.Split(strings.TrimSpace(line), " - ")

	var id int64
	fmt.Sscanf(hds[0], "%d", &id)
	cps := strings.Split(hds[1], " > ")
	for i := range cps {
		cps[i] = strings.TrimSpace(cps[i])
	}
	name, cps := cps[len(cps)-1], cps[:len(cps)-1]
	return &categoryRecord{
		ID:      id,
		Parents: cps,
		Name:    name,
	}, nil
}

type taxonomyData struct {
	// Set language when create TaxonomyData
	Language string
	LoadFunc func(string) ([]byte, error)

	data []*categoryRecord
}

func (td *taxonomyData) Filename() string {
	return fmt.Sprintf("taxonomy-with-ids.%s.txt", td.Language)
}
func (td *taxonomyData) Parse() error {
	data, err := td.LoadFunc(td.Filename())
	if err != nil {
		return err
	}
	b := bytes.NewBuffer(data)
	for {
		l, err := b.ReadString('\n')
		if err != nil {
			break
		}
		cr, err := parseLine(l)
		if err != nil {
			return err
		}
		if cr != nil {
			td.data = append(td.data, cr)
		}
	}
	return nil
}
