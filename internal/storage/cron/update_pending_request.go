package cron

import (
	"gorm.io/gorm"
)

func CronUpdateRequest(db *gorm.DB, log logger) {
	//updateOldPendingRequests(db, log, model.LimitRequest{}, string(model.LimitRequestStatusRejected))
	//updateOldPendingRequests(db, log, model.QuotaRequest{}, string(model.QuotaRequestStatusRejected))
}

// func updateOldPendingRequests[T any](db *gorm.DB, log logger, tableModel T, newStatus string) {
// 	log.Info("UpdateOldPendingRequests started for model")

// 	duration := time.Now().AddDate(0, 0, -3).Truncate(time.Minute) // 3 days ago
// 	newState := map[string]any{
// 		"Status": newStatus,
// 	}
// 	err := db.Model(tableModel).
// 		Where("updated_at < ? AND status = ?", duration, "Pending").
// 		Updates(newState).Error

// 	if err != nil {
// 		log.Error("Error updating old pending requests for model: ", sl.Err(err))
// 	} else {
// 		log.Info(fmt.Sprintf("Old pending requests updated successfully for %T", tableModel))
// 	}
// }
