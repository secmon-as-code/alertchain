package github

import (
	"io"

	"github.com/secmon-lab/alertchain/pkg/domain/model"
)

func ExecuteTemplate(w io.Writer, alert model.Alert) error {
	return issueTemplate.Execute(w, alert)
}
