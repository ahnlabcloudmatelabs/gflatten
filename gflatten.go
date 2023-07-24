package gflatten

import "reflect"

var WRAP = struct {
	SQUARE_BRACKET string
	SINGLE_QUOTE   string
	DOUBLE_QUOTE   string
	BACKTICK       string
	NONE           string
}{
	SQUARE_BRACKET: "[]",
	SINGLE_QUOTE:   "'",
	DOUBLE_QUOTE:   "\"",
	BACKTICK:       "`",
	NONE:           "",
}

type Option struct {
	ParameterDelimiter string
	ArrayDelimiter     string
	Prefix             string
	ParameterWrap      string
	ArrayWrap          string
}

func Flatten(src any, option Option) (map[string]any, error) {
	dest := make(map[string]any)
	mapSrc, err := toMap(src)
	if err != nil {
		return dest, err
	}

	flatten(option.Prefix, mapSrc, dest, option)
	return dest, nil
}

func flatten(prefix string, src map[string]any, dest map[string]any, option Option) {
	for k, v := range src {
		switch reflect.TypeOf(v).Kind() {

		case reflect.Slice:
			flattenArray(prefix, k, v, dest, option)

		case reflect.Map:
			flatten(
				parameterWrap(prefix, k, option),
				v.(map[string]any),
				dest,
				option,
			)

		default:
			dest[parameterWrap(prefix, k, option)] = v
		}
	}
}

func flattenArray(prefix string, key string, value any, dest map[string]any, option Option) {
	for i := 0; i < len(value.([]any)); i++ {
		if _, ok := value.([]any)[i].(map[string]any); ok {
			flatten(
				arrayWrap(parameterWrap(prefix, key, option), i, option),
				value.([]any)[i].(map[string]any),
				dest,
				option,
			)
		} else {
			dest[arrayWrap(parameterWrap(prefix, key, option), i, option)] = value.([]any)[i]
		}
	}
}
