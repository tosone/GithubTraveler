package headerlink

import (
	"regexp"
	"strings"
)

// Link ..
type Link map[string]string

// Parse ..
func Parse(links string) (res Link) {
	if links == "" {
		return
	}
	pagingReg := regexp.MustCompile(`rel="([a-z]+)"`)
	linkReg := regexp.MustCompile(`<(.*)>`)
	res = Link{}
	for _, link := range strings.Split(strings.TrimSpace(links), ",") {
		strs := strings.Split(strings.TrimSpace(link), ";")
		if len(strs) != 2 {
			continue
		}
		pagingResult := pagingReg.FindStringSubmatch(strs[1])
		linkResult := linkReg.FindStringSubmatch(strs[0])
		if len(pagingResult) == 2 && len(linkResult) == 2 {
			res[pagingResult[1]] = linkResult[1]
		}
	}
	return
}
