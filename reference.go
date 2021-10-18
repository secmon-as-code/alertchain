package alertchain

import "github.com/m-mizutani/alertchain/pkg/infra/ent"

type Reference struct {
	Source  string
	Title   string
	URL     string
	Comment string
}

type References []*Reference

func newReferences(bases []*ent.Reference) References {
	resp := make(References, len(bases))
	for i, ref := range bases {
		resp[i] = newReference(ref)
	}
	return resp
}

func newReference(base *ent.Reference) *Reference {
	return &Reference{
		Source:  base.Source,
		Title:   base.Title,
		URL:     base.URL,
		Comment: base.Comment,
	}
}

func (x *Reference) toEnt() *ent.Reference {
	return &ent.Reference{
		Source:  x.Source,
		Title:   x.Title,
		URL:     x.URL,
		Comment: x.Comment,
	}
}

func (x References) toEnt() []*ent.Reference {
	resp := make([]*ent.Reference, len(x))
	for i, ref := range x {
		resp[i] = ref.toEnt()
	}
	return resp
}
