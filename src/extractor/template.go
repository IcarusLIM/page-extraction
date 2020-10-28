package extract

import (
	"encoding/json"
)

type Template struct {
	ListName  string     `json:"listName,omitempty"`
	Method    string     `json:"method"`
	Selector  string     `json:"selector"`
	ListData  []Template `json:"listData,omitempty"`
	Attribute string     `json:"attribute,omitempty"`
	Regexp    string     `json:"regexp,omitempty"`
}

func NewTemplate(templateStr string) (*Template, error) {
	var template Template
	if err := json.Unmarshal([]byte(templateStr), &template); err != nil {
		return nil, err
	}
	return &template, nil
}
