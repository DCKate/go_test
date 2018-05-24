package kkutil

import (
	"fmt"
	"reflect"
)

//KKTag use for structure tag ex. `kk:"name"
const KKTag = "kk"

func serialize(mm map[string]interface{}, to interface{}) reflect.Value {
	st := reflect.TypeOf(to)
	newPtr := reflect.New(st)
	vl := newPtr.Elem()
	fmt.Println(st)
	fmt.Println(vl.Kind())

	if vl.Kind() == reflect.Struct {
		fmt.Println(vl.Kind())
		for ii := 0; ii < st.NumField(); ii++ {
			key := st.Field(ii).Tag.Get(KKTag)
			if vv, ok := mm[key]; ok {
				// exported field
				f := vl.FieldByName(st.Field(ii).Name)
				if f.IsValid() {
					// A Value can be changed only if it is
					// addressable and was not obtained by
					// the use of unexported struct fields.
					fmt.Println(f.CanSet())
					if f.CanSet() {
						switch assVal := vv.(type) {
						case int:
							if f.Kind() == reflect.Int {
								x := int64(assVal)
								if !f.OverflowInt(x) {
									vl.FieldByName(st.Field(ii).Name).SetInt(x)
								}
							}
						case string:
							if f.Kind() == reflect.String {
								vl.FieldByName(st.Field(ii).Name).SetString(assVal)
							}
						}
					}
				}
			}
		}
	}

	return newPtr
}
