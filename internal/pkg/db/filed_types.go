package db

import (
	"database/sql/driver"
	"strconv"
	"strings"
)

type Uint64s []uint64

func (u Uint64s) Value() (driver.Value, error) {
	var str string
	size := len(u)
	for idx, i := range u {
		str += strconv.FormatUint(i, 10)
		if idx + 1 < size {
			str += ","
		}
	}
	return "[" + str + "]", nil
}

func (u *Uint64s) Scan(src interface{}) error {
	str := strings.TrimRight(src.(string), "]")
	str = strings.TrimLeft(str, "[")
	strToList := strings.Split(str, ",")
	for _, c := range strToList {
		i, err := strconv.ParseUint(c, 10, 64)
		if err != nil {
			return err
		}
		*u = append(*u, i)
	}
	return nil
}

type Strings []string

func (s Strings) Value() (driver.Value, error) {
	var str string
	size := len(s)
	for idx, c := range s {
		str += c
		if idx + 1 < size {
			str += ","
		}
	}
	return "[" + str + "]", nil
}

func (s *Strings) Scan(src interface{}) error {
	str := strings.TrimRight(src.(string), "]")
	str = strings.TrimLeft(str, "[")
	*s = strings.Split(str, ",")
	return nil
}
