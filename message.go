package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

var url = os.Getenv("GOOGLE_CHAT_WEBHOOK_URL")

func init() {
	if url == "" {
		fmt.Println("GOOGLE_CHAT_WEBHOOK_URL is not set")
		os.Exit(1)
	}
}

type Message struct {
	AnalysedAt string `json:"analysedAt"`
	Branch     struct {
		IsMain bool   `json:"isMain"`
		Name   string `json:"name"`
		Type   string `json:"type"`
		URL    string `json:"url"`
	} `json:"branch"`
	ChangedAt string `json:"changedAt"`
	Project   struct {
		Key  string `json:"key"`
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"project"`
	Properties struct {
		SonarAnalysisDetectedci  string `json:"sonar.analysis.detectedci"`
		SonarAnalysisDetectedscm string `json:"sonar.analysis.detectedscm"`
	} `json:"properties"`
	QualityGate struct {
		Conditions []struct {
			ErrorThreshold string `json:"errorThreshold"`
			Metric         string `json:"metric"`
			Operator       string `json:"operator"`
			Status         string `json:"status"`
			Value          string `json:"value,omitempty"`
		} `json:"conditions"`
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"qualityGate"`
	Revision  string `json:"revision"`
	ServerURL string `json:"serverUrl"`
	Status    string `json:"status"`
	TaskID    string `json:"taskId"`
}

func (m *Message) validateMessage() error {

	if m.AnalysedAt == "" {
		return fmt.Errorf("Incorrect Format:")
	}

	return nil
}

func (m *Message) sendMessage() error {

	var bodyMessage strings.Builder

	if m.QualityGate.Status == "OK" {
		bodyMessage.WriteString("*SonarQube Quality Gate*\\n")
		bodyMessage.WriteString(fmt.Sprintf("Analysed at: %s\\n\\n", m.AnalysedAt))
		bodyMessage.WriteString(fmt.Sprintf("*Status*: %s \xE2\x9C\x85\\n\\n", m.QualityGate.Status))
		bodyMessage.WriteString(fmt.Sprintf("*Project:* " + m.Project.Name + "\\n"))
		switch m.Branch.Type {
		case "BRANCH":
			bodyMessage.WriteString(fmt.Sprintf("*Branch:* " + m.Branch.Name + "\\n"))
		case "PULL_REQUEST":
			bodyMessage.WriteString(fmt.Sprintf("*Pull request*: %s\\n", m.Branch.Name))
		}
		bodyMessage.WriteString(fmt.Sprintf("*Results:* " + m.Branch.URL + "\\n"))

	} else if m.QualityGate.Status == "ERROR" {

		bodyMessage.WriteString("*SonarQube Quality Gate*\\n")
		bodyMessage.WriteString(fmt.Sprintf("Analysed at: %s\\n\\n", m.AnalysedAt))
		bodyMessage.WriteString(fmt.Sprintf("*Status*: %s \xF0\x9F\x9A\xAB\\n\\n", m.QualityGate.Status))
		bodyMessage.WriteString(fmt.Sprintf("*Project:* " + m.Project.Name + "\\n"))
		switch m.Branch.Type {
		case "BRANCH":
			bodyMessage.WriteString(fmt.Sprintf("*Branch:* " + m.Branch.Name + "\\n"))
		case "PULL_REQUEST":
			bodyMessage.WriteString(fmt.Sprintf("*Pull request*: %s\\n", m.Branch.Name))
		}
		bodyMessage.WriteString(fmt.Sprintf("*Results:* " + m.Branch.URL + "\\n"))

	}

	client := &http.Client{}
	client.Timeout = 10 * time.Second

	json := []byte(`{"text": "` + bodyMessage.String() + `"}`)
	fmt.Println(string(json))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	fmt.Println(resp.Status)

	res, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(res))

	defer resp.Body.Close()

	return nil
}
