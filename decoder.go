package bencoder

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type decoder struct {
	rd *bufio.Reader
}

// Целое число записывается так: i<число в десятичной системе счисления>e. Число не должно начинаться с нуля, но число
// нуль записывается как i0e. Отрицательные числа записываются со знаком минуса перед числом.
// Число −42 будет выглядеть так «i-42e».
func (d *decoder) decodeInt() (int64, error) {
	s, err := d.rd.ReadString('e')
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(s[:len(s)-1], 10, 64)
}

func (d *decoder) decodeUint() (uint64, error) {
	s, err := d.rd.ReadString('e')
	if err != nil {
		return 0, err
	}
	return strconv.ParseUint(s[:len(s)-1], 10, 64)
}

// Строка байт: <размер>:<содержимое>. Размер — это положительное число в десятичной системе счисления, может быть
// нулём; содержимое — это непосредственно данные, представленные цепочкой байт, которые не подразумевают
// никакой символьной кодировки. Строка «spam» в этом формате будет выглядеть так «4:spam».
func (d *decoder) decodeString() (string, error) {
	s, err := d.rd.ReadString(':')
	if err != nil {
		return "", err
	}

	length, err := strconv.ParseUint(s[:len(s)-1], 10, 64)
	if err != nil {
		return "", err
	}

	buf := make([]byte, length)
	n, err := io.ReadFull(d.rd, buf)

	need := length - uint64(n)
	if need != 0 {
		return "", fmt.Errorf("could not read %d bytes", need)
	}

	return string(buf), nil
}

// Список (массив): l<содержимое>e. Содержимое включает в себя любые типы Bencode, следующие друг за другом.
// Список, состоящий из строки «spam» и числа 42, будет выглядеть так: «l4:spami42ee».
func (d *decoder) decodeArray() ([]any, error) {
	b, err := d.rd.ReadByte()
	if err != nil {
		return nil, err
	}

	if b != 'l' {
		return nil, fmt.Errorf("array must be started from 'l'")
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
	res := map[string]any{}

	for {
		// d.rd.Peek - нельзя, потому что контент может содержать 'e' и надо читать валидными блоками
		err := d.rd.UnreadByte()
		if err != nil {
			return nil, err
		}

		key, err := d.decodeUnknown()
		if err != nil {
			return nil, err
		}

		val, err := d.decodeUnknown()
		if err != nil {
			return nil, err
		}

		res[key.(string)] = val

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
		res, err = d.decodeInt()
		if err != nil {
			// overflow checking
			res, err = d.decodeUint()
		}
	case next == 'l':
		res, err = d.decodeArray()
	case next == 'd':
		res, err = d.decodeDict()
	}

	return res, err
}

func (d *decoder) peekNextByte() (byte, error) {
	b, err := d.rd.ReadByte()
	if err != nil {
		return 0, err
	}

	err = d.rd.UnreadByte()
	if err != nil {
		return 0, err
	}

	return b, nil
}
