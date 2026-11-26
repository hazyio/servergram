package wrapper

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm"
	"hazyio.com/servergram/agent/internal/common"
	"hazyio.com/servergram/agent/internal/db"
)

func RunJob(ctx context.Context, database *gorm.DB, jobID uint64, command string, trigger db.CronRunTrigger) error {
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
	var newRun *db.CronRun
	var saveFileLog *string
	if job.RunLogKeepCount > 0 {
		now := time.Now()
		logFileName, err := common.MakeLogFileName(jobID, now)
		if err != nil {
			return err
		}
		saveFileLog = logFileName
		newRun = &db.CronRun{
			CronJobID: int64(jobID),
			Status:    db.RunStatusRunning,
			Trigger:   trigger,
			StartedAt: now,
		}
		if err := database.Create(newRun).Error; err != nil {
			return err
		}
	}
	go func() {
		// remove old cron runs
		if err := db.PruneCronRuns(database, job.RunLogKeepCount, job.ID); err != nil {
			slog.Error("failed to prune cron runs", "error", err)
		}
	}()
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	if newRun != nil {
		go func() {
			ticker := time.NewTicker(time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-runCtx.Done():
					return
				case <-ticker.C:
					status, err := db.GetCronJobRunStatus(database, newRun.ID)
					if err != nil {
						slog.Error("failed to retrieve cron run status", "error", err)
						continue
					}
					newRun.Status = status

					if status != db.RunStatusRunning {
						slog.Info("cron run", "status", status, "run_id", newRun.ID)
						cancel()
						return
					}
				}
			}
		}()
	}
	slog.Info("starting cron job", "job_id", jobID, "command", command, "trigger", trigger, "is_managed", job.Managed, "save_file", saveFileLog)
	result, err := Run(runCtx, command, saveFileLog)
	slog.Info("finished cron job", "job_id", jobID, "status", result.ExitCode)

	job.LastRunAt = time.Now()

	if err := database.Save(job).Error; err != nil {
		slog.Error("failed to save cron job", "error", err)
	}
	if newRun != nil {

		if result.TimedOut {
			newRun.Status = db.RunStatusTimedOut
		} else if err != nil {
			newRun.Status = db.RunStatusFailed
		} else if result.ExitCode == 0 {
			newRun.Status = db.RunStatusSuccess
		}

		newRun.FinishedAt = time.Now()

		newRun.ExitCode = result.ExitCode
		if err := database.Save(newRun).Error; err != nil {
			return err
		}

	}
	if err != nil {
		slog.Error("failed to run command", "error", err, "command", command)
		return err
	}

	return nil
}
