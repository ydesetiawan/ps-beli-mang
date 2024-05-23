package helper

import (
	"fmt"
	"strings"
)

func ConvertStringListToString(list []string, separator string) string {
	if list == nil || len(list) == 0 {
		return ""
	}
	strList := make([]string, len(list))
	for i, v := range list {
		strList[i] = fmt.Sprintf("%s", v)
	}
	result := strings.Join(strList, separator)
	return result
}

func ConvertStringListToQuotedString(list []string, separator string) string {
	if list == nil || len(list) == 0 {
		return ""
	}
	strList := make([]string, len(list))
	for i, v := range list {
		strList[i] = fmt.Sprintf("'%s'", v)
	}
	result := strings.Join(strList, separator)
	return result
}

func ConvertIntListToString(intList []int, separator string) string {
	if intList == nil || len(intList) == 0 {
		return ""
	}
	strList := make([]string, len(intList))
	for i, v := range intList {
		strList[i] = fmt.Sprintf("%d", v)
	}
	result := strings.Join(strList, separator)
	return result
}

func ConvertInt64ListToString(intList []int64, separator string) string {
	if intList == nil || len(intList) == 0 {
		return ""
	}
	strList := make([]string, len(intList))
	for i, v := range intList {
		strList[i] = fmt.Sprintf("%d", v)
	}
	result := strings.Join(strList, separator)
	return result
}

func ConvertIntListToInt64(intList []int) []int64 {
	if intList == nil || len(intList) == 0 {
		return []int64{}
	}
	result := make([]int64, len(intList))
	for i, v := range intList {
		result[i] = int64(v)
	}
	return result
}

func ConvertInt64ListToInt(int64List []int64) []int {
	if int64List == nil || len(int64List) == 0 {
		return []int{}
	}
	result := make([]int, len(int64List))
	for i, v := range int64List {
		result[i] = int(v)
	}
	return result
}

func ConvertToUnitName(unitId int) string {
	switch unitId {
	case 1:
		return "Pcs"
	case 2:
		return "Kg"
	case 3:
		return "Ltr"
	case 4:
		return "Mtr"
	case 5:
		return "Set"
	case 6:
		return "Box"
	}
	return ""
}
