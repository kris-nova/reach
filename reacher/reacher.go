package reacher

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/elisarver/reach/lists"
	"github.com/elisarver/reach/tag"
	"github.com/elisarver/reach/target"
)

// Selector provides a statement goquery can use in a Find call.
type Selector interface {
	Select() string
}

// Mapper generates an approprirate goquery map
// function to retrieve a tag's attribute.
type Mapper interface {
	Map() func(int, *goquery.Selection) string
}

// SelectorMapper is the intersection of something that can select results and map them over functions.
type SelectorMapper interface {
	Selector
	Mapper
}

// SelectMap selects elements and maps them to response. Drops empty values.
func SelectMap(doc *goquery.Document, fm SelectorMapper) []string {
	return lists.DropEmpties(doc.Find(fm.Select()).Map(fm.Map()))
}

// TagSelectorMapper applies SelectorMapper to Tag
type TagSelectorMapper struct {
	*tag.Tag
}

// Select returns a tag's CSS select string.
func (tr TagSelectorMapper) Select() string {
	return tr.CSSSelector
}

// Map provides the selection function for a goquery.Map.
func (tr TagSelectorMapper) Map() func(int, *goquery.Selection) string {
	return func(_ int, sel *goquery.Selection) string {
		s, _ := sel.Attr(tr.Attribute)
		return s
	}
}

// documentFn exists to make testing possible without resorting to hardcoded function.
type documentFn func(string) (*goquery.Document, error)

// TargetReacher is any function that fetches tags on targets.
type TargetReacher func([]target.Target, []*tag.Tag) ([]string, error)

// genReachTargets binds the appropriate function to generate a document
// to the ReachTargets func.
func genReachTargets(fn documentFn) TargetReacher {
	if fn == nil {
		fn = goquery.NewDocument
	}
	return func(ts []target.Target, tags []*tag.Tag) ([]string, error) {
		var output []string
		for _, t := range ts {
			resp, err := fn(t.String())
			if err != nil {
				return []string{}, err
			}
			for _, t := range tags {
				output = append(output, SelectMap(resp, TagSelectorMapper{Tag: t})...)
			}
		}
		return output, nil
	}
}

var (
	// ReachTargets takes a list of targets, a list of tags, and a fetcher, fetches the targets with the
	// fetcher, and finds the tags in the document.
	ReachTargets = genReachTargets(nil)
)
