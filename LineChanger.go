// package main read config json and change file
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// LineConfig root of all config
type LineConfig []LineConfigElement

// LineConfigElement each component
type LineConfigElement struct {
	Path        string   `json:"path"`
	Description string   `json:"description"`
	Encoding    string   `json:"encoding"`
	Change      []Change `json:"change"`
}

// Change find and change to
type Change struct {
	Find     string            `json:"find"`
	Changeto map[string]string `json:"changeto"`
}

func processLine(lc LineConfigElement) {
	//fmt.Printf("line : %+v\n", lc)

	for _, cc := range lc.Change {
		for k, v := range cc.Changeto {
			fmt.Println(k + " - " + v)
		}

	}
}

func main() {
	lineJSON, err := ioutil.ReadFile("example.json")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	var lines LineConfig
	err = json.Unmarshal([]byte(lineJSON), &lines)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	for _, line := range lines {
		processLine(line)
	}
	//fmt.Printf("lines : %+v", lines)
	//Birds : [{Species:pigeon Description:} {Species:eagle Description:bird of prey}]
}
