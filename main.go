package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gopkg.in/yaml.v2"
)

type dynatrace struct {
	url      string
	token    string
	problems struct {
		Result struct {
			TotalOpenProblemsCount int `json:"totalOpenProblemsCount"`
			OpenProblemCounts      struct {
				Inf         int `json:"INFRASTRUCTURE"`
				Service     int `json:"SERVICE"`
				Application int `json:"APPLICATION"`
				Environment int `json:"ENVIRONMENT"`
			} `json:"openProblemCounts"`
		} `json:"result"`
	}
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

	err = json.Unmarshal(output, &d.problems)
	if err != nil {
		panic(err)
	}
}

func main() {

	d := dynatrace{}
	d.getConfig("config.yml")
	d.getProblems()

	fmt.Printf("\nDynatrace Issues:\n\tINFRASTRUCTURE: %v\n\tSERVICE: %v\n\tAPPLICATION: %v\n\tENVIRONMENT: %v\n", d.problems.Result.OpenProblemCounts.Inf, d.problems.Result.OpenProblemCounts.Service, d.problems.Result.OpenProblemCounts.Application, d.problems.Result.OpenProblemCounts.Environment)
}
