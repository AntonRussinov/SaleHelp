package gormscopes

import (
	"context"
	"errors"
	"salehelp/internal/storage/params/filter"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type (
	FilterField struct {
		FieldName        string
		ValueType        filter.ValueType
		CompareType      filter.CompareType
		AdditionalValues []any
	}

	FilterMap map[string]FilterField
)

var ErrUnknownValueType = errors.New("unknown value type")

func Filter(ctx context.Context, filterMap FilterMap) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		filterFields := filter.GetFromContext(ctx)
		if len(filterFields) == 0 {
			return db
		}

		query := db

		for _, filterField := range filterFields {
			if filterField.Field == "tag" {
				query = convertFilterForTag(query, filterField, filterMap)
			} else {
				cond, val, ok := makeFilterCondition(filterField, filterMap)
				if !ok {
					continue
				}

				query = query.Where(cond, val...)
			}
		}

		return query
	}
}

func convertFilterForTag(db *gorm.DB, filterField filter.Param, filterMap FilterMap) *gorm.DB {
	conditions := []string{}
	combinedCondition := ""
	args := []any{}

	for _, tagName := range filterField.Values {
		cond, _, ok := makeFilterCondition(filterField, filterMap)
		if !ok {
			continue
		}

		conditions = append(conditions, cond)
		args = append(args, tagName)
	}

	if len(conditions) > 0 {
		combinedCondition = "(" + strings.Join(conditions, " OR ") + ")"
	}

	return db.Where(combinedCondition, args...)
}

func convertFilterValue(value string, valueType filter.ValueType) (any, error) {
	var converted any
	var err error

	switch valueType {
	case filter.ValueTypeBool:
		converted, err = strconv.ParseBool(value)
	case filter.ValueTypeInt:
		converted, err = strconv.Atoi(value)
	case filter.ValueTypeFloat:
		converted, err = strconv.ParseFloat(value, 64)
	case filter.ValueTypeString:
		converted = value
	case filter.ValueTypeUUID:
		converted, err = uuid.Parse(value)
	case filter.ValueTypeDate:
		converted, err = time.Parse("2006-01-02", value)
	default:
		err = ErrUnknownValueType
	}

	if err != nil {
		return nil, err
	}

	return converted, nil
}

func makeFilterCondition(filterField filter.Param, filterMap FilterMap) (string, []any, bool) {
	field, ok := filterMap[filterField.Field]
	if !ok {
		return "", nil, false
	}

	name := field.FieldName
	if name == "" {
		name = filterField.Field
	}

	var filterValue any
	var err error

	compareType := filterField.CompareType
	if field.CompareType != "" {
		compareType = field.CompareType
	}

	if compareType == filter.CompareTypeIn {
		filterValue = filterField.Values
	} else {
		filterValue, err = convertFilterValue(filterField.Value, field.ValueType)
	}

	if err != nil {
		return "", nil, false
	}

	condition := buildFilterCondition(name, compareType)

	values := make([]any, 0, len(field.AdditionalValues)+1)
	values = append(values, field.AdditionalValues...)
	values = append(values, filterValue)

	return condition, values, true
}

func compareTypesStatements() map[filter.CompareType]string {
	return map[filter.CompareType]string{
		filter.CompareTypeEqual:          " = ?",
		filter.CompareTypeNotEqual:       " != ?",
		filter.CompareTypeGreaterThan:    " > ?",
		filter.CompareTypeGreaterOrEqual: " >= ?",
		filter.CompareTypeLessThan:       " < ?",
		filter.CompareTypeLessOrEqual:    " <= ?",
		filter.CompareTypeIn:             " IN (?)",
		filter.CompareTypeNotIn:          " NOT IN (?)",
		filter.CompareTypeLike:           " ILIKE '%' || ? || '%'",
		filter.CompareTypeNotLike:        " NOT ILIKE '%' || ? || '%'",
		filter.CompareTypeStartsWith:     " ILIKE ? || '%'",
		filter.CompareTypeEndsWith:       " ILIKE '%' || ?",
	}
}

func buildFilterCondition(name string, compareType filter.CompareType) string {
	statements := compareTypesStatements()

	statement, ok := statements[compareType]
	if !ok {
		statement = " = ?"
	}

	return name + statement
}
