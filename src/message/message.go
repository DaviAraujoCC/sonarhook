package message

import (
	"bytes"
	"fmt"
	"net/http"
	"sonarhook/src/config"
	"sonarhook/util"
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

type MessageConstructor interface {
	ParseMessage() (string, error)
	SendMessage(text string) error
}

type messageConstructor struct {
	message Message
}

func NewMessage(message Message) MessageConstructor {
	return &messageConstructor{message}
}

func (mc *messageConstructor) ParseMessage() (string, error) {

	if mc.message.AnalysedAt == "" {
		return "", fmt.Errorf("Incorrect Format:")
	}

	if config.Status != "" && mc.message.QualityGate.Status != config.Status {
		return "", fmt.Errorf("Ignoring status: %s", mc.message.QualityGate.Status)
	}

	var bodyMessage strings.Builder

	bodyMessage.WriteString("*SonarQube Quality Gate*\\n")

	bodyMessage.WriteString(fmt.Sprintf("Analysed at: %s\\n\\n", util.ParseTime(mc.message.AnalysedAt)))

	switch mc.message.QualityGate.Status {
	case "OK":
		bodyMessage.WriteString("*Status*: PASS \xE2\x9C\x85\\n\\n")
	case "ERROR":
		bodyMessage.WriteString("*Status*: FAILED \xF0\x9F\x9A\xAB\\n\\n")
	}

	bodyMessage.WriteString(fmt.Sprintf("*Project:* " + mc.message.Project.Name + "\\n"))

	switch mc.message.Branch.Type {
	case "BRANCH":
		bodyMessage.WriteString(fmt.Sprintf("*Branch:* " + mc.message.Branch.Name + "\\n"))
	case "PULL_REQUEST":
		bodyMessage.WriteString(fmt.Sprintf("*Pull request*: ID %s\\n", mc.message.Branch.Name))
	}

	bodyMessage.WriteString(fmt.Sprintf("*Results:* " + mc.message.Branch.URL + "\\n"))

	return bodyMessage.String(), nil
}

func (mc *messageConstructor) SendMessage(text string) error {

	client := &http.Client{}
	client.Timeout = 10 * time.Second

	json := []byte(`{"text": "` + text + `"}`)

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
