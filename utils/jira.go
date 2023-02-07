package utils

import "github.com/andygrunwald/go-jira"

func NewJiraClient(config *Config) (*jira.Client, error) {
	tp := jira.PATAuthTransport{
		Token: config.JiraToken,
	}
	jiraClient, err := jira.NewClient(tp.Client(), config.JiraHost)
	if err != nil {
		return nil, err
	}
	return jiraClient, nil
}
