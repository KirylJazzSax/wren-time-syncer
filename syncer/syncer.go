package syncer

import "wren-time-syncer/fetcher"

type IssuesSyncer interface {
	SyncIssues(tasks []fetcher.Task, issues []fetcher.Issue, force bool)
	SyncIssue(task fetcher.Task, issue fetcher.Issue)
}
