package interfaces

import (
	"io"

	"github.com/m-mizutani/alertchain/pkg/domain/model"
	"github.com/m-mizutani/alertchain/pkg/domain/types"
)

type Schema interface {
	ID() types.SchemaID
	Parse(r io.Reader) ([]model.Alert, error)
}
