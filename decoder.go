package bencoder

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type decoder struct {
	rdr *bufio.Reader
}

// Целое число записывается так: i<число в десятичной системе счисления>e. Число не должно начинаться с нуля, но число
// нуль записывается как i0e. Отрицательные числа записываются со знаком минуса перед числом.
// Число −42 будет выглядеть так «i-42e».
//
//	Не используем метод, так как нужен generic (Method cannot have type parameters).

func (d *decoder) decodeInteger() (any, error) {
	// check first symbol
	if b, err := d.rdr.ReadByte(); err != nil {
		return nil, err
	} else if b != 'i' {
		return nil, fmt.Errorf("integer must start with 'i'")
	}

	s, err := d.rdr.ReadString('e')
	if err != nil {
		return 0, err
	}

	if len(s) < 2 {
		return 0, fmt.Errorf("invalid integer: %s", s)
	}

	s = s[:len(s)-1]

	n, err := strconv.ParseInt(s, 10, 64)
	if err == nil {
		return n, nil
	}

	if errors.Is(err, strconv.ErrRange) {
		return strconv.ParseUint(s, 10, 64)
	}

	return 0, err
}

// Строка байт: <размер>:<содержимое>. Размер — это положительное число в десятичной системе счисления, может быть
// нулём; содержимое — это непосредственно данные, представленные цепочкой байт, которые не подразумевают
// никакой символьной кодировки. Строка «spam» в этом формате будет выглядеть так «4:spam».
func (d *decoder) decodeString() (string, error) {
	s, err := d.rdr.ReadString(':')
	if err != nil {
		return "", err
	}

	s = s[:len(s)-1]

	length, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return "", fmt.Errorf("string must start with length %w", err)
	}

	buf := make([]byte, length)
	n, err := io.ReadFull(d.rdr, buf)

	need := length - uint64(n)
	if need != 0 {
		return "", fmt.Errorf("could not read %d bytes", need)
	}

	return string(buf), nil
}

// Список (массив): l<содержимое>e. Содержимое включает в себя любые типы Bencode, следующие друг за другом.
// Список, состоящий из строки «spam» и числа 42, будет выглядеть так: «l4:spami42ee».
func (d *decoder) decodeArray() ([]any, error) {
	// check first symbol
	if b, err := d.rdr.ReadByte(); err != nil {
		return nil, err
	} else if b != 'l' {
		return nil, fmt.Errorf("array must start with 'l'")
	}

	var res []any

	for {
		val, err := d.decodeUnknown()
		if err != nil {
			return nil, err
		}

		res = append(res, val)

		next, err := d.peekNextByte()
		if err != nil {
			return nil, err
		} else if next == 'e' {
			break
		}
	}

	return res, nil
}

// Словарь: d<содержимое>e. Содержимое состоит из пар ключ-значение, которые следуют друг за другом.
// Ключи могут быть только строкой байт и должны быть упорядочены в лексикографическом порядке.
// Значение может быть любым элементом Bencode. Если сопоставить ключам «bar» и «foo» значения «spam» и 42,
// получится: «d3:bar4:spam3:fooi42ee». (Если добавить пробелы между элементами, будет легче понять
// структуру: "d 3:bar 4:spam 3:foo i42e e".)
func (d *decoder) decodeDict() (map[string]any, error) {
	// check first symbol
	if b, err := d.rdr.ReadByte(); err != nil {
		return nil, err
	} else if b != 'd' {
		return nil, fmt.Errorf("array must start with 'd'")
	}

	res := map[string]any{}

	for {
		key, err := d.decodeString()
		if err != nil {
			return nil, err
		}

		val, err := d.decodeUnknown()
		if err != nil {
			return nil, err
		}

		res[key] = val

		next, err := d.peekNextByte()
		if err != nil {
			return nil, err
		} else if next == 'e' {
			break
		}
	}

	return res, nil
}

func (d *decoder) decodeUnknown() (any, error) {
	var err error

	next, err := d.peekNextByte()
	if err != nil {
		return nil, err
	}

	var res any

	switch true {
	case next >= '1' && next <= '9':
		res, err = d.decodeString()
	case next == 'i':
		res, err = d.decodeInteger()
	case next == 'l':
		res, err = d.decodeArray()
	case next == 'd':
		res, err = d.decodeDict()
	}
	return res, err
}

//func parseTypedInteger[T int64 | uint64](s string) (n T, err error) {
//	// интересный механизм: обращение к result и приведение к интерфейсу
//	switch any(n).(type) {
//	case int64:
//		val, e := strconv.ParseInt(s, 10, 64)
//		if e != nil {
//			err = e
//			return
//		}
//
//		n = T(val)
//	case uint64:
//		val, e := strconv.ParseUint(s, 10, 64)
//		if e != nil {
//			err = e
//			return
//		}
//
//		n = T(val)
//	}
//
//	// более понятно - return n, err
//	return n, err
//}

func (d *decoder) peekNextByte() (byte, error) {
	bb, err := d.rdr.Peek(1)
	if err != nil {
		return 0, err
	}

	return bb[0], nil
}
