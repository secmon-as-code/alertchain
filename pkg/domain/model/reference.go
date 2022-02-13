package model

import (
	"github.com/google/uuid"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Reference struct {
	ID      types.ReferenceID
	AlertID types.AlertID

	Source  string
	Title   string
	URI     string
	Comment string
}

type References []*Reference

func (x *Alert) NewReference(src, title, uri, comment string) *Reference {
	return &Reference{
		ID:      types.ReferenceID(uuid.NewString()),
		AlertID: x.ID,

		Source:  src,
		Title:   title,
		URI:     uri,
		Comment: comment,
	}
}
