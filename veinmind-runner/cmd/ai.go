package main

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"os"

	"github.com/chaitin/libveinmind/go/cmd"
	"github.com/chaitin/veinmind-common-go/service/report/event"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/log"
	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/sdk"
)

var aiCmd = &cmd.Command{
	Use:   "analyze",
	Short: "Analyze Veinmind Report With OpenAI.",
	RunE:  AnalyzeAI,
}

func AnalyzeAI(cmd *cmd.Command, args []string) error {
	log.GetModule(log.AIAnalyzerKey).Info("Initializing AI environment......")
	token, err := cmd.Flags().GetString("openai-token")
	if err != nil {
		return err
	}
	if token == "" {
		return errors.New("empty openai token")
	}
	r, err := cmd.Flags().GetString("report")
	if err != nil {
		return err
	}
	prefix, _ := cmd.Flags().GetString("prefix")
	log.GetModule(log.AIAnalyzerKey).Info("Reading report data......")
	reportFile, err := os.Open(r)
	if err != nil {
		log.GetModule(log.AIAnalyzerKey).Errorf("Open report data error: %s", err)
		return nil
	}
	data, err := io.ReadAll(reportFile)
	if err != nil {
		log.GetModule(log.AIAnalyzerKey).Errorf("Read report data error: %s", err)
		return nil
	}
	events := make([]*event.Event, 0)
	err = json.Unmarshal(data, &events)
	if err != nil {
		log.GetModule(log.AIAnalyzerKey).Errorf("Format report data error: %s", err)
		return nil
	}
	return AnalyzeReport(cmd.Context(), token, prefix, events)
}

func AnalyzeReport(ctx context.Context, token string, prefix string, events []*event.Event) error {
	total := func() int {
		cnt := 0
		for _, e := range events {
			if e.EventType == event.Info {
				continue
			}
			cnt += 1
		}
		return cnt
	}()
	log.GetModule(log.AIAnalyzerKey).Infof("Totally %d events, starting analyze......", total)
	for i, e := range events {
		if e.EventType == event.Info {
			continue
		}
		content, _ := json.Marshal(e)
		stream, err := sdk.DialogueStream(ctx, token, prefix, string(content))
		if err != nil {
			if errors.Unwrap(err) == io.EOF {
				log.GetModule(log.AIAnalyzerKey).Errorf("Network cannot be connected. Please check your network connection settings")
			} else {
				log.GetModule(log.AIAnalyzerKey).Errorf("Network error: %s", err)
			}
			return nil
		}
		log.GetModule(log.AIAnalyzerKey).Infof("[%d/%d] Analyzing report events: %s ......", i+1, total, e.ID)
		log.GetModule(log.AIAnalyzerKey).Info("")
		err = sdk.Read(stream)
		if err != nil {
			return err
		}
	}
	log.GetModule(log.AIAnalyzerKey).Info("Analyzed Over.")
	return nil
}

func init() {
	aiCmd.Flags().StringP("openai-token", "", "", "OpenAI token")
	aiCmd.Flags().StringP("report", "r", "report.json", "report (json) file path")
	aiCmd.Flags().StringP("prefix", "p", "", "training openai limit sentence")
}
