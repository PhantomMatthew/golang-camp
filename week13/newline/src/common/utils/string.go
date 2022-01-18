package utils

import "strings"

func ReplaceLastWord(s, old, new string) string {
	if strings.HasSuffix(s,old) {
		return s[0 : len(s)-len(old)]+new
	}else{
		return s
	}
	//i := 0
	//for m := 1; m <= n; m++ {
	//	x := strings.Index(s[i:], old)
	//	if x < 0 {
	//		break
	//	}
	//	i += x
	//	if m == n {
	//		return s[:i] + new + s[i+len(old):]
	//	}
	//	i += len(old)
	//}
	//return s
}