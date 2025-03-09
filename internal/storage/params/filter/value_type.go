package filter

type ValueType string

const (
	ValueTypeString ValueType = "string"
	ValueTypeInt    ValueType = "int"
	ValueTypeFloat  ValueType = "float"
	ValueTypeBool   ValueType = "bool"
	ValueTypeUUID   ValueType = "uuid"
	ValueTypeDate   ValueType = "date"
)
