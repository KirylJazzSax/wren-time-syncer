package syncer

import (
	"fmt"
	"strconv"
	"sync"
	"time"
	"wren-time-syncer/fetcher"
	"wren-time-syncer/utils"

	"github.com/andygrunwald/go-jira"
	"github.com/jedib0t/go-pretty/progress"
)

type JiraIssueSyncer struct {
	*sync.WaitGroup
	*jira.Client
	progress.Writer
}

func (jis *JiraIssueSyncer) SyncIssues(tasks []fetcher.Task, issues []fetcher.Issue, force bool) {
	go jis.Writer.Render()

	tracker := utils.NewTracker("sync tasks", int64(len(tasks)))
	jis.Writer.AppendTracker(&tracker)
	jis.WaitGroup.Add(len(tasks))
	for i, t := range tasks {
		go func(task fetcher.Task, idx int) {
			if (force || (!force && !issues[idx].IsPossiblySynced())) && !task.Synced() {
				jis.SyncIssue(task, issues[idx])
			}

			tracker.Increment(1)
			if tracker.IsDone() {
				time.Sleep(time.Millisecond * 100)
				jis.Writer.Stop()
				fmt.Println()
			}
			jis.WaitGroup.Done()
		}(t, i)
	}
	jis.WaitGroup.Wait()
}

func (jis *JiraIssueSyncer) SyncIssue(task fetcher.Task, issue fetcher.Issue) {
	jiraTime := jira.Time(task.GetFrom())
	_, _, err := jis.Client.Issue.AddWorklogRecord(task.GetKey(), &jira.WorklogRecord{
		TimeSpent: strconv.FormatUint(task.GetMinutes(), 10) + "m",
		Comment:   task.GetComment(),
		Started:   &jiraTime,
	})
	if err != nil {
		issue.SetSyncError(err)
		task.SetSynced(false)
		return
	}

	issue.SetSyncError(nil)
	task.SetSynced(true)
}

func NewJiraIssueSyncer(client *jira.Client, pw progress.Writer) IssuesSyncer {
	return &JiraIssueSyncer{
		&sync.WaitGroup{},
		client,
		pw,
	}
}
