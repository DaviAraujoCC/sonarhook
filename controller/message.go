package controller

import (
	"bytes"
	"fmt"
	"net/http"
	"sonarhook/config"
	"strings"
	"time"
)

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

	if config.Status != "" && m.Status != config.Status {
		return fmt.Errorf("Ignoring status: %s", m.Status)
	}

	return nil
}

func (m *Message) sendMessage() error {

	var bodyMessage strings.Builder

	bodyMessage.WriteString("*SonarQube Quality Gate*\\n")

	t, _ := time.Parse(time.RFC3339, m.AnalysedAt)

	bodyMessage.WriteString(fmt.Sprintf("Analysed at: %s\\n\\n", t.Format("2006-01-02 15:04:05")))

	switch m.QualityGate.Status {
	case "OK":
		bodyMessage.WriteString("*Status*: PASS \xE2\x9C\x85\\n\\n")
	case "ERROR":
		bodyMessage.WriteString("*Status*: FAILED \xF0\x9F\x9A\xAB\\n\\n")
	}

	bodyMessage.WriteString(fmt.Sprintf("*Project:* " + m.Project.Name + "\\n"))

	switch m.Branch.Type {
	case "BRANCH":
		bodyMessage.WriteString(fmt.Sprintf("*Branch:* " + m.Branch.Name + "\\n"))
	case "PULL_REQUEST":
		bodyMessage.WriteString(fmt.Sprintf("*Pull request*: ID %s\\n", m.Branch.Name))
	}

	bodyMessage.WriteString(fmt.Sprintf("*Results:* " + m.Branch.URL + "\\n"))

	client := &http.Client{}
	client.Timeout = 10 * time.Second

	json := []byte(`{"text": "` + bodyMessage.String() + `"}`)
	fmt.Println(string(json))

	req, err := http.NewRequest("POST", config.GoogleChatWebhookURL, bytes.NewBuffer(json))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}
