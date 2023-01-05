package pipeline

import (
	"bytes"
	"container/list"
	"encoding/json"
	"fmt"
	"html/template"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/mrsimonemms/conveyor-belt/pkg/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/exp/slices"
)

type Pipeline struct {
	Error    *Job                         `json:"-"`        // The error is used as a circuit-breaker
	Name     string                       `json:"name"`     // The pipeline name
	Jobs     *list.List                   `json:"-"`        // The jobs to run in stage order
	Response map[string]map[string]Result `json:"response"` // The collated responses from each call - keys are "stage" and "job"
	Stages   []string                     `json:"stages"`   // Stages and their ordering
	Time     struct {
		Start time.Time `json:"start"` // Used to calculate the execution time
		Total Duration  `json:"total"` // Used to display for the output
	} `json:"executionTime"`
}

type Stage struct {
	Stage string
	Jobs  []Job
}

type CompletionAction func(*Pipeline, zerolog.Logger)

func (p *Pipeline) addResponse(job Job, result *Result) {
	// Save the response data
	if p.Response == nil {
		p.Response = make(map[string]map[string]Result)
	}
	if _, ok := p.Response[job.Stage]; !ok {
		p.Response[job.Stage] = make(map[string]Result)
	}

	p.Response[job.Stage][job.Name] = *result
}

func (p *Pipeline) stopTimer() time.Duration {
	p.Time.Total.Duration = time.Since(p.Time.Start)

	return p.Time.Total.Duration
}

func (p *Pipeline) Run(postActions ...CompletionAction) error {
	// Record the start time
	p.Time.Start = time.Now()

	var wg sync.WaitGroup

	pipelineLog := log.Logger.With().Str("pipelineId", uuid.NewString()).Str("pipelineName", p.Name).Logger()

	pipelineLog.Info().Msg("Pipeline triggered")

	for stageList := p.Jobs.Front(); stageList != nil; stageList = stageList.Next() {
		wgDone := make(chan bool)
		errChan := make(chan error)

		stage := stageList.Value.(Stage)

		stageLog := pipelineLog.With().Str("stage", stage.Stage).Logger()

		wg.Add(len(stage.Jobs))

		for _, job := range stage.Jobs {
			jobLog := stageLog.With().Str("job", job.Name).Logger()

			go func(job Job) {
				defer wg.Done()

				jobLog.Info().Msg("Executing job")

				if result, err := job.Exec(p); err != nil {
					errChan <- err
					jobLog.Error().Err(err).Msg("Job errored")
				} else {
					jobLog.Info().Msg("Job completed")
					p.addResponse(job, result)
				}
			}(job)
		}

		go func() {
			stageLog.Debug().Msg("Waiting for stage to run")

			wg.Wait()
			close(wgDone)
		}()

		select {
		case <-wgDone:
			stageLog.Info().Msg("Stage completed successfully")
			continue
		case err := <-errChan:
			close(errChan)

			p.stopTimer()

			pipelineLog = pipelineLog.With().Dur("duration", p.Time.Total.Duration).Logger()

			pipelineLog.Error().Err(err).Msg("Pipeline errored - triggering error job")

			if err := p.TriggerError(err); err != nil {
				// The trigger error has failed - store error internally
				pipelineLog.Error().Err(err).Msg("Pipeline error job failed")
				return err
			}
			// Error has successfully triggered
			pipelineLog.Debug().Msg("Pipeline error message sent")
			return nil
		}
	}

	// Calculate the duration
	p.stopTimer()

	for _, c := range postActions {
		c(p, pipelineLog)
	}

	pipelineLog.Info().Dur("duration", p.Time.Total.Duration).Msg("Pipeline completed successfully")

	return nil
}

func (p *Pipeline) TriggerError(errorToSend error) error {
	if p.Error == nil {
		log.Error().Err(errorToSend).Msg("No error trigger configured")

		return errorToSend
	}

	if _, err := p.Error.Exec(p); err != nil {
		log.Error().Err(err).Msg("Error trigger failed")
		return err
	}

	return nil
}

// Build the pipeline from the config
func Build(cfg *config.Config) (*Pipeline, error) {
	stages := cfg.Spec.Stages

	// Get the jobs per stage
	jobs := make(map[string][]Job)
	for _, j := range cfg.Spec.Jobs {
		if !slices.Contains(stages, j.Stage) {
			return nil, fmt.Errorf("%w: %s", ErrUnknownStage, j.Stage)
		}

		jobs[j.Stage] = append(jobs[j.Stage], Job{
			PipelineJob: j,
		})
	}

	// Store the jobs in a linked list
	jobList := list.New()
	for _, s := range stages {
		jobList.PushBack(Stage{
			Stage: s,
			Jobs:  jobs[s],
		})
	}

	if jobList.Len() == 0 {
		return nil, ErrNoJobs
	}

	var errorJob *Job
	if cfg.Spec.Error != nil {
		errorJob = &Job{
			PipelineJob: *cfg.Spec.Error,
		}
	}

	return &Pipeline{
		Error:  errorJob,
		Name:   cfg.Metadata.Name,
		Jobs:   jobList,
		Stages: stages,
	}, nil
}

// Converts any template variables to the value and converts to JSON
func parseData(p *Pipeline, data *map[string]any) ([]byte, error) {
	strData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	t, err := template.New("pipeline").Parse(string(strData))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, p)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
