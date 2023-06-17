package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type Parameter struct {
	ID    types.ParamID       `json:"id"`
	Name  types.ParamName     `json:"name"`
	Value types.ParamValue    `json:"value"`
	Type  types.ParameterType `json:"type"`
}

type Parameters []Parameter

func (x Parameters) Copy() Parameters {
	newParams := make(Parameters, len(x))
	for i, p := range x {
		newParams[i] = Parameter{
			ID:    p.ID,
			Name:  p.Name,
			Value: p.Value,
			Type:  p.Type,
		}
	}
	return newParams
}

func TidyParameters(params Parameters) Parameters {
	var ret Parameters

	idMap := map[types.ParamID]int{}

	for _, p := range params {
		if p.ID == "" {
			p.ID = types.NewParamID()
		}

		if _, ok := idMap[p.ID]; ok {
			ret[idMap[p.ID]] = p
		} else {
			ret = append(ret, p)
			idMap[p.ID] = len(ret) - 1
		}
	}

	return ret
}
