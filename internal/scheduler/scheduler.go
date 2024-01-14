// The scheduler package is a light wrapper around the gochron package
// that provides a simple interface for scheduling jobs.
//
// This will help keep our cmd package clean and simple.
package scheduler

import (
	"errors"
	"log/slog"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/google/uuid"
	"github.com/masoncitemple4/certwatch/internal/defaults"
)

var (
	HostsRequiredError = errors.New("must specify WithHosts option when calling New in order to register a new scheduler")
)

type Host struct {
	Hostname string `json:"hostname"`
	Port     string `json:"port"`
}

type Scheduler struct {
	logger *slog.Logger

	Hosts []Host

	srv gocron.Scheduler

	// convenience property
	jMap map[uuid.UUID]gocron.Job
}

type OptionFn func(*Scheduler)

func WithLogger(logger *slog.Logger) OptionFn {
	return func(s *Scheduler) {
		s.logger = logger
	}
}

func WithHosts(hosts []Host) OptionFn {
	return func(s *Scheduler) {
		s.Hosts = hosts
	}
}

func New(opts ...OptionFn) (*Scheduler, error) {
	s := &Scheduler{}

	if len(opts) < 1 {
		opts = defaultOpts()
	}

	for _, opt := range opts {
		opt(s)
	}

	if len(s.Hosts) == 0 {
		return nil, HostsRequiredError
	}

	if s.logger == nil {
		defaultOpts()[0](s)
	}

	tmp, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	s.srv = tmp

	var jobIds []uuid.UUID

	// TODO: Add configurable duration.
	job, err := s.srv.NewJob(
		gocron.DurationJob(
			5*time.Minute,
		),
		gocron.NewTask(
			checkCertsJob,
			s,
			s.Hosts,
		),
		gocron.WithStartAt(gocron.WithStartImmediately()),
		gocron.WithEventListeners(
			gocron.AfterJobRuns(func(jobID uuid.UUID, jobName string) {
				s.logger.Info("job success", slog.String("job_id", jobID.String()), slog.String("job_name", jobName))
			}),
			gocron.AfterJobRunsWithError(func(jobID uuid.UUID, jobName string, err error) {
				s.logger.Error("job error", slog.String("job_id", jobID.String()), slog.String("job_name", jobName), slog.String("error", err.Error()))
			}),
		),
	)
	if err != nil {
		return nil, err
	}

	jobIds = append(jobIds, job.ID())

	return s, nil
}

type CloseFn func()

func (s *Scheduler) Run() (CloseFn, error) {
	s.logger.Info("starting scheduler")

	// start the scheduler
	s.srv.Start()

	// when you're done, shut it down
	closeFn := func() {
		if err := s.srv.Shutdown(); err != nil {
			s.logger.Error("service error", slog.String("desc", "could error stopping service"), slog.String("error", err.Error()))
		}
		s.logger.Info("service stopped")
	}

	return closeFn, nil
}

// NOTE: purposefully omitted hosts
// to avoid committing sensitive data to the repo.
func defaultOpts() []OptionFn {

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{AddSource: true})
	logger := slog.New(handler).With(slog.String("machine_id", defaults.MachineName()))

	loggerFn := WithLogger(logger)

	return []OptionFn{loggerFn}
}

func (s *Scheduler) RunJob(jid uuid.UUID) error {
	if err := s.jMap[jid].RunNow(); err != nil {
		return err
	}

	return nil
}
