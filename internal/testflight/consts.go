package testflight

import "regexp"

const UserAgentListURL = "https://raw.githubusercontent.com/fake-useragent/fake-useragent/refs/heads/main/src/fake_useragent/data/browsers.jsonl"
const AwesomeTestFlightURL = "https://raw.githubusercontent.com/pluwen/awesome-testflight-link/refs/heads/main/README.md"

var RegexAwesomeTestFlight = regexp.MustCompile(
	`(^|\n)\|\s*(.*?)\s*\|\s*\[(https?://testflight\.apple\.com/join/\w+)][^|]*\s*\|\s*([YFND])`,
)
