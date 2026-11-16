package db

type CronRunLog struct {
	ID        int64  `json:"id" gorm:"primaryKey"`
	RunID     int64  `json:"run_id"`
	Line      string `json:"line"`
	Timestamp int64  `json:"timestamp"`
}
