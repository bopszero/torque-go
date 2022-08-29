package utils

import (
	"reflect"

	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gopkg.in/yaml.v3"
)

func DumpData(fromObj interface{}, toObj interface{}) error {
	toValue := reflect.Indirect(reflect.ValueOf(toObj))
	if !toValue.CanSet() || !toValue.CanAddr() {
		return IssueErrorf("DumpData cannot set value to target object")
	}
	fromValue := reflect.Indirect(reflect.ValueOf(fromObj))
	if toValue.Type() != fromValue.Type() && toValue.Kind() != reflect.Interface {
		return IssueErrorf(
			"DumpData cannot dump from type %v to type %v",
			fromValue.Type(), toValue.Type(),
		)
	}
	toValue.Set(fromValue)
	return nil
}

func DumpDataByJSON(fromObject interface{}, toObject interface{}) error {
	fromJSON, err := comutils.JsonEncode(fromObject)
	if err != nil {
		return WrapError(err)
	}
	if err := comutils.JsonDecode(fromJSON, toObject); err != nil {
		return WrapError(err)
	}
	return nil
}

func DumpDataByYAML(fromObject interface{}, toObject interface{}) error {
	fromYAML, err := yaml.Marshal(fromObject)
	if err != nil {
		return WrapError(err)
	}
	if err := yaml.Unmarshal(fromYAML, toObject); err != nil {
		return WrapError(err)
	}
	return nil
}
