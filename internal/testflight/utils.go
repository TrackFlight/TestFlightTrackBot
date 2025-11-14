package testflight

import "fmt"

func ChangeLinkResolution(link string, size int) string {
	return RegexImageSize.ReplaceAllString(link, fmt.Sprintf("%dx%d", size, size))
}
