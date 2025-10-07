package utils

import (
	"strconv"
	"strings"
)

func EncodeSimpleString(s string) []byte {
	return []byte("+" + s + "\r\n")
}

func EncodeError(err string) []byte {
	return []byte("-" + err + "\r\n")
}

func EncodeInteger(i int) []byte {
	return []byte(":" + strconv.FormatInt(int64(i), 10) + "\r\n")
}

func EncodeBulkString(s string) []byte {
	if s == "" {
		return []byte("$-1\r\n")
	}
	return []byte("$" + strconv.Itoa(len(s)) + "\r\n" + s + "\r\n")
}

func EncodeArray(arr []string) []byte {
	var b strings.Builder
	b.WriteString("*" + strconv.Itoa(len(arr)) + "\r\n")
	for _, a := range arr {
		b.WriteString("$" + strconv.Itoa(len(a)) + "\r\n")
		b.WriteString(a + "\r\n")
	}
	return []byte(b.String())
}
