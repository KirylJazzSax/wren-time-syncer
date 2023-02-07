package fetcher

import "time"

type Task interface {
	GetKey() string
	GetFrom() time.Time
	GetTo() time.Time
	GetComment() string
	GetCategoryName() string
	GetMinutes() uint64
	Synced() bool
	SetSynced(f bool)
}

type IssuesFetcher interface {
	FetchIssues(tasks []Task, d time.Time) []Issue
}

type Issue interface {
	GetKey() string
	IsPossiblySynced() bool
	GetFetchError() error
	GetSyncError() error
	SetFetchError(e error)
	SetSyncError(e error)
}

type TaskFetcher interface {
	FetchTasks(d time.Time) ([]Task, error)
}
