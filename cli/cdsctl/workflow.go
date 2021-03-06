package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/fsamin/go-dump"
	"github.com/spf13/cobra"

	"github.com/ovh/cds/cli"
	"github.com/ovh/cds/sdk"
)

var (
	workflowCmd = cli.Command{
		Name:  "workflow",
		Short: "Manage CDS workflow",
	}

	workflow = cli.NewCommand(workflowCmd, nil,
		[]*cobra.Command{
			cli.NewListCommand(workflowListCmd, workflowListRun, nil),
			cli.NewListCommand(workflowHistoryCmd, workflowHistoryRun, nil),
			cli.NewGetCommand(workflowShowCmd, workflowShowRun, nil),
			cli.NewDeleteCommand(workflowDeleteCmd, workflowDeleteRun, nil),
			cli.NewCommand(workflowRunManualCmd, workflowRunManualRun, nil),
			cli.NewCommand(workflowStopCmd, workflowStopRun, nil),
			cli.NewCommand(workflowExportCmd, workflowExportRun, nil),
			cli.NewCommand(workflowImportCmd, workflowImportRun, nil),
			cli.NewCommand(workflowPullCmd, workflowPullRun, nil),
			cli.NewCommand(workflowPushCmd, workflowPushRun, nil),
			workflowArtifact,
		})
)

var workflowListCmd = cli.Command{
	Name:  "list",
	Short: "List CDS workflows",
	Args: []cli.Arg{
		{Name: "project-key"},
	},
}

func workflowListRun(v cli.Values) (cli.ListResult, error) {
	w, err := client.WorkflowList(v["project-key"])
	if err != nil {
		return nil, err
	}
	return cli.AsListResult(w), nil
}

var workflowHistoryCmd = cli.Command{
	Name:  "history",
	Short: "History of a CDS workflow",
	Args: []cli.Arg{
		{Name: "project-key"},
		{Name: "workflow-name"},
	},
	OptionalArgs: []cli.Arg{
		{
			Name: "offset",
			IsValid: func(s string) bool {
				match, _ := regexp.MatchString(`[0-9]?`, s)
				return match
			},
			Weight: 1,
		},
		{
			Name: "limit",
			IsValid: func(s string) bool {
				match, _ := regexp.MatchString(`[0-9]?`, s)
				return match
			},
			Weight: 2,
		},
	},
}

func workflowHistoryRun(v cli.Values) (cli.ListResult, error) {
	var offset int64
	if v.GetString("offset") != "" {
		var errn error
		offset, errn = v.GetInt64("offset")
		if errn != nil {
			return nil, errn
		}
	}

	var limit int64
	if v.GetString("limit") != "" {
		var errl error
		limit, errl = v.GetInt64("limit")
		if errl != nil {
			return nil, errl
		}
	}

	w, err := client.WorkflowRunList(v["project-key"], v["workflow-name"], offset, limit)
	if err != nil {
		return nil, err
	}
	return cli.AsListResult(w), nil
}

var workflowShowCmd = cli.Command{
	Name:  "show",
	Short: "Show a CDS workflow",
	Args: []cli.Arg{
		{Name: "project-key"},
		{Name: "workflow-name"},
	},
	OptionalArgs: []cli.Arg{
		{
			Name: "run-number",
			IsValid: func(s string) bool {
				match, _ := regexp.MatchString(`[0-9]?`, s)
				return match
			},
			Weight: 1,
		},
	},
}

func workflowShowRun(v cli.Values) (interface{}, error) {
	var runNumber int64
	if v.GetString("run-number") != "" {
		var errl error
		runNumber, errl = v.GetInt64("run-number")
		if errl != nil {
			return nil, errl
		}
	}

	if runNumber == 0 {
		w, err := client.WorkflowGet(v["project-key"], v["workflow-name"])
		if err != nil {
			return nil, err
		}
		return *w, nil
	}

	w, err := client.WorkflowRunGet(v["project-key"], v["workflow-name"], runNumber)
	if err != nil {
		return nil, err
	}

	var tags []string
	for _, tag := range w.Tags {
		tags = append(tags, fmt.Sprintf("%s:%s", tag.Tag, tag.Value))
	}

	type wtags struct {
		sdk.WorkflowRun
		Payload string `cli:"payload"`
		Tags    string `cli:"tags"`
	}

	var payload []string
	if v, ok := w.WorkflowNodeRuns[w.Workflow.RootID]; ok {
		if len(v) > 0 {
			e := dump.NewDefaultEncoder(new(bytes.Buffer))
			e.Formatters = []dump.KeyFormatterFunc{dump.WithDefaultLowerCaseFormatter()}
			e.ExtraFields.DetailedMap = false
			e.ExtraFields.DetailedStruct = false
			e.ExtraFields.Len = false
			e.ExtraFields.Type = false
			pl, errm1 := e.ToStringMap(v[0].Payload)
			if errm1 != nil {
				return nil, errm1
			}
			for k, kv := range pl {
				payload = append(payload, fmt.Sprintf("%s:%s", k, kv))
			}
			payload = append(payload)
		}
	}

	wt := &wtags{*w, strings.Join(payload, " "), strings.Join(tags, " ")}
	return *wt, nil
}

var workflowDeleteCmd = cli.Command{
	Name:  "delete",
	Short: "Delete a CDS workflow",
	Args: []cli.Arg{
		{Name: "project-key"},
		{Name: "workflow-name"},
	},
}

func workflowDeleteRun(v cli.Values) error {
	err := client.WorkflowDelete(v["project-key"], v["workflow-name"])
	if err != nil && v.GetBool("force") && sdk.ErrorIs(err, sdk.ErrWorkflowNotFound) {
		fmt.Println(err.Error())
		os.Exit(0)
	}
	return err
}

var workflowStopCmd = cli.Command{
	Name:  "stop",
	Short: "Stop a CDS workflow or a specific node name",
	Long:  "Stop a CDS workflow or a specific node name",
	Example: `
		cdsctl workflow stop MYPROJECT myworkflow 5 # To stop a workflow run on number 5
		cdsctl workflow stop MYPROJECT myworkflow 5 compile # To stop a workflow node run on workflow run 5
	`,
	Args: []cli.Arg{
		{Name: "project-key"},
		{Name: "workflow-name"},
		{Name: "run-number"},
	},
	OptionalArgs: []cli.Arg{
		{Name: "node-name"},
	},
}

func workflowStopRun(v cli.Values) error {
	var fromNodeID int64
	runNumber, errp := strconv.ParseInt(v.GetString("run-number"), 10, 64)
	if errp != nil {
		return fmt.Errorf("run-number invalid: not a integer")
	}

	if v.GetString("node-name") != "" {
		if runNumber <= 0 {
			return fmt.Errorf("You can use flag node-name without flag run-number")
		}
		wr, err := client.WorkflowRunGet(v["project-key"], v["workflow-name"], runNumber)
		if err != nil {
			return err
		}
		for _, wnrs := range wr.WorkflowNodeRuns {
			if wnrs[0].WorkflowNodeName == v.GetString("node-name") {
				fromNodeID = wnrs[0].ID
				break
			}
		}
		if fromNodeID == 0 {
			return fmt.Errorf("Node not found")
		}
	}

	if fromNodeID != 0 {
		wNodeRun, err := client.WorkflowNodeStop(v["project-key"], v["workflow-name"], runNumber, fromNodeID)
		if err != nil {
			return err
		}
		fmt.Printf("Workflow node %s from workflow %s #%d has been stopped\n", v.GetString("node-name"), v["workflow-name"], wNodeRun.Number)
	} else {
		w, err := client.WorkflowStop(v["project-key"], v["workflow-name"], runNumber)
		if err != nil {
			return err
		}
		fmt.Printf("Workflow %s #%d has been stopped\n", v["workflow-name"], w.Number)
	}

	return nil
}
