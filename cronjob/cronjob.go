package cronjob

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/config"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/domain"
	"github.com/robfig/cron/v3"
)

var (
	cr *cron.Cron
)

type CronJob struct {
	*cron.Cron
}

func Init() (err error) {
	cr = cron.New()
	cr.Start()
	return nil
}

func GetCron() *CronJob {
	return &CronJob{cr}
}

func (cr *CronJob) AddJob(sched string, exe func()) (cron.EntryID, error) {
	if !config.Config.CronConfig.IsEnableCron {
		return 0, domain.ErrCronJobIsDisabled
	}

	jobId, err := cr.AddFunc(sched, exe)
	if err != nil {
		return 0, domain.ErrCronJobAddJobFailed
	}

	return jobId, nil
}

func (cr *CronJob) RemoveJob(jobId cron.EntryID) {
	cr.Remove(jobId)
}
