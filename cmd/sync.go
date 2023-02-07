package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"wren-time-syncer/fetcher"
	"wren-time-syncer/syncer"
	"wren-time-syncer/utils"

	"github.com/jedib0t/go-pretty/table"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var date string

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync wren tasks",
	Long:  "will display tasks what we need to sync then you will choose do you need to sync it or not",
	Run: func(cmd *cobra.Command, args []string) {
		config, err := utils.LoadConfig(".")
		if err != nil {
			log.Fatal(err)
		}

		dayToSync := time.Now()
		if date != "" {
			dayToSync, _ = time.Parse("2006-01-02", date)
		}

		jiraClient, err := utils.NewJiraClient(&config)
		if err != nil {
			log.Fatal(err)
		}

		pw := utils.NewProgress()

		jiraFetcher, err := fetcher.NewJiraIssueFetcher(jiraClient, pw, &config)
		if err != nil {
			log.Fatal(err)
		}

		pc := fetcher.NewPolcodeTaskFetcher(config)
		is := syncer.NewJiraIssueSyncer(jiraClient, pw)
		tasks, err := pc.FetchTasks(dayToSync)
		if err != nil {
			log.Fatal(err)
		}

		tasks = filterWrenTasks(tasks)

		issues := jiraFetcher.FetchIssues(tasks, dayToSync)

		prompt := promptui.Select{
			Label: "Select what to do please",
			Items: preparePromptItems(tasks, issues),
		}
		confirmPrompt := promptui.Prompt{
			Label:     "You sure",
			IsConfirm: true,
		}

		for {
			renderTable(tasks, issues, utils.NewTable())

			prompt.Items = preparePromptItems(tasks, issues)

			i, r, err := prompt.Run()
			if err != nil {
				log.Fatal(err)
			}

			if r == utils.ExitMessage {
				os.Exit(0)
			}

			if r == utils.SyncAllMessage || r == utils.SyncOnlyNotSyncedMessage {
				is.SyncIssues(tasks, issues, false)
				break
			}

			if r == utils.SyncAllForceMessage {
				s, _ := confirmPrompt.Run()
				if strings.EqualFold(s, "y") {
					is.SyncIssues(tasks, issues, true)
				}
				break
			}

			if issues[i].IsPossiblySynced() {
				s, _ := confirmPrompt.Run()
				if strings.EqualFold(s, "y") {
					is.SyncIssue(tasks[i], issues[i])
				}
				continue
			}

			is.SyncIssue(tasks[i], issues[i])
		}

		renderTable(tasks, issues, utils.NewTable())
		fmt.Println("Bye")
	},
}

func init() {
	SyncCmd.Flags().StringVar(&date, "date", "", "Day you want sync tasks in format 2006-01-02")
}

func renderTable(tasks []fetcher.Task, issues []fetcher.Issue, tbl table.Writer) {
	header := table.Row{"Key", "Comment", "Date", "From", "To", "Spent time", "Sync Status"}
	tbl.AppendHeader(header)

	rows := make([]table.Row, len(tasks))
	for i, t := range tasks {
		syncStatus := utils.SyncStatusGoodToGo

		if issues[i].IsPossiblySynced() {
			syncStatus = utils.SyncStatusPossiblySynced
		}

		if t.Synced() {
			syncStatus = utils.SyncStatusSynced
		}

		if issues[i].GetFetchError() != nil {
			syncStatus = issues[i].GetFetchError().Error()
		}

		if issues[i].GetSyncError() != nil {
			syncStatus = issues[i].GetSyncError().Error()
		}

		rows[i] = table.Row{
			t.GetKey(),
			t.GetComment(),
			t.GetFrom().Format("2006-01-02"),
			t.GetFrom().Format("15:04"),
			t.GetTo().Format("15:04"),
			fmt.Sprintf("%dm", t.GetMinutes()),
			syncStatus,
		}
	}
	tbl.AppendRows(rows)
	tbl.Render()
}

func preparePromptItems(tasks []fetcher.Task, issues []fetcher.Issue) (res []string) {
	hasMightSynced := false
	var sb strings.Builder
	for i, task := range tasks {
		syncJiraError := issues[i].GetSyncError()
		isPossiblySynced := issues[i].IsPossiblySynced()
		isSynced := task.Synced()

		sb.WriteString(fmt.Sprintf("%s ", task.GetKey()))

		if isSynced {
			sb.WriteString(fmt.Sprintf("- %s", utils.SyncStatusSynced))
		}

		if isPossiblySynced && !isSynced && syncJiraError == nil {
			hasMightSynced = true
			sb.WriteString("- task could be synced")
		}

		if syncJiraError != nil {
			sb.WriteString(fmt.Sprintf("- error: %s", syncJiraError.Error()))
		}

		res = append(res, sb.String())
		sb.Reset()
	}

	res = append(res, utils.SyncAllMessage)
	if hasMightSynced {
		res = append(res, utils.SyncAllForceMessage)
	}
	res = append(res, utils.ExitMessage)
	return
}

func filterWrenTasks(tasks []fetcher.Task) []fetcher.Task {
	res := make([]fetcher.Task, 0)
	for _, t := range tasks {
		if strings.Compare(t.GetCategoryName(), utils.WrenCategoryName) == 0 {
			res = append(res, t)
		}
	}

	return res
}
