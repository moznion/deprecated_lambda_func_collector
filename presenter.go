package collector

type Presenter interface {
	Render(deprecatedFunctions []*DeprecatedFunction) (string, error)
}

type Output struct {
	ARN               string `json:"arn"`
	DeprecatedRuntime string `json:"deprecated_runtime"`
	DeprecationDate   string `json:"deprecation_date"`
}

func TransformToOutput(deprecatedFunctions []*DeprecatedFunction) []Output {
	outputs := make([]Output, len(deprecatedFunctions))

	for i, c := range deprecatedFunctions {
		outputs[i] = Output{
			ARN:               *c.Conf.FunctionArn,
			DeprecatedRuntime: string(c.Conf.Runtime),
			DeprecationDate:   c.DeprecationDate,
		}
	}

	return outputs
}
