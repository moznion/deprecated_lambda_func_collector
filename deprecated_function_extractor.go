package collector

import (
	_ "embed"

	"github.com/aws/aws-sdk-go-v2/service/lambda/types"
	"github.com/goccy/go-yaml"
)

//go:embed deprecated_runtime.yaml
var deprecatedRuntimeSourceYAMLString []byte

type DeprecatedRuntimeSource struct {
	Runtime         string `yaml:"runtime"`
	DeprecationDate string `yaml:"deprecation_date"`
	Planned         bool   `yaml:"planned"`
}

type DeprecatedRuntime struct {
	DeprecatedRuntimeToDate map[string]string
}

type DeprecatedFunctionExtractor struct {
	deprecatedRuntime         *DeprecatedRuntime
	deprecationPlannedRuntime *DeprecatedRuntime
}

func NewDeprecatedFunctionExtractor() (*DeprecatedFunctionExtractor, error) {
	var deprecatedRuntimeSource []DeprecatedRuntimeSource
	err := yaml.Unmarshal(deprecatedRuntimeSourceYAMLString, &deprecatedRuntimeSource)
	if err != nil {
		return nil, err
	}

	deprecatedRuntime := &DeprecatedRuntime{
		DeprecatedRuntimeToDate: map[string]string{},
	}
	deprecationPlannedRuntime := &DeprecatedRuntime{
		DeprecatedRuntimeToDate: map[string]string{},
	}

	for _, d := range deprecatedRuntimeSource {
		dst := deprecatedRuntime
		if d.Planned {
			dst = deprecationPlannedRuntime
		}
		dst.DeprecatedRuntimeToDate[d.Runtime] = d.DeprecationDate
	}

	return &DeprecatedFunctionExtractor{
		deprecatedRuntime:         deprecatedRuntime,
		deprecationPlannedRuntime: deprecationPlannedRuntime,
	}, nil
}

func (f *DeprecatedFunctionExtractor) Extract(configurations []types.FunctionConfiguration, includeFuturePlanned bool) []*DeprecatedFunction {
	var deprecatedFunctions []*DeprecatedFunction
	for _, conf := range configurations {
		if date := f.deprecatedRuntime.DeprecatedRuntimeToDate[string(conf.Runtime)]; date != "" {
			deprecatedFunctions = append(deprecatedFunctions, &DeprecatedFunction{
				Conf:            conf,
				DeprecationDate: date,
			})
			continue
		}

		if includeFuturePlanned {
			if date := f.deprecationPlannedRuntime.DeprecatedRuntimeToDate[string(conf.Runtime)]; date != "" {
				deprecatedFunctions = append(deprecatedFunctions, &DeprecatedFunction{
					Conf:            conf,
					DeprecationDate: date,
				})
				continue
			}
		}
	}

	return deprecatedFunctions
}
