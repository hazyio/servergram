package wrapper

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"hazyio.com/servergram/agent/internal/db"
)

func RunJob(database *gorm.DB, jobID uint64, command string, trigger db.CronRunTrigger) error {
	job, err := db.FindCronJobByID(database, jobID)
	if err != nil {
		// TODO: implement a better way to handle this error
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var newRun *db.CronRun
	if job.MetaKeepCount > 0 {
		newRun = &db.CronRun{
			CronJobID: int64(jobID),
			Status:    db.RunStatusRunning,
			Trigger:   trigger,
			StartedAt: time.Now(),
		}
		if err := database.Create(newRun).Error; err != nil {
			return err
		}
	}
	go func() {
		// remove old cron runs
		if err := db.PruneCronRuns(database, job.MetaKeepCount, job.ID); err != nil {
			slog.Error("failed to prune cron runs", "error", err)
		}
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

	result, err := Run(ctx, command)
	job.LastRunAt = time.Now()
	if err != nil {
		job.Status = db.CronJobStatusFailed

	} else {
		job.Status = db.CronJobStatusSuccess

	}
	if err := database.Save(job).Error; err != nil {
		slog.Error("failed to save cron job", "error", err)
	}

	if job.MetaKeepCount > 0 && newRun != nil {
		if err != nil {
			newRun.Status = db.RunStatusFailed
		} else {
			newRun.Status = db.RunStatusSuccess
		}
		newRun.FinishedAt = time.Now()
		newRun.ExitCode = result.ExitCode
		if err := database.Save(newRun).Error; err != nil {
			return err
		}

	}

	return nil
}
