package backups

// GetLatest backup from a list.
func GetLatest(list *BackupList) Backup {
	var backup Backup

	for _, item := range *list {
		if item.Completed.After(backup.Completed) {
			backup = item
		}
	}

	return backup
}