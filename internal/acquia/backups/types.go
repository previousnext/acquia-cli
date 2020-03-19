package backups

import "time"

// BackupList contains a list of backups.
type BackupList []Backup

// Backup performed by the Acquia platform.
type Backup struct {
	ID        int64  `json:"id"`
	Completed time.Time  `json:"completed_at"`
	Environment Environment `json:"environment"`
	Database Database `json:"database"`
}

// Environment associated with a Backup.
type Environment struct {
	ID string `json:"id"`
	Name string `json:"name"`
}

// Database associated with a Backup.
type Database struct {
	ID int64 `json:"id"`
	Name string `json:"name"`
}