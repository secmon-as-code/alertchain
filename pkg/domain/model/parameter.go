package model

import "github.com/m-mizutani/alertchain/pkg/domain/types"

type Parameter struct {
	Key   types.ParamKey
	Value types.ParamValue
	Type  types.ParameterType
}

type Parameters []Parameter

func (x Parameters) Overwrite(src Parameters) Parameters {
	resp := x[:]
	exists := map[types.ParamKey]int{}
	for i, p := range x {
		exists[p.Key] = i
	}

	for _, p := range src {
		if idx, ok := exists[p.Key]; ok {
			resp[idx] = p
		} else {
			resp = append(resp, p)
		}
	}

	return resp
}
