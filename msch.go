// The MIT License (MIT)
//
// Copyright (c) 2016 Fredy Wijaya
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const (
	mavenURL string = "http://search.maven.org/solrsearch/select?q="
)

var (
	keyword string
	gradle  bool
	maven   bool
)

type searchResult struct {
	Response struct {
		Docs []struct {
			Group    string `json:"g"`
			Artifact string `json:"a"`
			Version  string `json:"latestVersion"`
		} `json:"docs"`
	} `json:"response"`
}

func search(keyword string) error {
	res, err := http.Get(mavenURL + url.QueryEscape(keyword))
	if err != nil {
		return err
	}
	defer res.Body.Close()
	decoder := json.NewDecoder(res.Body)
	var searchResult searchResult
	err = decoder.Decode(&searchResult)
	if err != nil {
		return err
	}
	for _, doc := range searchResult.Response.Docs {
		if gradle {
			fmt.Println(fmt.Sprintf("%s:%s:%s", doc.Group, doc.Artifact, doc.Version))
		} else if maven {
			fmt.Println(fmt.Sprintf(`<dependency>
    <groupId>%s</groupId>
    <artifactId>%s</artifactId>
    <version>%s</version>
</dependency>`, doc.Group, doc.Artifact, doc.Version))
			fmt.Println()
		}
	}
	return nil
}

func init() {
	flag.StringVar(&keyword, "keyword", "", "Search keyword")
	flag.BoolVar(&gradle, "gradle", false, "Gradle format")
	flag.BoolVar(&maven, "maven", false, "Maven format")
}

func validateArgs() {
	if len(keyword) == 0 {
		errorAndExit("--keyword option is required")
	}
	if !gradle && !maven {
		errorAndExit("Either --gradle or maven option is required")
	}
	if gradle && maven {
		errorAndExit("--gradle and --maven options are mutually exclusive")
	}
}

func errorAndExit(msg interface{}) {
	fmt.Println("Error:", msg)
	os.Exit(1)
}

func main() {
	flag.Parse()
	validateArgs()
	err := search(keyword)
	if err != nil {
		errorAndExit(err)
	}
}
