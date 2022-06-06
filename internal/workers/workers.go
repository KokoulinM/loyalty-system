package workers

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/KokoulinM/go-musthave-diploma-tpl/cmd/gophermart/config"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/models"
	"github.com/KokoulinM/go-musthave-diploma-tpl/internal/tasks"
	"github.com/rs/zerolog"
)

type WorkerPoolJobStore interface {
	GetJobToExecute(ctx context.Context, maxCount int) ([]models.JobStoreRow, error)
	ExecuteJob(ctx context.Context, jobID string) error
	IncreaseCounter(ctx context.Context, jobID string, count int) error
}

type Job struct {
	ID    string
	Func  func(ctx context.Context) error
	Count int
}

type WorkerPool struct {
	jobStore         WorkerPoolJobStore
	taskStore        *tasks.TaskStore
	numOfWorkers     int
	inputCh          chan Job
	logger           *zerolog.Logger
	maxJobRetryCount int
}

func New(jobStore WorkerPoolJobStore, taskStore *tasks.TaskStore, cfg *config.ConfigWorkerPool,
	logger *zerolog.Logger) *WorkerPool {
	wp := &WorkerPool{
		jobStore:         jobStore,
		taskStore:        taskStore,
		numOfWorkers:     cfg.NumOfWorkers,
		inputCh:          make(chan Job, cfg.PoolBuffer),
		logger:           logger,
		maxJobRetryCount: cfg.MaxJobRetryCount,
	}
	return wp
}

func (wp *WorkerPool) Run(ctx context.Context) {
	wg := &sync.WaitGroup{}
	for i := 0; i < wp.numOfWorkers; i++ {
		wg.Add(1)
		go func(i int) {
			wp.logger.Log().Msgf("Worker #%v start ", i)
		outer:
			for {
				select {
				case job := <-wp.inputCh:
					err := job.Func(ctx)
					if err != nil {
						wp.logger.Error().Msgf("Error on worker #%v: %v\n", i, err.Error())
						err = wp.jobStore.IncreaseCounter(ctx, job.ID, job.Count)
						if err != nil {
							wp.logger.Error().Msgf("Error with increase job counter with job:%v error:%v", job.ID, err.Error())
						}
						continue
					}
					err = wp.jobStore.ExecuteJob(ctx, job.ID)
					if err != nil {
						wp.logger.Error().Msgf("Error with executing job: %v, error: %v", job.ID, err.Error())
					}
				case <-ctx.Done():
					break outer
				}

			}
			wp.logger.Log().Msgf("Worker #%v close\n", i)
			wg.Done()
		}(i)
	}
	sch := wp.scheduler(ctx)
	wg.Wait()
	close(wp.inputCh)
	sch.Stop()
}

func (wp *WorkerPool) push(task Job) {
	wp.inputCh <- task
}

func (wp *WorkerPool) scheduler(ctx context.Context) *time.Ticker {
	ticker := time.NewTicker(time.Second * 5)
	wp.logger.Log().Msg("start scheduler")
	go func() {
		for {
			select {
			case <-ticker.C:
				wp.logger.Log().Msg("ticker tick")
				wp.transferTaskToWorkerPool(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
	return ticker
}

func (wp *WorkerPool) transferTaskToWorkerPool(ctx context.Context) {
	jobs, err := wp.jobStore.GetJobToExecute(ctx, wp.maxJobRetryCount)
	if err != nil {
		wp.logger.Error().Msgf("Error occured in getting task in worker pool: %v", err.Error())
		return
	}
	for _, job := range jobs {

		task, ok := wp.taskStore.MapOfTask[job.Type]

		if !ok {
			wp.logger.Error().Msgf("Get job of unknown type: %v", job.Type)
			continue
		}
		parameters := make(map[string]string)
		err := json.Unmarshal([]byte(job.Parameters), &parameters)
		if err != nil {
			wp.logger.Error().Msgf("Error with parce parameters, job_id: %v, err: %v", job.ID, err.Error())
			continue
		}
		function, err := task.CreateFunction(parameters)
		if err != nil {
			wp.logger.Error().Msgf("Wrong paramenters for function, job_id: %v, err: %v", job.ID, err.Error())
			continue
		}
		jobToPush := Job{
			ID:   job.ID,
			Func: function,
		}

		wp.push(jobToPush)
	}
}
