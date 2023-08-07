package util

import "strings"

var view_allow_exts = []string{".html", ".htm", ".shtml", ".xhtml"}

func GetViewExtsString() string {
	return strings.Join(view_allow_exts, ",")
}

func GetViewExtsSlice() []string {
	return view_allow_exts
}

func CheckExt(allowExts []string, ext string) bool {
	for _, allow_ext := range allowExts {
		if strings.EqualFold(allow_ext, ext) {
			return true
		}
	}
	return false
}
