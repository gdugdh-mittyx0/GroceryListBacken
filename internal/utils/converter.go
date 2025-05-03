package utils

import (
	"errors"
	"strconv"
)

var ErrConvert = errors.New("err_convert")

func StringToInt(str string) (int, error) {
	res, err := strconv.Atoi(str)
	if err != nil {
		return 0, ErrConvert
	}
	return res, nil
}

func StringToInt64(str string) (int64, error) {
	res, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, ErrConvert
	}
	return res, nil
}

func StringToUint(str string) (uint, error) {
	res, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, ErrConvert
	}
	return uint(res), nil
}

func StringToUint64(str string) (uint64, error) {
	res, err := strconv.ParseUint(str, 10, 64)
	if err != nil {
		return 0, ErrConvert
	}
	return res, nil
}

func StringToBool(str string) (bool, error) {
	res, err := strconv.ParseBool(str)
	if err != nil {
		return false, ErrConvert
	}
	return res, nil
}

func StringToFloat64(str string) (float64, error) {
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return 0, ErrConvert
	}
	return res, nil
}

func IntToString(num int) string {
	return strconv.Itoa(num)
}

func Int64ToString(num int64) string {
	return strconv.Itoa(int(num))
}

func StringsToUint32(strings []string) ([]uint32, error) {
	ints := make([]uint32, len(strings))

	for i, s := range strings {
		num, err := strconv.Atoi(s)
		if err != nil {
			return nil, ErrConvert
		}
		ints[i] = uint32(num)
	}

	return ints, nil
}
