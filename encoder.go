package bencoder

import (
	"bytes"
	"maps"
	"slices"
	"sort"
	"strconv"
)

type encoder struct {
	bb bytes.Buffer
}

// Целое число записывается так: i<число в десятичной системе счисления>e. Число не должно начинаться с нуля, но число
// нуль записывается как i0e. Отрицательные числа записываются со знаком минуса перед числом.
// Число −42 будет выглядеть так «i-42e».
func (e *encoder) encodeInt(i int64) {
	e.bb.WriteByte('i')
	e.bb.WriteString(strconv.FormatInt(i, 10))
	e.bb.WriteByte('e')
}

func (e *encoder) encodeUint(i uint64) {
	e.bb.WriteByte('i')
	e.bb.WriteString(strconv.FormatUint(i, 10))
	e.bb.WriteByte('e')
}

// Строка байт: <размер>:<содержимое>. Размер — это положительное число в десятичной системе счисления, может быть
// нулём; содержимое — это непосредственно данные, представленные цепочкой байт, которые не подразумевают
// никакой символьной кодировки. Строка «spam» в этом формате будет выглядеть так «4:spam».
func (e *encoder) encodeString(s string) {
	//e.bb.WriteString(strconv.Itoa(utf8.RuneCountInString(s)))
	e.bb.WriteString(strconv.Itoa(len(s))) // именно кол-во байт
	e.bb.WriteByte(':')
	e.bb.WriteString(s)
}

// Список (массив): l<содержимое>e. Содержимое включает в себя любые типы Bencode, следующие друг за другом.
// Список, состоящий из строки «spam» и числа 42, будет выглядеть так: «l4:spami42ee».
func (e *encoder) encodeArray(arr []any) {
	e.bb.WriteByte('l')

	for _, a := range arr {
		e.encodeUnknown(a)
	}

	e.bb.WriteByte('e')
}

// Словарь: d<содержимое>e. Содержимое состоит из пар ключ-значение, которые следуют друг за другом.
// Ключи могут быть только строкой байт и должны быть упорядочены в лексикографическом порядке.
// Значение может быть любым элементом Bencode. Если сопоставить ключам «bar» и «foo» значения «spam» и 42,
// получится: «d3:bar4:spam3:fooi42ee». (Если добавить пробелы между элементами, будет легче понять
// структуру: "d 3:bar 4:spam 3:foo i42e e".)
func (e *encoder) encodeDict(dict map[string]any) {
	// RFC requires that keys be sorted
	keys := sort.StringSlice(slices.Collect(maps.Keys(dict)))
	keys.Sort()

	e.bb.WriteByte('d')
	for _, key := range keys {
		val := dict[key]
		e.encodeString(key)
		e.encodeUnknown(val)
	}
	e.bb.WriteByte('e')
}

func (e *encoder) encodeUnknown(val any) {
	switch v := val.(type) {
	case string:
		e.encodeString(v)
	case int:
		e.encodeInt(int64(v))
	case int8:
		e.encodeInt(int64(v))
	case int16:
		e.encodeInt(int64(v))
	case int32:
		e.encodeInt(int64(v))
	case int64:
		e.encodeInt(v)
	case uint:
		e.encodeUint(uint64(v))
	case uint8:
		e.encodeUint(uint64(v))
	case uint16:
		e.encodeUint(uint64(v))
	case uint32:
		e.encodeUint(uint64(v))
	case uint64:
		e.encodeUint(v)
	case []string:
		e.encodeArray(toAnySlice(v))
	case []int:
		e.encodeArray(toAnySlice(v))
	case []int8:
		e.encodeArray(toAnySlice(v))
	case []int16:
		e.encodeArray(toAnySlice(v))
	case []int32:
		e.encodeArray(toAnySlice(v))
	case []int64:
		e.encodeArray(toAnySlice(v))
	case []uint:
		e.encodeArray(toAnySlice(v))
	case []uint8:
		e.encodeArray(toAnySlice(v))
	case []uint16:
		e.encodeArray(toAnySlice(v))
	case []uint32:
		e.encodeArray(toAnySlice(v))
	case []uint64:
		e.encodeArray(toAnySlice(v))
	case []any:
		e.encodeArray(v)
	case map[string]any:
		e.encodeDict(v)
	}
}

func toAnySlice[T any](sl []T) []any {
	res := make([]any, len(sl))
	for i, v := range sl {
		res[i] = v
	}
	return res
}

func (e *encoder) bytes() []byte {
	return e.bb.Bytes()
}
