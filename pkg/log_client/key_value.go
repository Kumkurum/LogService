package log_client

import "fmt"

type KeyValue struct {
	Key   string
	Value string
}

func ConvertToMap(kvPairs ...KeyValue) map[string]string {
	result := make(map[string]string)
	for _, pair := range kvPairs {
		result[pair.Key] = pair.Value
	}
	return result
}
func ConvertToMessage(kvPairs ...KeyValue) string {
	message := ""
	for _, kv := range kvPairs {
		message += fmt.Sprintf("%s: %s; ", kv.Key, kv.Value)
	}
	return message
}
