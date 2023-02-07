package fetcher

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"wren-time-syncer/utils"

	"github.com/andygrunwald/go-jira"
	"github.com/jedib0t/go-pretty/progress"
)

type JiraIssue struct {
	*jira.Issue
	PossiblySynced bool
	ErrorFetching  error
	ErrorSync      error
}

func (issue *JiraIssue) GetKey() string {
	return issue.Key
}

func (issue *JiraIssue) IsPossiblySynced() bool {
	return issue.PossiblySynced
}

func (issue *JiraIssue) GetFetchError() error {
	return issue.ErrorFetching
}

func (issue *JiraIssue) GetSyncError() error {
	return issue.ErrorSync
}

func (issue *JiraIssue) SetFetchError(e error) {
	issue.ErrorFetching = e
}
func (issue *JiraIssue) SetSyncError(e error) {
	issue.ErrorSync = e
}

type JiraIssuesFetcher struct {
	*jira.Client
	*sync.WaitGroup
	progress.Writer
	*utils.Config
}

func (jif *JiraIssuesFetcher) FetchIssues(tasks []Task, d time.Time) []Issue {
	jif.WaitGroup.Add(len(tasks))
	fetchedIssues := make([]JiraIssue, len(tasks))

	go jif.Writer.Render()

	tracker := utils.NewTracker("fetching tasks", int64(len(tasks)))
	jif.Writer.AppendTracker(&tracker)
	unixStr := strconv.FormatInt(d.UnixMilli(), 10)
	for i, t := range tasks {
		go func(task Task, idx int) {
			jqlStr := fmt.Sprintf(
				"worklogDate = %s AND worklogAuthor = %s AND issuekey = %s AND worklogComment ~ %s",
				unixStr,
				jif.Config.Username,
				task.GetKey(),
				fmt.Sprintf("\"%s\"", task.GetComment()),
			)

			issue, _, err := jif.Client.Issue.Get(task.GetKey(), &jira.GetQueryOptions{})
			issues := []jira.Issue{}
			if err == nil {
				issues, _, _ = jif.Client.Issue.Search(jqlStr, &jira.SearchOptions{})
			}

			fetchedIssues[idx] = JiraIssue{
				issue,
				len(issues) > 0,
				err,
				nil,
			}

			tracker.Increment(1)
			if tracker.IsDone() {
				time.Sleep(time.Millisecond * 100)
				jif.Writer.Stop()
				fmt.Println()
			}
			jif.WaitGroup.Done()
		}(t, i)
	}
	jif.WaitGroup.Wait()

	issues := make([]Issue, 0)
	for _, fi := range fetchedIssues {
		f := fi
		issues = append(issues, &f)
	}
	return issues
}

func NewJiraIssueFetcher(jiraClient *jira.Client, pw progress.Writer, config *utils.Config) (IssuesFetcher, error) {
	return &JiraIssuesFetcher{
		jiraClient,
		&sync.WaitGroup{},
		pw,
		config,
	}, nil
}
