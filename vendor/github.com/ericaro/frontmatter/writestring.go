package frontmatter

import (
	"reflect"
)

var (
	stringType = reflect.ValueOf("txt").Type()
)

//Read a String out of a field identified by tag/tagvalue pair
func ReadString(v interface{}, tag, tagval string) (txt string, err error) {
	fv := findField(v, tag, tagval) //find the field
	if fv == nil {                  //there was no such field, report it as an error
		return "", ErrNoContentField
	}
	//convert to string type ()
	if !fv.Type().ConvertibleTo(stringType) {
		return "", ErrWrongContentFieldType //has the wrong type
	}
	//conversion has been checked above
	txt = fv.Convert(stringType).Interface().(string)
	return

}

//WriteString set 'txt' on the field tagged with 'tag'
func WriteString(v interface{}, tag, tagval, txt string) error {
	fv := findField(v, tag, tagval)
	if fv == nil {
		return ErrNoContentField
	}

	if !fv.CanSet() {
		return ErrUnexported
	}

	fvt := fv.Type()
	//attempt the only two cases: string convertible to field type
	// or *string is
	switch {

	case reflect.TypeOf(txt).ConvertibleTo(fvt):
		fv.Set(reflect.ValueOf(txt).Convert(fvt))

	case reflect.TypeOf(&txt).ConvertibleTo(fvt):
		fv.Set(reflect.ValueOf(&txt).Convert(fvt))

	default:
		return ErrWrongContentFieldType
	}
	return nil
}

//findField by tagkey tagvalue pair
func findField(v interface{}, tag, tagval string) *reflect.Value {

	val := reflect.ValueOf(v)
	ty := val.Type()

	// if it was a pointer to a struct open the underlying type, instead
	if ty.Kind() == reflect.Ptr {
		ty = ty.Elem() // the actual underneath type
		val = val.Elem()
	}

	for i := 0; i < ty.NumField(); i++ { //for each field, (yeah, I know it's a bit gross)
		//get the field, and the fields tag
		if x := ty.Field(i).Tag.Get(tag); x == tagval {
			f := val.Field(i) // get the field value
			return &f         //and return it
		}
	}
	return nil
}
