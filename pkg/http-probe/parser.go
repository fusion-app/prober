package httpprobe

import (
	"fmt"
	"github.com/savaki/jq"
	"strconv"
)

type ValueTypeName string

const (
	String ValueTypeName = "string"
	Float  ValueTypeName = "float"
	Int    ValueTypeName = "int"
	Bool   ValueTypeName = "Bool"
)

func PKUAPIParse(src []byte) ([]byte, error) {
	selectors := []string{".data", ".result", ".response", "returnJSONStr"}

	var dataBody = src
	var parseBody interface{}
	var err error
	for _, selector := range selectors {
		parseBody, err = JQParse(dataBody, selector, String)
		if err != nil {
			return nil, fmt.Errorf("JQ '%s' Parse error: %+v", selector, err)
		}
		//log.Printf("After jq '%s', result :%s", selector, parseBody.(string))
		dataBody = []byte(parseBody.(string))
	}
	if string(dataBody) == "failed" {
		return nil, fmt.Errorf("PKU API result is failed")
	}
	return dataBody, nil
}

func JQParse(jsonData []byte, selector string, typeName ValueTypeName) (interface{}, error) {
	op, err := jq.Parse(selector)
	if err != nil {
		return nil, fmt.Errorf("Invalid jq error: %+v ", err.Error())
	}
	valueBytes, err := op.Apply(jsonData)
	if err != nil {
		return nil, fmt.Errorf("Not found value in json error: %+v ", err.Error())
	}
	switch typeName {
	case String:
		if val, err := strconv.Unquote(string(valueBytes)); err == nil {
			return val, nil
		} else {
			return nil, fmt.Errorf("Parse string value error: %+v ", err.Error())
		}
	case Float:
		if val, err := strconv.ParseFloat(string(valueBytes), 64); err == nil {
			return val, nil
		} else {
			return nil, fmt.Errorf("Parse float value error: %+v ", err.Error())
		}
	case Bool:
		if val, err := strconv.ParseBool(string(valueBytes)); err == nil {
			return val, nil
		} else {
			return nil, fmt.Errorf("Parse boolean value error: %+v ", err.Error())
		}
	case Int:
		if val, err := strconv.ParseInt(string(valueBytes), 10, 64); err == nil {
			return val, nil
		} else {
			return nil, fmt.Errorf("Parse boolean value error: %+v ", err.Error())
		}
	default:
		return nil, nil
	}
}
