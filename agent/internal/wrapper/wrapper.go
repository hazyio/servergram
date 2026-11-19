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
		slog.Error("failed to find cron job", "error", err, "job_id", jobID)
		// TODO: implement a better way to handle this error
		return err
	}
	if job.Status == db.CronJobStatusInactive {
		slog.Info("skipping inactive cron job", "job_id", job.ID)
		return nil
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
					slog.Error("failed to retrieve cron run status", "error", statusErr)
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
	slog.Info("starting cron job", "job_id", jobID)
	result, err := Run(ctx, command)
	job.LastRunAt = time.Now()

	if err := database.Save(job).Error; err != nil {
		slog.Error("failed to save cron job", "error", err)
	}

	if newRun != nil {
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
