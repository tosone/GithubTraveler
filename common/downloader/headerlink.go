package downloader

import (
	"regexp"
	"strings"
)

type link map[string]string

// headerLink parse GitHub header's link for next page num
func headerLink(links string) (res link) {
	if links == "" {
		return
	}
	pagingReg := regexp.MustCompile(`rel="([a-z]+)"`)
	linkReg := regexp.MustCompile(`<(.*)>`)
	res = link{}
	for _, link := range strings.Split(strings.TrimSpace(links), ",") {
		strList := strings.Split(strings.TrimSpace(link), ";")
		if len(strList) != 2 {
			continue
		}
		pagingResult := pagingReg.FindStringSubmatch(strList[1])
		linkResult := linkReg.FindStringSubmatch(strList[0])
		if len(pagingResult) == 2 && len(linkResult) == 2 {
			res[pagingResult[1]] = linkResult[1]
		}
	}
	return
}
