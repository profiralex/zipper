package zipper

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
)

const (
	zipperTag        = "zip"
	zipperDefaultTag = "zip-default"
)

// Provider provides data from a certain source
type Provider interface {
	Lookup(string) (string, error)
}

// Load retrieves the value for struct keys from Provider
// zip tag defines the value key to lookup
// zip-default tag defines the default value if value not found by provider
// if strict parameter is true then the process stops if loader cannot find a key and
// no default value provided
func Load(model interface{}, provider Provider, strict bool) error {
	return load(model, provider, strict)
}

func load(model interface{}, provider Provider, strict bool) error {
	t := reflect.TypeOf(model).Elem()
	v := reflect.ValueOf(model).Elem()

	for i := 0; i < t.NumField(); i++ {
		typeField := t.Field(i)
		valueField := v.Field(i)

		if typeField.Type.Kind() == reflect.Struct {
			value := valueField.Addr().Interface()
			err := load(value, provider, strict)
			if err != nil {
				return fmt.Errorf("Unable to set value for %s: %w", typeField.Name, err)
			}
			continue
		}

		key, ok := typeField.Tag.Lookup(zipperTag)
		if !ok {
			continue
		}

		value, err := provider.Lookup(key)
		if err != nil {
			defaultValue, ok := typeField.Tag.Lookup(zipperDefaultTag)
			if !ok {
				if strict {
					return fmt.Errorf("Unable to load value for field %s: %w", typeField.Name, err)
				}
				continue
			}

			value = defaultValue
		}

		err = assignValue(&v, &typeField, value)
		if err != nil {
			return fmt.Errorf("Unable to set value for field %s: %w", typeField.Name, err)
		}
	}

	return nil
}

func assignValue(reflectValue *reflect.Value, typeField *reflect.StructField, value string) error {
	switch typeField.Type.Kind() {
	case reflect.String:
		reflectValue.FieldByIndex(typeField.Index).SetString(value)
	case reflect.Bool:
		value, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		reflectValue.FieldByIndex(typeField.Index).SetBool(value)
	case reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64:
		value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		reflectValue.FieldByIndex(typeField.Index).SetInt(value)
	case reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64:
		value, err := strconv.ParseUint(value, 10, 64)
		if err != nil {
			return err
		}
		reflectValue.FieldByIndex(typeField.Index).SetUint(value)
	case reflect.Float32, reflect.Float64:
		value, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		reflectValue.FieldByIndex(typeField.Index).SetFloat(value)
	default:
		return fmt.Errorf("Could not set value for %s type", typeField.Type.Name())
	}
	return nil
}

// EnvProvider loads values from environment
type EnvProvider struct {
}

// Lookup looks up for value in the environment
func (p *EnvProvider) Lookup(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("Value not found in the environment")
	}

	return value, nil
}
