package access

type Action string

const (
	ActionManage Action = "MANAGE"
	ActionRead   Action = "READ"
	ActionCreate Action = "CREATE"
	ActionUpdate Action = "UPDATE"
	ActionDelete Action = "DELETE"
)
