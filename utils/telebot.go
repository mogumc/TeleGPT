// 文本预处理
// @author MoGuQAQ
// @version 1.0.0

package utils

import (
	"fmt"
	"regexp"
)

func Contains(slice []int64, item int64) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func Is_command(str string, command string, uname string) bool {
	pattern := fmt.Sprintf(`^\/%s(?:@%s)?(?:\s+\S.*)?$`, command, uname)
	re := regexp.MustCompile(pattern)
	if re.MatchString(str) {
		return true
	} else {
		return false
	}
}

func Get_command(str string, command string, uname string) string {
	pattern := fmt.Sprintf(`^\/%s(@%s)?\s+(.*)$`, command, uname)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(str)
	if len(matches) > 2 {
		return matches[2]
	} else {
		return ""
	}
}

func Is_at(str string, uname string) bool {
	pattern := fmt.Sprintf(`^@%s(?:\s+\S.*)?$`, uname)
	re := regexp.MustCompile(pattern)
	if re.MatchString(str) {
		return true
	} else {
		return false
	}
}

func Get_at(str string, uname string) string {
	pattern := fmt.Sprintf(`^@%s?\s+(.*)$`, uname)
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(str)
	if len(matches) > 2 {
		return matches[2]
	} else {
		return ""
	}
}
