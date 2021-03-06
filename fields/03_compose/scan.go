package main

import (
	"fmt"
	"strings"
)

type jsonConfig struct {
	Fields []jsonField
}

type jsonField struct {
	Name       string
	Type       string
	Val        interface{}
	Multiplier interface{}
}

//gistsnip:start:scan
func (config *jsonConfig) Scan(fields ...Field) error {
	for _, dst := range fields {
		name := dst.Name()
		src, err := config.findField(name)
		if err != nil {
			return err
		}

		err = dst.Assign(src.Type, src.Val)
		if err != nil {
			return err
		}
	}
	return nil
}

//gistsnip:end:scan

func (config *jsonConfig) findField(name string) (*jsonField, error) {
	for i := 0; i < len(config.Fields); i++ {
		field := &config.Fields[i]
		if strings.EqualFold(field.Name, name) {
			return field, nil
		}
	}
	return nil, fmt.Errorf("unable to find field " + name)
}
