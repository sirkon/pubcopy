package pubcopy

import (
	"log"
	"reflect"
	"runtime/debug"

	"github.com/pkg/errors"
)

// Copy copies src into dst using a filtering provided by filter implementation
func Copy(src interface{}, dst interface{}, filter Filter) (err error) {
	rSrc := reflect.ValueOf(src)
	rDst := reflect.ValueOf(dst)
	if rDst.Kind() != reflect.Ptr {
		return errors.Errorf("invalid destination, must be a pointer, got %T", dst)
	}
	if rSrc.Type() != rDst.Type().Elem() {
		return errors.Errorf("invalid destination, must be *%T, got %T", src, dst)
	}

	return pubcopy(rSrc, rDst, filter)
}

func pubcopy(src, dst reflect.Value, filter Filter) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = errors.Errorf("reflect error for %T %v:\n%s", src.Interface(), r, string(debug.Stack()))
		}
	}()
	typeOf := src.Type()
	switch typeOf.Kind() {
	case reflect.Ptr:
		res := reflect.New(typeOf.Elem())
		if err := pubcopy(src.Elem(), res, filter); err != nil {
			return errors.WithMessage(err, "pointer dereference")
		}
		dst.Elem().Set(res)
	case reflect.Struct:
		res := reflect.New(typeOf)
		for i := 0; i < typeOf.NumField(); i++ {
			field := typeOf.Field(i)
			if len(field.Name) == 0 {
				if !filter.Type(field.Type) {
					continue
				}
			}
			if !filter.Name(field.Name) {
				continue
			}

			nDst := reflect.New(field.Type)
			if err := pubcopy(src.Field(i), nDst, filter); err != nil {
				return errors.WithMessagef(err, "field %s of %T", field.Name, src.Interface())
			}
			res.Elem().Field(i).Set(nDst.Elem())
		}
		dst.Elem().Set(res.Elem())
	case reflect.Map:
		res := reflect.MakeMap(typeOf)
		iter := src.MapRange()
		for iter.Next() {
			if err := setMapItem(iter, filter, src, res); err != nil {
				return errors.WithMessage(err, "map iteration")
			}
		}
		if err := setMap(dst, res); err != nil {
			return errors.WithMessage(err, "map copy")
		}
	case reflect.Slice:
		if err := initSlice(dst, src); err != nil {
			return errors.WithMessage(err, "setting up a destination slice")
		}
		for i := 0; i < src.Len(); i++ {
			newItem := reflect.New(src.Type().Elem())
			if err := pubcopy(src.Index(i), newItem, filter); err != nil {
				return errors.WithMessagef(err, "copying slice %T element", src.Interface())
			}
			log.Println(dst.Type())
			if err := appendSliceItem(dst, newItem.Elem()); err != nil {
				return errors.WithMessage(err, "append slice item")
			}
		}
	case reflect.Array:
		dst.Elem().Set(src)
	case reflect.Chan:
		return errors.New("channels are now allowed to copy")
	default:
		dst.Elem().Set(src)
	}
	return nil
}

func initSlice(dst reflect.Value, src reflect.Value) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = errors.Errorf("%s", r)
		}
	}()

	dst.Elem().Set(reflect.New(src.Type()).Elem())
	return
}

func appendSliceItem(dst reflect.Value, newItem reflect.Value) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = errors.Errorf("%s", r)
		}
	}()
	dst.Elem().Set(reflect.Append(dst.Elem(), newItem))
	return
}

func setMap(dst reflect.Value, res reflect.Value) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = errors.Errorf("%s", r)
		}
	}()

	dst.Elem().Set(res)
	return
}

func setMapItem(iter *reflect.MapIter, filter Filter, src reflect.Value, res reflect.Value) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = errors.Errorf("%s", r)
		}
	}()

	nKey := reflect.New(iter.Key().Type())
	if err := pubcopy(iter.Key(), nKey, filter); err != nil {
		return errors.WithMessagef(err, "copying key of %T", src.Interface())
	}
	nValue := reflect.New(iter.Value().Type())
	if err := pubcopy(iter.Value(), nValue, filter); err != nil {
		return errors.WithMessagef(err, "copying value of %T for key %s", src.Interface(), iter.Key().Interface())
	}
	if err := setMapItemValue(res, nKey, nValue); err != nil {
		return errors.WithMessage(err, "setting map item")
	}
	return nil
}

func setMapItemValue(res reflect.Value, nKey reflect.Value, nValue reflect.Value) (err error) {
	defer func() {
		r := recover()
		if r != nil {
			err = errors.Errorf("%s", r)
		}
	}()
	res.SetMapIndex(nKey.Elem(), nValue.Elem())
	return
}
