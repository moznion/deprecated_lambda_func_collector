package collector

import (
	"encoding/json"
	"fmt"
)

type JSONPresenter struct {
}

func (p *JSONPresenter) Render(deprecatedFunctions []*DeprecatedFunction) (string, error) {
	out, err := json.Marshal(TransformToOutput(deprecatedFunctions))
	if err != nil {
		return "", fmt.Errorf("failed to JSON marshalling for the output: %w", err)
	}
	return string(out), nil
}
