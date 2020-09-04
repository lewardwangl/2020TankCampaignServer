package utils

import "encoding/json"

func Struct2str(obj interface{}) string {
	str, _ := json.Marshal(obj)
	return string(str)
}
