package search

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/joaosoft/errors"
)

func GetEnv() string {
	env := os.Getenv("env")
	if env == "" {
		env = "local"
	}

	return env
}

func Exists(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func ReadFile(file string, obj interface{}) ([]byte, error) {
	var err error

	if !Exists(file) {
		return nil, errors.New(errors.ErrorLevel, 0, "file don't exist")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	if obj != nil {
		if err := json.Unmarshal(data, obj); err != nil {
			return nil, err
		}
	}

	return data, nil
}

func ReadFileLines(file string) ([]string, error) {
	lines := make([]string, 0)

	if !Exists(file) {
		return nil, errors.New(errors.ErrorLevel, 0, "file don't exist")
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func WriteFile(file string, obj interface{}) error {
	if !Exists(file) {
		return errors.New(errors.ErrorLevel, 0, "file don't exist")
	}

	jsonBytes, _ := json.MarshalIndent(obj, "", "    ")
	if err := ioutil.WriteFile(file, jsonBytes, 0644); err != nil {
		return err
	}

	return nil
}

func read(columns []string, rows *sql.Rows, value reflect.Value) (int, error) {

	value = value.Elem()
	isSlice := value.Kind() == reflect.Slice
	count := 0

	// load each row
	for rows.Next() {
		var elem reflect.Value
		if isSlice {
			elem = reflect.New(value.Type().Elem()).Elem()
		} else {
			elem = value
		}

		// load field values
		fields, err := getFields(loadOptionRead, columns, elem)
		if err != nil {
			return 0, err
		}

		// scan values from row
		err = rows.Scan(fields...)
		if err != nil {
			return 0, err
		}

		count++
		if isSlice {
			value.Set(reflect.Append(value, elem))
		} else {
			break
		}
	}

	return count, nil
}

func getFields(loadOption loadOption, columns []string, object reflect.Value) ([]interface{}, error) {
	var fields []interface{}

	// add columns to a map
	mapColumns := make(map[string]bool)
	for _, name := range columns {
		mapColumns[name] = true
	}

	mappedValues := make(map[string]interface{})
	loadColumnStructValues(loadOption, columns, mapColumns, object, mappedValues)

	for _, name := range columns {
		fields = append(fields, mappedValues[name])
	}

	return fields, nil
}

func loadColumnStructValues(loadOption loadOption, columns []string, mapColumns map[string]bool, object reflect.Value, mappedValues map[string]interface{}) {
	switch object.Kind() {
	case reflect.Ptr:
		if !object.IsNil() {
			loadColumnStructValues(loadOption, columns, mapColumns, object.Elem(), mappedValues)
		}
	case reflect.Struct:
		t := object.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" && !field.Anonymous {
				// unexported
				continue
			}
			tag := field.Tag.Get(string(loadOption))
			if tag == "-" {
				// ignore
				continue
			}

			if tag == "" {
				tag = field.Tag.Get(string(loadOptionDefault))
				if tag == "-" || tag == "" {
					// ignore
					continue
				}
			}

			if _, ok := mapColumns[tag]; ok {
				mappedValues[tag] = object.Field(i).Addr().Interface()
			}
		}
	default:
		mappedValues[columns[0]] = object.Addr().Interface()
	}
}

func loadStructValues(loadOption loadOption, object reflect.Value, columns *[]string, mappedValues map[string]reflect.Value) {
	switch object.Kind() {
	case reflect.Ptr:
		if !object.IsNil() {
			loadStructValues(loadOption, object.Elem(), columns, mappedValues)
		}
	case reflect.Struct:
		t := object.Type()
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.PkgPath != "" && !field.Anonymous {
				// unexported
				continue
			}

			tag := field.Tag.Get(string(loadOption))
			if tag == "-" {
				// ignore
				continue
			}

			if tag == "" {
				tag = field.Tag.Get(string(loadOptionDefault))
				if tag == "-" || tag == "" {
					// ignore
					continue
				}
			}

			if _, ok := mappedValues[tag]; !ok {
				mappedValues[tag] = object.Field(i)
				if columns != nil {
					*columns = append(*columns, tag)
				}
			}
		}
	}
}
