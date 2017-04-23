package taxonomy

import (
	"fmt"
	"strings"
	"sync"

	"sort"

	"github.com/xreception/google-taxonomy-go/data"
)

const KeyLanguage = "en-US"

func NewTaxonomy(langs []string) (*Taxonomy, error) {
	tx := Taxonomy{
		rootIndex: make(categoryIDIndex),
		idIndex:   make(categoryIDIndex),
		langDict:  make(languageDictionary),
		nameIndex: make(categoryNameIndex),
	}
	td := taxonomyData{
		Language: KeyLanguage,
		LoadFunc: data.Asset,
	}
	td.Parse()
	err := tx.init(&td)
	if err != nil {
		return nil, err
	}
	tds := []*taxonomyData{}
	for _, lang := range langs {
		if lang == KeyLanguage {
			continue
		}
		td := taxonomyData{
			Language: lang,
			LoadFunc: data.Asset,
		}
		if err := td.Parse(); err != nil {
			return nil, err
		}
		tds = append(tds, &td)
	}
	for _, td := range tds {
		tx.loadLanguage(td)
	}
	return &tx, nil
}

type Taxonomy struct {
	rootIndex categoryIDIndex
	idIndex   categoryIDIndex
	langDict  languageDictionary
	mux       sync.RWMutex

	// nameIndex only used when init
	nameIndex categoryNameIndex
}

func (t *Taxonomy) init(data *taxonomyData) error {
	t.mux.Lock()
	defer func() {

		// release memory of nameIndex
		t.nameIndex = nil
		t.mux.Unlock()
	}()
	ld := make(map[int64]string)
	for _, r := range data.data {
		c := category{
			ID:       r.ID,
			Children: make(map[int64]*category),
		}
		if _, ok := t.idIndex[r.ID]; !ok {
			t.idIndex[r.ID] = &c
			t.rootIndex[r.ID] = &c
			t.nameIndex[r.Name] = &c
		}
		ld[r.ID] = r.Name
	}
	for _, r := range data.data {
		for i, cn := range r.Parents {
			var chn string
			if i < len(r.Parents)-1 {
				chn = r.Parents[i+1]
			} else {
				chn = r.Name
			}
			c, ch := t.nameIndex[cn], t.nameIndex[chn]
			if c == nil {
				return fmt.Errorf("%s is not found in nameIndex", cn)
			}
			if ch == nil {
				return fmt.Errorf("%s is not found in nameIndex", chn)
			}
			c.AppendChild(ch)
			ch.SetParent(c)
			delete(t.rootIndex, ch.ID)
		}
	}
	t.langDict[data.Language] = ld
	return nil
}
func (t *Taxonomy) hasLanguage(lang string) bool {
	t.mux.Lock()
	defer t.mux.Unlock()
	_, has := t.langDict[lang]
	return has
}
func (t *Taxonomy) loadLanguage(data *taxonomyData) {
	t.mux.Lock()
	defer t.mux.Unlock()
	ld := make(map[int64]string)
	if _, ok := t.langDict[data.Language]; ok {
		return
	}
	for _, r := range data.data {
		ld[r.ID] = r.Name
	}
	t.langDict[data.Language] = ld
}
func (t *Taxonomy) GetRootsCategoryInfo(lang string) ([]*CategoryInfo, error) {
	infs := []*CategoryInfo{}
	for _, c := range t.rootIndex {
		inf := c.GetInfo()
		if err := t.translate(inf, lang); err != nil {
			return []*CategoryInfo{}, nil
		}
		infs = append(infs, inf)
	}
	return infs, nil
}

func (t *Taxonomy) GetCategoryInfo(id int64, lang string) (*CategoryInfo, error) {
	t.mux.RLock()
	defer t.mux.RLock()
	c, _ := t.idIndex[id]
	if c == nil {
		return nil, nil
	}
	inf := c.GetInfo()
	if err := t.translate(inf, lang); err != nil {
		return nil, err
	}
	return inf, nil
}
func (t *Taxonomy) translateCategoryData(dat *categoryData, lang string) error {
	ld := t.langDict[lang]
	if ld == nil {
		return fmt.Errorf("language %s is not found", lang)
	}
	if name, ok := ld[dat.ID]; ok {
		dat.Name = name
	} else if name, ok := t.langDict[KeyLanguage][dat.ID]; ok {
		dat.Name = name
	} else {
		return fmt.Errorf("%d not found", dat.ID)
	}
	return nil
}
func (t *Taxonomy) translate(inf *CategoryInfo, lang string) error {
	cd := categoryData{
		ID: inf.ID,
	}
	if err := t.translateCategoryData(&cd, lang); err != nil {
		return err
	}
	inf.Name = cd.Name
	for _, dat := range inf.Fullpath {
		err := t.translateCategoryData(dat, lang)
		if err != nil {
			return err
		}
	}
	for _, dat := range inf.Children {
		err := t.translateCategoryData(dat, lang)
		if err != nil {
			return err
		}
	}
	inf.Language = lang
	return nil
}

type categoryData struct {
	ID   int64
	Name string
}

func (dat *categoryData) String() string {
	if len(dat.Name) > 0 {
		return dat.Name
	}
	return fmt.Sprintf("ID(%d)", dat.ID)
}

// CategoryInfo describes a category path and it's children
type CategoryInfo struct {
	Language string
	ID       int64
	Name     string
	Fullpath []*categoryData
	Children []*categoryData
}

func (inf *CategoryInfo) String() string {
	var (
		prs = []string{}
		chs = []string{}
	)
	for _, pr := range inf.Fullpath {
		prs = append(prs, pr.String())
	}

	fp := strings.Join(prs, " > ")
	if len(inf.Children) == 0 {
		return fp
	}
	for _, ch := range inf.Children {
		chs = append(chs, ch.String())
	}
	sort.Strings(chs)
	return strings.Join([]string{
		fp,
		fmt.Sprintf("{%s}", strings.Join(chs, " | ")),
	}, " > ")
}

type category struct {
	Parent   *category
	ID       int64
	Children map[int64]*category
	mux      sync.Mutex
}

func (c *category) SetParent(parent *category) error {
	c.mux.Lock()
	defer c.mux.Unlock()
	if c.Parent != nil && c.Parent != parent {
		return fmt.Errorf("%d already have a parent, insert duplicated", c.ID)
	}
	c.Parent = parent
	return nil
}
func (c *category) AppendChild(ch *category) {
	c.mux.Lock()
	defer c.mux.Unlock()
	if _, ok := c.Children[ch.ID]; !ok {
		c.Children[ch.ID] = ch
	}
}

func (c *category) GetInfo() *CategoryInfo {
	inf := CategoryInfo{
		ID:       c.ID,
		Fullpath: []*categoryData{},
		Children: []*categoryData{},
	}
	for id := range c.Children {
		inf.Children = append(inf.Children, &categoryData{
			ID: id,
		})
	}
	prs := []int64{c.ID}
	it := c.Parent
	for it != nil {
		prs = append(prs, it.ID)
		it = it.Parent
	}
	for i := range prs {
		inf.Fullpath = append(inf.Fullpath, &categoryData{
			ID: prs[len(prs)-i-1],
		})
	}
	return &inf
}
