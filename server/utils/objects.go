package utils

import (
	"encoding/json"
	"io/ioutil"
)

func ToIndentedJSON(content interface{}) string {
	obj, err := json.MarshalIndent(content, "", "  ")
	if err != nil {
		panic(err)
	}
	return string(obj)
}

func UnmarshallJSONFromFile(fileName string, object interface{}) error {
	byts, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	return json.Unmarshal(byts, object)
}
