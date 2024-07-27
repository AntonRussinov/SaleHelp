package access

type Resource string

const (
	ResourceAll           Resource = "ALL"
	ResourceUsers         Resource = "Users"
	ResourceUserQuotas    Resource = "UserQuotas"
	ResourceQuotaRequests Resource = "QuotaRequests"
	ResourceUserLimits    Resource = "UserLimits"
	ResourceLimitRequests Resource = "ResourceLimitRequests"
)
