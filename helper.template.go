package main

import (
	"fmt"
	"strconv"
	"time"
)

// 格式化时间
func DateFormat(date time.Time, layout string) string {
	return date.Format(layout)
}

// 截取字符串
func Substring(source string, start, end int) string {
	rs := []rune(source)
	length := len(rs)
	if start < 0 {
		start = 0
	}
	if end > length {
		end = length
	}
	return string(rs[start:end])
}

func Str2Float(str string) float64 {

	f, _ := strconv.ParseFloat(str, 32)
	return f
}

func Str2Int(str string) int {

	i, _ := strconv.Atoi(str)
	return i
}

func FloatLt(str string) bool {
	if Str2Float(str) < 10.0 {
		return true
	}
	return false
}

func IntEq(val int) bool {
	if val == 0 {
		return true
	}
	return false
}

func ToDiscount(money string, discount string) string {
	fmoney := Str2Float(money)
	fdiscount := Str2Float(discount)
	if Str2Float(discount) < 10.0 {
		return fmt.Sprintf("%.0f", fmoney*fdiscount/10)
	}
	return fmt.Sprintf("%.0f", fmoney)
}

func Sub(money string, discount string) string {
	fmoney := Str2Float(money)
	fdiscount := Str2Float(discount)
	if Str2Float(discount) < 10.0 {
		return fmt.Sprintf("%.0f", fmoney-fmoney*fdiscount/10)
	}
	return "0.00"
}
