package fetcher

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
	"wren-time-syncer/utils"
)

type PolcodeTaskFetcher struct {
	*http.Client
	*utils.Config
}

type PolcodeTask struct {
	Id           uint      `json:"id"`
	CategoryId   uint      `json:"categoryId"`
	CategoryName string    `json:"categoryName"`
	From         time.Time `json:"fromDate"`
	To           time.Time `json:"toDate"`
	Comment      string    `json:"comment"`
	Minutes      uint64    `json:"minutes"`
	IsSynced     bool
}

func (pc *PolcodeTask) GetFrom() time.Time {
	return pc.From
}
func (pc *PolcodeTask) GetTo() time.Time {
	return pc.To
}
func (pc *PolcodeTask) GetComment() string {
	return strings.Join(strings.Split(pc.Comment, ":")[1:], ":")
}
func (pc *PolcodeTask) GetKey() string {
	return strings.Split(pc.Comment, ":")[0]
}
func (pc *PolcodeTask) GetCategoryName() string {
	return pc.CategoryName
}
func (pc *PolcodeTask) Synced() bool {
	return pc.IsSynced
}
func (pc *PolcodeTask) SetSynced(f bool) {
	pc.IsSynced = f
}
func (pc *PolcodeTask) GetMinutes() uint64 {
	return pc.Minutes
}

func (pc *PolcodeTaskFetcher) FetchTasks(d time.Time) ([]Task, error) {
	pcTasks := make([]PolcodeTask, 0)
	req, _ := http.NewRequest("GET", fmt.Sprintf("%sapi/time/%s", pc.Config.LinkHost, d.Format("2006-01-02")), nil)
	req.Header.Set("Authorization", pc.Config.LinkAuthHeader)
	resp, err := pc.Client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(body, &pcTasks)
	if err != nil {
		return nil, err
	}

	tasks := make([]Task, 0)
	for _, t := range pcTasks {
		n := t
		tasks = append(tasks, &n)
	}

	return tasks, nil
}

func NewPolcodeTaskFetcher(config utils.Config) TaskFetcher {
	return &PolcodeTaskFetcher{
		&http.Client{},
		&config,
	}
}
