package log_client

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
