package hasher

import (
	"io"
	"reflect"

	"github.com/cespare/xxhash/v2"
	"github.com/webmafia/fast"
)

var _ io.Writer = (*Hasher)(nil)

type Hasher struct {
	dig xxhash.Digest
}

func (h *Hasher) Reset() {
	h.dig.Reset()
}

func (h *Hasher) Hash() uint64 {
	return h.dig.Sum64()
}

//go:inline
func (h *Hasher) Write(b []byte) (int, error) {
	return h.dig.Write(b)
}

func (h *Hasher) WriteString(s string) (int, error) {
	return h.Write(fast.StringToBytes(s))
}

func (h *Hasher) WriteInt(v int) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteInt8(v int8) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteInt16(v int16) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteInt32(v int32) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteInt64(v int64) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteUint(v uint) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteUint8(v uint8) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteUint16(v uint16) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteUint32(v uint32) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteUint64(v uint64) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteFloat32(v float32) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteFloat64(v float64) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteBool(v bool) (int, error) {
	return h.Write(toBytes(&v))
}

func (h *Hasher) WriteAny(v any) (int, error) {
	switch value := v.(type) {
	case int:
		return h.WriteInt(value)
	case int8:
		return h.WriteInt8(value)
	case int16:
		return h.WriteInt16(value)
	case int32:
		return h.WriteInt32(value)
	case int64:
		return h.WriteInt64(value)
	case uint:
		return h.WriteUint(value)
	case uint8:
		return h.WriteUint8(value)
	case uint16:
		return h.WriteUint16(value)
	case uint32:
		return h.WriteUint32(value)
	case uint64:
		return h.WriteUint64(value)
	case float32:
		return h.WriteFloat32(value)
	case float64:
		return h.WriteFloat64(value)
	case bool:
		return h.WriteBool(value)
	case string:
		return h.WriteString(value)
	case []byte:
		return h.Write(value)
	default:
		// For other types, use reflection to handle structs, slices, arrays, and pointers.
		return h.writeReflect(v)
	}
}

func (h *Hasher) writeReflect(v any) (int, error) {

	// Use reflection to determine the kind of value.
	val := reflect.ValueOf(v)
	switch val.Kind() {

	case reflect.Pointer:
		// If it's a pointer, dereference it and hash the value it points to.
		if val.IsNil() {
			return h.Write([]byte{})
		}
		return h.WriteAny(val.Elem().Interface())

	case reflect.Struct:
		// For structs, iterate through fields and hash them individually.
		total := 0
		for i := 0; i < val.NumField(); i++ {
			fld := val.Field(i)

			if !fld.CanInterface() {
				continue
			}

			n, err := h.WriteAny(fld.Interface())
			if err != nil {
				return total, err
			}
			total += n
		}
		return total, nil

	case reflect.Slice, reflect.Array:
		// For slices or arrays, iterate through the elements and hash them.
		total := 0
		for i := 0; i < val.Len(); i++ {
			n, err := h.WriteAny(val.Index(i).Interface())
			if err != nil {
				return total, err
			}
			total += n
		}
		return total, nil
	}

	return 0, nil
}

func Hash(v any) uint64 {
	var h Hasher
	h.Reset()
	h.WriteAny(v)
	return h.Hash()
}
