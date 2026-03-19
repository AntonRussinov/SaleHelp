package api

type NotificationSetting struct {
	Key            string `json:"key" validate:"required"`
	Enabled        bool   `json:"enabled"`
	Title          string `json:"title,omitempty"`
	Description    string `json:"description,omitempty"`
	Category       string `json:"category,omitempty"`
	DefaultEnabled bool   `json:"defaultEnabled,omitempty"`
}

type ListNotificationSettingsResponse struct {
	Items []NotificationSetting `json:"items"`
}

type UpdateNotificationSettingsItem struct {
	Key     string `json:"key" validate:"required"`
	Enabled bool   `json:"enabled"`
}

type UpdateNotificationSettingsRequest struct {
	Items []UpdateNotificationSettingsItem `json:"items" validate:"required,dive"`
}
