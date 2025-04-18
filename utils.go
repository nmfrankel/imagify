package main

import (
	"fmt"
	"strconv"
	"strings"
)

func (i *intSlice) String() string {
	return fmt.Sprintf("%v", *i)
}

func (i *intSlice) Set(value string) error {
	// Trim brackets if they exist
	value = strings.Trim(value, "[]")
	if value == "" {
		return nil
	}

	// Split string by comma
	strValues := strings.Split(value, ",")
	for _, str := range strValues {
		num, err := strconv.Atoi(strings.TrimSpace(str))
		if err != nil {
			return err
		}
		*i = append(*i, num)
	}
	return nil
}
