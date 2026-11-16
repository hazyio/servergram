package wrapper

import (
	"context"
	"time"

	"gorm.io/gorm"
	"hazyio.com/servergram/agent/internal/cron"
	"hazyio.com/servergram/agent/internal/db"
)

func RunJob(database *gorm.DB, jobID uint64, command string, trigger db.CronRunTrigger) error {
	job, err := db.FindCronJobByID(database, jobID)
	if err != nil {
		// TODO: implement a better way to handle this error
		return err
	}
	newRun := &db.CronRun{
		JobID:     int64(jobID),
		Status:    db.RunStatusRunning,
		Trigger:   trigger,
		StartedAt: time.Now(),
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := database.Create(newRun).Error; err != nil {
		return err
	}
	cmd, err := cron.GetCronJobCommandByHash(ctx, job.Hash)
	if err != nil {
		return err
	}

	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				statusErr := newRun.RetrieveStatus(database)
				if statusErr != nil {
					return
				}
				if newRun.Status != db.RunStatusRunning {
					// cancel the context to stop the command
					cancel()
					return
				}
			}
		}
	}()
	result, err := Run(ctx, cmd)
	if err != nil {
		newRun.UpdateStatus(database, db.RunStatusFailed)
		return err
	}
	newRun.FinishedAt = time.Now()
	newRun.ExitCode = result.ExitCode
	newRun.UpdateStatus(database, db.RunStatusSuccess)
	return nil
}
