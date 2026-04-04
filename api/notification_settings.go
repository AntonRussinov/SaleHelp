package api

type NotificationSetting struct {
	Key            string `json:"key"                      validate:"required" example:"fbs_new_order"`
	Enabled        bool   `json:"enabled"                                      example:"true"`
	Title          string `json:"title,omitempty"                              example:"Новый заказ FBS"`
	Description    string `json:"description,omitempty"                        example:"Уведомление о новом сборочном задании FBS"`
	Category       string `json:"category,omitempty"                           example:"orders"`
	DefaultEnabled bool   `json:"defaultEnabled,omitempty"                     example:"true"`
}

type ListNotificationSettingsResponse struct {
	Items []NotificationSetting `json:"items"`
}

type UpdateNotificationSettingsItem struct {
	Key     string `json:"key"     validate:"required" example:"fbs_new_order"`
	Enabled bool   `json:"enabled"                     example:"false"`
}

type UpdateNotificationSettingsRequest struct {
	Items []UpdateNotificationSettingsItem `json:"items" validate:"required,dive"`
}
