package jsonlines

import (
	"fmt"
	"reflect"
	"io"
	"bufio"
	"encoding/json"
)

func getOriginalSlice(ptrToSlice interface{}) (slice reflect.Value, err error) {
	ptr2sl := reflect.TypeOf(ptrToSlice)
	if ptr2sl.Kind() != reflect.Ptr {
    		return reflect.ValueOf(nil), fmt.Errorf("expected pointer to slice, got %s", ptr2sl.Kind())
	}

	originalSlice := reflect.Indirect(reflect.ValueOf(ptrToSlice))
	sliceType := originalSlice.Type()
	if sliceType.Kind() != reflect.Slice {
    		return reflect.ValueOf(nil), fmt.Errorf("expected pointer to slice, got pointer to %s", sliceType.Kind())
	}
	return originalSlice, nil
}

// Decode reads the next JSON Lines-encoded value that reads
// from r and stores it in the slice pointed to by ptrToSlice.
func Decode(r io.Reader, ptrToSlice interface{}) error {
	originalSlice, err := getOriginalSlice(ptrToSlice)
	if err != nil {
		return err
	}

	slElem := originalSlice.Type().Elem()
	scanner := bufio.NewReader(r)
	for {
		item, err := scanner.ReadBytes('\n')
		if err != nil {
			log.Print(string(item))
			return err
		}

		//create new object
		newObj := reflect.New(slElem).Interface()
		err = json.Unmarshal(item, newObj)

		ptrToNewObj := reflect.Indirect(reflect.ValueOf(newObj))
		originalSlice.Set(reflect.Append(originalSlice, ptrToNewObj))

		if err != nil {
			log.Print(string(item))
			return err
		}
	}
}

// Encode writes the JSON Lines encoding of ptrToSlice to the w stream
func Encode(w io.Writer, ptrToSlice interface{}) error {
	originalSlice, err := getOriginalSlice(ptrToSlice)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(w)
	for i := 0; i < originalSlice.Len(); i++ {
		elem := originalSlice.Index(i).Interface()
		err = enc.Encode(elem)
		if err != nil {
			return err
		}
        }
	return nil
}
