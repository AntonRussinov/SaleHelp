package cron

import (
	"gorm.io/gorm"
)

func DeleteOldSoftDeletedRecords(db *gorm.DB, log logger) {
	log.Info("DeleteOldSoftDeletedRecords started")

	// duration := time.Now().AddDate(0, 0, -7).Truncate(time.Minute) // 7 days ago

	// err := db.Unscoped().Where("deleted_at IS NOT NULL AND deleted_at < ?", duration).Delete(&model.Template{}).Error
	// if err != nil {
	// 	log.Error("Error deleting old soft deleted records for Template table: ", sl.Err(err))
	// }

	// err = db.Unscoped().Where("deleted_at IS NOT NULL AND deleted_at < ?", duration).Delete(&model.Application{}).Error
	// if err != nil {
	// 	log.Error("Error deleting old soft deleted records for Application table: ", sl.Err(err))
	// }

	// if err == nil {
	// 	log.Info("Old soft deleted records deleted successfully")
	// }
}
