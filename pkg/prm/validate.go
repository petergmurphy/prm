//nolint:structcheck,unused
package prm

import (
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/rs/zerolog/log"
)

type ValidateExitCode int64

const (
	VALIDATION_PASS ValidateExitCode = iota
	VALIDATION_FAILED
	VALIDATION_ERROR
)

// Validate allows a lits of tool names to be executed against
// the codeDir.
//
// Tools can be empty, in which case we expect that a local
// configuration file (validate.yml) will contain a list of
// tools to run.
func (p *Prm) Validate(tool *Tool, args []string, outputSettings OutputSettings) error {

	// is the tool available?
	err := p.Backend.GetTool(tool, p.RunningConfig)
	if err != nil {
		log.Error().Msgf("Failed to validate with tool: %s/%s", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		return err
	}

	// the tool is available so execute against it
	exit, err := p.Backend.Validate(tool, args, p.RunningConfig, DirectoryPaths{codeDir: p.CodeDir, cacheDir: p.CacheDir}, outputSettings)

	switch exit {
	case VALIDATION_PASS:
		log.Info().Msgf("Tool %s/%s validated successfully", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		log.Info().Msg("PASS")
	case VALIDATION_FAILED:
		log.Error().Msgf("Tool %s/%s validation returned at least one failure", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id)
		log.Error().Msg("FAIL")
	case VALIDATION_ERROR:
		log.Error().Msgf("Tool %s/%s encountered errored during validation %s", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id, err)
		log.Error().Msg("ERROR")
	default:
		log.Info().Msgf("Tool %s/%s exited with code %d", tool.Cfg.Plugin.Author, tool.Cfg.Plugin.Id, exit)
	}

	return err
}

type result struct {
	toolName string
	err      error
	exit     ValidateExitCode
}

func (p *Prm) ValidateTools(toolsInfo []ToolInfo, workerCount int) error {
	noOfTools := len(toolsInfo)
	jobs := make(chan ToolInfo, noOfTools)
	results := make(chan result, noOfTools)

	for w := 1; w <= workerCount; w++ {
		go p.worker(w, jobs, results)
	}

	for _, tool := range toolsInfo {
		jobs <- tool
	}
	close(jobs)

	var resultsList []result
	for a := 1; a <= noOfTools; a++ {
		resultsList = append(resultsList, <-results)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tool Name", "Validation Exit Code", "Error"})
	table.SetBorder(false)
	for _, result := range resultsList {
		errString := "No Error"
		if result.err != nil {
			errString = result.err.Error()
		}
		table.Append([]string{result.toolName, fmt.Sprint(result.exit), errString})
	}
	table.Render()
	return nil // TODO: CHECK ERROR
}

func (p *Prm) worker(id int, jobs <-chan ToolInfo, output chan<- result) {
	for job := range jobs {
		toolName := job.Tool.Cfg.Plugin.Id
		log.Info().Msgf("Worker %d: Validating with the %s tool", id, toolName)
		err := p.Backend.GetTool(job.Tool, p.RunningConfig)
		if err != nil {
			log.Error().Msgf("Failed to validate with tool: %s/%s", job.Tool.Cfg.Plugin.Author, job.Tool.Cfg.Plugin.Id)
			output <- result{toolName: toolName, err: err, exit: VALIDATION_ERROR}
			continue
		}

		// the tool is available so execute against it
		exit, err := p.Backend.Validate(job.Tool, job.Args, p.RunningConfig, DirectoryPaths{codeDir: p.CodeDir, cacheDir: p.CacheDir}, job.OutputSettings)
		if err != nil {
			output <- result{toolName: toolName, err: err, exit: exit}
		} else {
			output <- result{toolName: toolName, err: err, exit: exit}
		}
	}
}
