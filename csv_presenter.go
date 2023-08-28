package collector

import (
	"bytes"
	"encoding/csv"
	"log"
)

type CSVPresenter struct {
	Delimiter rune
}

func (p *CSVPresenter) Render(deprecatedFunctions []*DeprecatedFunction) (string, error) {
	records := make([][]string, len(deprecatedFunctions))

	for i, d := range deprecatedFunctions {
		records[i] = []string{*d.Conf.FunctionArn, string(d.Conf.Runtime), d.DeprecationDate}
	}

	buff := &bytes.Buffer{}
	w := csv.NewWriter(buff)
	w.Comma = p.Delimiter

	for _, record := range records {
		if err := w.Write(record); err != nil {
			log.Fatalln("error writing record to csv:", err)
		}
	}

	w.Flush()

	return buff.String(), nil
}
