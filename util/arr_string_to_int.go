package util

import (
	"strconv"
	"strings"
)

func ArrStringToInt(arr []string) []uint {
	var res []uint
	for _, v := range arr {
		res = append(res, StringToInt(v))
	}
	return res
}

func StringToInt(str string) uint {
	res, _ := strconv.Atoi(str)
	return uint(res)
}

func StringToArrInt(str string, sep string) []uint {
	arr := strings.Split(str, sep)
	return ArrStringToInt(arr)
}

func ArrIntToSingleString(arr []uint, delimiter string) string {
	var res string
	for _, v := range arr {
		res += strconv.FormatInt(int64(v), 10) + delimiter
	}
	return res[:len(res)-1]
}
