package extract

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

var caseDir = "../../test/extractor"

type Case struct {
	url      string
	html     string
	template *Template
	result   string
}

func TestExtractor(t *testing.T) {
	// load case
	var cases map[string]*Case = map[string]*Case{}
	if files, err := ioutil.ReadDir(caseDir); err == nil {
		for _, file := range files {
			fullPath := caseDir + "/" + file.Name()
			namePair := strings.Split(file.Name(), ".")
			if _, exist := cases[namePair[0]]; !exist {
				cases[namePair[0]] = &Case{
					template: &Template{},
				}
			}
			c, _ := cases[namePair[0]]
			switch namePair[1] {
			case "html":
				content, _ := ioutil.ReadFile(fullPath)
				c.html = string(content)
			case "json":
				content, _ := ioutil.ReadFile(fullPath)
				if err := json.Unmarshal(content, c.template); err != nil {
					t.Log(err.Error())
				}
			case "txt":
				f, _ := os.Open(fullPath)
				reader := bufio.NewReader(f)
				for {
					if line, _, err := reader.ReadLine(); err == nil {
						if strings.HasPrefix(string(line), "# url") {
							url, _, _ := reader.ReadLine()
							c.url = strings.TrimSpace(string(url))
						} else if strings.HasPrefix(string(line), "# answer") {
							result, _, _ := reader.ReadLine()
							c.result = string(result)
						}
					} else {
						break
					}
				}
			}

		}
	}
	for k, v := range cases {
		resJson := Extract(v.url, v.html, v.template).([]Map)
		resByte, _ := json.Marshal(resJson[0])
		resStr := string(resByte)

		var answerJson *Map = &Map{}
		json.Unmarshal([]byte(v.result), answerJson)
		answerByte, _ := json.Marshal(answerJson)
		answerStr := string(answerByte)

		if answerStr != resStr {
			t.Errorf("%s diff:\nt: %s\nf: %s", k, answerStr, resStr)
		}
	}
}
