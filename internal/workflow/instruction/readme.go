package instruction

import (
	"fmt"
	"io/ioutil"
	"strings"
	"sync"
)

var once sync.Once
var README_SECTIONS []string

func loadReadMe() []string {
	if len(README_SECTIONS) <= 0 {
		once.Do(func() {
			source, _ := ioutil.ReadFile("./README.md")
			readme := string(source)
			secs := strings.Split(readme, "## ")
			README_SECTIONS = secs[1:]
			for _, sec := range secs[1:] {
				README_SECTIONS = append(README_SECTIONS, fmt.Sprintf("## %s", sec))
			}
		})
	}
	return README_SECTIONS
}
