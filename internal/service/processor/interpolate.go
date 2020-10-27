package processor

import (
	"bytes"
	"encoding/json"
	"text/template"

	"gitlab.com/distributed_lab/logan/v3/errors"
)

func interpolate(tmpl string, payload []byte) ([]byte, error) {
	t := template.New("tmpl")
	t, err := t.Parse(tmpl)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse template")
	}

	p := make(map[string]interface{})
	if err = json.Unmarshal(payload, &p); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal payload")
	}

	var res bytes.Buffer
	if err = t.Execute(&res, p); err != nil {
		return nil, errors.Wrap(err, "failed to execute template")
	}

	return res.Bytes(), nil
}
