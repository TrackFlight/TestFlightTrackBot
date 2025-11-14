package testflight

import "regexp"

var (
	RegexAppName = regexp.MustCompile(
		`<title>Join the ([^<]+?) beta - TestFlight - Apple</title>`,
	)
	RegexAppIcon = regexp.MustCompile(
		`background-image: url\((https://[A-Za-z0-9./_-]+)\);`,
	)
	RegexDescription = regexp.MustCompile(
		`<p class="step3">([^<]*?)</p>`,
	)
	RegexStatus = regexp.MustCompile(
		`<div class="beta-status">\n?.*?\n?\s*<span>(.+)</span>\n?\s*</div>`,
	)
	RegexLink = regexp.MustCompile(
		`^https?://testflight\.apple\.com/join/[A-Za-z0-9]+$`,
	)
	RegexAwesomeTestFlight = regexp.MustCompile(
		`(^|\n)\|\s*(.*?)\s*\|\s*\[(https?://testflight\.apple\.com/join/\w+)][^|]*\s*\|\s*([YFND])`,
	)
	RegexImageSize = regexp.MustCompile(
		`\d+x\d+\w*-?\d*`,
	)
)
