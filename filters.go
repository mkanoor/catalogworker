package main

import (
	"github.com/jmespath/go-jmespath"
	log "github.com/sirupsen/logrus"
	"strings"
)

// Filter stores the parsed value of the JMESPath expression
type Filter struct {
	Value          string
	ReplaceResults bool
}

// Apply the JMESPath filter to the JSON body recieved from
// Ansible Tower
func (f *Filter) Apply(jsonBody map[string]interface{}) (map[string]interface{}, error) {
	precompiled, err := jmespath.Compile(f.Value)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	result, err := precompiled.Search(jsonBody)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	if f.ReplaceResults {
		jsonBody["results"] = result
	} else {
		jsonBody = result.(map[string]interface{})
	}
	return jsonBody, nil
}

// Parse the filter value which can be a string or map.
// The map is typically used when working with a single object response
// The string filter value is used when working with a list response which
// can contain multiple objects and the filter needs to be applied to each
// object and the results collection be updated.
func (f *Filter) Parse(element interface{}) {
	switch element.(type) {
	case string:
		f.Value = element.(string)
		f.ReplaceResults = true
	case map[string]interface{}:
		var sb strings.Builder
		for key, value := range element.(map[string]interface{}) {
			switch value.(type) {
			case string:
				if sb.Len() == 0 {
					sb.WriteString("{")
				} else {
					sb.WriteString(",")
				}
				sb.WriteString(key + ":" + value.(string))
			}
		}
		sb.WriteString("}")
		f.Value = sb.String()
	}
}
