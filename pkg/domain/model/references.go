package model

type References []Reference

type Reference struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func (r Reference) Copy() Reference {
	newRef := Reference{
		Title: r.Title,
		URL:   r.URL,
	}
	return newRef
}

func (refs References) Copy() References {
	newRefs := make(References, len(refs))
	for i, ref := range refs {
		newRefs[i] = ref.Copy()
	}
	return newRefs
}
