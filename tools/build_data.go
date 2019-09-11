package tools

import (
	"fmt"
	"os"
	"strconv"
)

type DataGenerateRule struct {
	MaxKeyNum   int
	MaxStrLen   int
	MaxValLen   int
	MaxIntVal   int
	MaxFloatVal int
	MaxMulti    int
	Silent      bool
	Latency     bool
}

func getData(r *Rand, dgr *DataGenerateRule, prefix, name string) string {
	var (
		result string
		//记录start end情况下start的值
		first      = true
		firstInt   = -1
		firstChar  = 0
		firstFloat = -0.1
	)

	switch name {
	case "int":
		result = strconv.Itoa(r.RandInt(dgr.MaxIntVal))
	case "ilimit":
		iVal := r.RandInt(dgr.MaxIntVal)
		tmp := r.RandInt(2)
		if tmp%2 == 0 {
			iVal = -iVal
		}
		if first {
			firstInt = iVal
			first = false
		} else {
			if firstInt > 0 && iVal > 0 && firstInt > iVal {
				iVal += firstInt
			}
			first = true
		}
		result = strconv.Itoa(iVal)
	case "float":
		result = strconv.FormatFloat(r.RandFloat(dgr.MaxFloatVal), 'f', -1, 64)
	case "flimit":
		fVal := r.RandFloat(dgr.MaxFloatVal)
		if first {
			firstFloat = fVal
		} else {
			if firstFloat > fVal {
				fVal += firstFloat
			}
		}
		switch r.RandInt(3) {
		case 1:
			result = strconv.FormatFloat(fVal, 'f', -1, 64)
		case 2:
			result = "(" +
				strconv.FormatFloat(fVal, 'f', -1, 64)
		case 3:
			if first {
				result = "-inf"
				firstFloat = -0.1
			} else {
				result = "+inf"
			}
		}

		if first { //用于区分-inf和+inf
			first = false
		} else {
			first = true
		}
	case "string":
		result = r.RandRangeKey(prefix, dgr.MaxKeyNum)
	case "string_v":
		result = r.RandString(dgr.MaxValLen)
	case "slimit":
		iVal := r.RandInt(25)
		if first {
			if iVal == 25 {
				iVal--
			}
			firstChar = iVal
		} else {
			if firstChar >= iVal {
				iVal = firstChar + 1
			}
		}
		sVal := string('a'+iVal) + r.RandString(dgr.MaxStrLen)
		switch r.RandInt(3) {
		case 1:
			result = "(" + sVal
		case 2:
			result = "[" + sVal
		case 3:
			if first {
				result = "-"
				firstChar = 0
			} else {
				result = "+"
			}
		}
		if first {
			first = false
		} else {
			first = true
		}
	case "position":
		tmp := r.RandInt(2)
		if tmp%2 == 0 {
			result = "BEFORE"
		} else {
			result = "AFTER"
		}
	case "serialized":
		result = "\x0c\xc3@R@\x8b\x04\x8b\x00\x00\x00u \x03\n\x0c\x00\x00\x01" +
			"a\x03\x130.10\xe0\x05\x00\x031\x15\x01b@\x17\x002\xe0\x05\x16@\x17" +
			"\x00c`\x17\x009\xe0\x06\x00\x02\x15\x01d@\x17\x004 #\xe0\x03\x00\x0b2" +
			"\x15\x01e\x03\x030.5\x05\x01f@\x1f\x005\xe0\x066\x018\xff\x06\x006F_\xbe\xd0\xec*\x0b"
	case "match":
		length := r.RandInt(dgr.MaxStrLen)
		for length > 0 {
			length--
			switch r.RandInt(3) {
			case 1:
				result += string('a' + r.RandInt(26) - 1)
			case 2:
				result += "?"
			case 3:
				result += "*"
			}
		}
	case "aggrgate":
		switch r.RandInt(3) {
		case 1:
			result = "SUM"
		case 2:
			result = "MIN"
		case 3:
			result = "MAX"
		}
	case "sortModle":
		switch r.RandInt(2) {
		case 1:
			result = "ASC"
		case 2:
			result = "DESC"
		}
	default:
		result = name
	}

	return result
}

func BuildData(r *Rand, dgr *DataGenerateRule, prefix, conf interface{}, numKeys *int) []string {
	var (
		result []string
	)

	switch confType := conf.(type) {
	case []interface{}:
		paramList := conf.([]interface{})
		for ind, param := range paramList {
			switch param.(type) {
			case []interface{}:
				tmp := r.RandInt(2)
				if tmp%2 == 1 {
					result = append(result, BuildData(r, dgr, prefix, param, numKeys)...)
				}
			default:
				if param.(string) == "numkeys" {
					*numKeys = r.RandInt(dgr.MaxMulti) + 2
					result = append(result, strconv.Itoa(*numKeys))
				} else if param.(string) == "etc" {
					loop := 0
					if *numKeys != 0 {
						loop = *numKeys - 2 //保证key的个数与后续传入相同
					} else {
						loop = r.RandInt(dgr.MaxMulti)
					}
					for loop > 0 {
						result = append(result, BuildData(r, dgr, prefix, paramList[:ind], numKeys)...)
						loop--
					}
				} else {
					result = append(result, getData(r, dgr, prefix.(string), param.(string)))
				}
			}
		}
	default:
		fmt.Println("Unknown build data conf type:", confType)
		os.Exit(1)
	}

	return result
}
