package filter

import (
	"reflect"
	"regexp"
	"strings"
)

const (
	personalDataFilterTagName = "pdfilter"
	noFilterFlagName          = "nofilter"
	tagConfigSeparator        = ","
)

type personalDataFilter struct {
	mask                   string
	matchFilterFunc        *MatchFilterFunc
	personalDataRegExp     *regexp.Regexp
	personalDataProperties []string
}

func (filter *personalDataFilter) RemovePersonalData(input interface{}) interface{} {
	if input == nil {
		return nil
	}

	inputType := reflect.TypeOf(input)

	// We don't need to filter zero values.
	if reflect.DeepEqual(input, reflect.Zero(inputType).Interface()) {
		return input
	}

	switch inputType.Kind() {
	case reflect.String:
		return filter.handleString(input)
	case reflect.Slice:
		inputValue := reflect.ValueOf(input)
		res := reflect.MakeSlice(inputType, inputValue.Len(), inputValue.Cap())
		return filter.handleCollection(inputValue, res)
	case reflect.Array:
		// reflect.New will create pointer value. We need the dereferenced value.
		// The reflect.Ptr case will make sure to return pointer.
		res := reflect.New(inputType).Elem()
		inputArray := reflect.ValueOf(input)
		return filter.handleCollection(inputArray, res)
	case reflect.Map:
		return filter.handleMap(input)
	case reflect.Struct:
		return filter.handleStruct(input)
	case reflect.Ptr:
		return filter.handlePointer(input)
	default:
		return input
	}
}

func (filter *personalDataFilter) handleString(input interface{}) interface{} {
	var filtered string
	if filter.matchFilterFunc != nil {
		filtered = filter.personalDataRegExp.ReplaceAllStringFunc(input.(string), *filter.matchFilterFunc)
	} else {
		filtered = filter.personalDataRegExp.ReplaceAllString(input.(string), filter.mask)
	}

	return filtered
}

func (filter *personalDataFilter) handleCollection(input, res reflect.Value) interface{} {
	for i := 0; i < input.Len(); i++ {
		v := input.Index(i)
		filteredValue := filter.RemovePersonalData(v.Interface())
		res.Index(i).Set(reflect.ValueOf(filteredValue))
	}

	return res.Interface()
}

func (filter *personalDataFilter) handleMap(input interface{}) interface{} {
	// reflect.New will create pointer value. We need the dereferenced value.
	// The reflect.Ptr case will make sure to return pointer.
	mapValue := reflect.ValueOf(input)
	res := reflect.MakeMap(mapValue.Type())
	keys := mapValue.MapKeys()

	for _, k := range keys {
		v := mapValue.MapIndex(k)

		// The v will be interface for map[string]interface{} and
		// string for map[string]string. That's why we need to call the Interface() method.
		// The Interface() method will not return interface{} object which contains another interface{} object.
		// It will return interface{} object which holds a type. With this object we can treat map[string]interface{}
		// and map[string]string the same way.
		valueInterface := v.Interface()
		realValue := reflect.ValueOf(valueInterface)
		if k.Kind() == reflect.String && filter.isFieldPersonalDataString(realValue, k.Interface().(string)) {
			res.SetMapIndex(k, reflect.ValueOf(filter.mask))
		} else {
			filteredValue := filter.RemovePersonalData(valueInterface)
			res.SetMapIndex(k, reflect.ValueOf(filteredValue))
		}
	}

	return res.Interface()
}

func (filter *personalDataFilter) handleStruct(input interface{}) interface{} {
	inputValue := reflect.ValueOf(input)
	inputType := inputValue.Type()
	// reflect.New will create pointer value. We need the dereferenced value.
	// The reflect.Ptr case will make sure to return pointer.
	inputValueCopy := reflect.New(inputType).Elem()

	for i := 0; i < inputValue.NumField(); i++ {
		field := inputType.Field(i)
		if len(field.PkgPath) > 0 {
			// The field is private (https://golang.org/pkg/reflect/#StructField).
			// Currently we can't set unexported fields - https://golang.org/pkg/reflect/#Value.CanSet.
			continue
		}

		fieldConfig := filter.getTagConfig(field)
		fieldValue := inputValue.Field(i)
		if fieldConfig != nil && fieldConfig.NoFilter {
			inputValueCopy.Field(i).Set(fieldValue)
			continue
		}

		var filteredField interface{}
		if filter.isFieldPersonalDataString(fieldValue, field.Name) {
			filteredField = filter.mask
		} else {
			filteredField = filter.RemovePersonalData(fieldValue.Interface())
		}

		inputValueCopy.Field(i).Set(reflect.ValueOf(filteredField))
	}

	return inputValueCopy.Interface()
}

func (filter *personalDataFilter) handlePointer(input interface{}) interface{} {
	derefedInterface := reflect.ValueOf(input).Elem().Interface()
	// The value returned from RemovePersonalData method will not be pointer.
	// It is used only in this method and that's why we can't call the
	// Addr method on it. It is not addressable (https://golang.org/pkg/reflect/#Value.CanAddr).
	value := filter.RemovePersonalData(derefedInterface)

	reflectValue := reflect.ValueOf(value)
	reflectType := reflect.TypeOf(value)
	// We need to return pointer in this case. That's why we need to make the value
	// returned from the RemovePersonalData method addressable (https://golang.org/pkg/reflect/#Value.CanAddr).
	// Since we don't know what type we need, we can only use dynamic slice. Each element in the dynamic slice
	// is addressable.
	slice := reflect.MakeSlice(reflect.SliceOf(reflectType), 1, 1)
	slice.Index(0).Set(reflectValue)

	return slice.Index(0).Addr().Interface()
}

func (filter *personalDataFilter) isFieldPersonalDataString(value reflect.Value, fieldName string) bool {
	return value.Kind() == reflect.String && indexOfString(filter.personalDataProperties, strings.ToLower(fieldName)) >= 0
}

func (filter *personalDataFilter) getTagConfig(field reflect.StructField) *filterTagConfig {
	config, ok := field.Tag.Lookup(personalDataFilterTagName)
	if !ok {
		return nil
	}

	values := strings.Split(config, tagConfigSeparator)
	noFilter := indexOfString(values, noFilterFlagName) >= 0

	return &filterTagConfig{
		NoFilter: noFilter,
	}
}
