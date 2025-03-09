package filter

type CompareType string

const (
	CompareTypeEqual          CompareType = "eq"
	CompareTypeNotEqual       CompareType = "ne"
	CompareTypeGreaterThan    CompareType = "gt"
	CompareTypeGreaterOrEqual CompareType = "gte"
	CompareTypeLessThan       CompareType = "lt"
	CompareTypeLessOrEqual    CompareType = "lte"
	CompareTypeIn             CompareType = "in"
	CompareTypeNotIn          CompareType = "nin"
	CompareTypeLike           CompareType = "like"
	CompareTypeNotLike        CompareType = "nlike"
	CompareTypeStartsWith     CompareType = "starts"
	CompareTypeEndsWith       CompareType = "ends"
)

func validCompareTypes() map[CompareType]struct{} {
	return map[CompareType]struct{}{
		CompareTypeEqual:          {},
		CompareTypeNotEqual:       {},
		CompareTypeGreaterThan:    {},
		CompareTypeGreaterOrEqual: {},
		CompareTypeLessThan:       {},
		CompareTypeLessOrEqual:    {},
		CompareTypeIn:             {},
		CompareTypeNotIn:          {},
		CompareTypeLike:           {},
		CompareTypeNotLike:        {},
		CompareTypeStartsWith:     {},
		CompareTypeEndsWith:       {},
	}
}

func (compType CompareType) IsValid() bool {
	validTypes := validCompareTypes()

	_, ok := validTypes[compType]

	return ok
}
