package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"net/http"
)

type dynatrace struct {
	url      string
	token    string
	problems map[string]int
}

func (d *dynatrace) getConfig(c string) {

	type config struct {
		Url   string `yaml:"dynatraceURL"`
		Token string `yaml:"dynatraceToken"`
	}

	var conf config

	source, err := ioutil.ReadFile(c)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(source, &conf)
	if err != nil {
		panic(err)
	}

	d.url = conf.Url
	d.token = conf.Token

}

func (d *dynatrace) getProblems() {

	type problemCounts struct {
		Inf         int `json:"INFRASTRUCTURE"`
		Service     int `json:"SERVICE"`
		Application int `json:"APPLICATION"`
		Environment int `json:"ENVIRONMENT"`
	}

	type results struct {
		TotalOpenProblemsCount int           `json:"totalOpenProblemsCount"`
		OpenProblemCounts      problemCounts `json:"openProblemCounts"`
	}

	type result struct {
		Result results `json:"result"`
	}

	path := "/api/v1/problem/status"
	api := fmt.Sprintf("%s%s", d.url, path)

	client := &http.Client{}
	request, err := http.NewRequest("GET", api, nil)
	if err != nil {
		panic(err)
	}

	// Set Basic Auth header
	authToken := fmt.Sprintf("Api-Token %s", d.token)
	request.Header.Set("Authorization", authToken)

	res, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	output, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		panic(err)
	}

	var dynatraceProblems result
	err = json.Unmarshal(output, &dynatraceProblems)
	d.problems["Infrastructure"] = dynatraceProblems.Result.OpenProblemCounts.Inf
	d.problems["Service"] = dynatraceProblems.Result.OpenProblemCounts.Service
	d.problems["Application"] = dynatraceProblems.Result.OpenProblemCounts.Application
	d.problems["Environment"] = dynatraceProblems.Result.OpenProblemCounts.Environment
}

func main() {

	d := dynatrace{
		problems: make(map[string]int),
	}
	d.getConfig("config.yml")
	d.getProblems()
	fmt.Printf("%v\n", d.problems)

}
