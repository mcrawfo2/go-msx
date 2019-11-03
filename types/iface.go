package types

func InterfaceSliceToStringSlice(source []interface{}) (target []string) {
	for _, v := range source {
		target = append(target, v.(string))
	}
	return
}
