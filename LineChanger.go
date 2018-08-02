// package main read config json and change file
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
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
	Find        string            `json:"find"`
	Description string            `json:"description"`
	Changeto    map[string]string `json:"changeto"`
	findPrefix  string
	findPostfix string
}

func readFile(path, encoding string) ([]string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	if encoding == "euckr" {
		euckrDec := korean.EUCKR.NewDecoder()
		got, _, err := transform.Bytes(euckrDec, data)
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(got), "\n")
		return lines, nil
	} else if encoding == "utf16bom" {
		utfDec := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder()
		got, _, err := transform.Bytes(utfDec, data)
		if err != nil {
			return nil, err
		}
		lines := strings.Split(string(got), "\n")
		return lines, nil
	}
	if encoding != "utf8" {
		return nil, fmt.Errorf("unknown encoding %s", encoding)
	}
	lines := strings.Split(string(data), "\n")
	return lines, nil
}

func writeFile(lines []string, path, encoding string) error {
	onestring := strings.Join(lines, "\n")
	// if string has \r normalize \r\n
	if strings.Index(onestring, "\r") != -1 {
		onestring = strings.Replace(onestring, "\r", "", -1)
		onestring = strings.Replace(onestring, "\n", "\r\n", -1)
	}

	if encoding == "euckr" {
		euckrEnc := korean.EUCKR.NewEncoder()
		got, _, err := transform.Bytes(euckrEnc, []byte(onestring))
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, []byte(got), 0644)
		if err != nil {
			return err
		}
		return nil
	}
	if encoding == "utf16bom" {
		utfEnc := unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewEncoder()
		got, _, err := transform.Bytes(utfEnc, []byte(onestring))
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(path, []byte(got), 0644)
		if err != nil {
			return err
		}
		return nil
	}
	err := ioutil.WriteFile(path, []byte(onestring), 0644)
	if err != nil {
		return err
	}
	return nil
}

func processItem(buildstage string, lc LineConfigElement, basePath string) error {
	// fmt.Println("-------------------------")
	// fmt.Println("Path : " + basePath + lc.Path)
	// fmt.Println("buildstage : " + buildstage)
	// fmt.Println(lc.Description)

	// preprocess find
	for i, cc := range lc.Change {
		arr := strings.Split(cc.Find, "[[tochange]]")
		if len(arr) != 2 {
			return fmt.Errorf("invalid config Find format : %s in %s", cc.Find, lc.Path)

		}
		lc.Change[i].findPrefix = arr[0]
		lc.Change[i].findPostfix = arr[1]
	}

	// read source file
	fileLines, err := readFile(basePath+lc.Path, lc.Encoding)
	if err != nil {
		return err
	}

	for _, cc := range lc.Change {

		// process line by line
		for i, text := range fileLines {
			// if line doesn't have findPrefix continue
			pre := strings.Index(text, cc.findPrefix)
			if pre == -1 {
				continue
			}
			// if findPostfix exists and [line after findPrefix] doesn't have findPostfix continue
			if len(cc.findPostfix) > 0 && strings.Index(text[pre:], cc.findPostfix) == -1 {
				continue
			}

			//fmt.Println("Found line : " + text)
			first := text[:pre+len(cc.findPrefix)]
			remain := text[pre+len(cc.findPrefix):]
			second := ""
			if len(cc.findPostfix) > 0 {
				post := strings.Index(remain, cc.findPostfix)
				if post != -1 {
					second = remain[post:]
				}
			}

			var replace string
			// if config has stage text, make first + val + second, else make first + second
			if val, ok := cc.Changeto[buildstage]; ok {
				replace = first + val + second
			} else {
				replace = first + second
			}
			fileLines[i] = replace
		}
	}

	// write back to source file
	//err = writeFile(fileLines, basePath+lc.Path+".conv", lc.Encoding)
	err = writeFile(fileLines, basePath+lc.Path, lc.Encoding)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("LineChanger configjson buildstage (basepath)")
		return
	}
	configjson := os.Args[1] // config path
	buildstage := os.Args[2] // build stage to apply
	basePath := "./"
	if len(os.Args) >= 4 {
		basePath = os.Args[3] // base path if specified
	}

	// read config json and unmarshall
	lineJSON, err := ioutil.ReadFile(configjson)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	var lines LineConfig
	err = json.Unmarshal([]byte(lineJSON), &lines)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	// process config
	for _, line := range lines {
		err := processItem(buildstage, line, basePath)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	}
}
