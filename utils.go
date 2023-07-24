package gflatten

import (
	"encoding/json"
	"strconv"
)

func toMap(src any) (dest map[string]any, err error) {
	data, err := json.Marshal(src)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &dest)
	return
}

func parameterWrap(prefix, key string, option Option) string {
	_prefix := createPrefix(prefix, option.ParameterDelimiter)

	switch option.ParameterWrap {
	case WRAP.SQUARE_BRACKET:
		return _prefix + "[" + key + "]"
	case WRAP.SINGLE_QUOTE:
		return _prefix + "'" + key + "'"
	case WRAP.DOUBLE_QUOTE:
		return _prefix + "\"" + key + "\""
	case WRAP.BACKTICK:
		return _prefix + "`" + key + "`"
	default:
		return _prefix + key
	}
}

func arrayWrap(prefix string, index int, option Option) string {
	key := strconv.Itoa(index)
	_prefix := createPrefix(prefix, option.ArrayDelimiter)

	switch option.ArrayWrap {
	case WRAP.SQUARE_BRACKET:
		return _prefix + "[" + key + "]"
	case WRAP.SINGLE_QUOTE:
		return _prefix + "'" + key + "'"
	case WRAP.DOUBLE_QUOTE:
		return _prefix + "\"" + key + "\""
	case WRAP.BACKTICK:
		return _prefix + "`" + key + "`"
	default:
		return _prefix + key
	}
}

func createPrefix(prefix, delimiter string) string {
	if prefix == "" {
		return ""
	}

	return prefix + delimiter
}
