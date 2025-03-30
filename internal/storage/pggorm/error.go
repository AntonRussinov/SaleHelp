package pggorm

import (
	"fmt"
	"salehelp/internal/storage"

	"github.com/jackc/pgx/v5/pgconn"
)

func foreignKeyConstraints() map[string]string {
	fieldmap := map[string]string{
		"applications_function_group_id_fkey":                "functionGroup",
		"applications_logo_id_fkey":                          "logoFile",
		"application_documents_file_id_fkey":                 "file",
		"templates_logo_id_fkey":                             "logo",
		"chat_documents_file_id_fkey":                        "file",
		"application_users_user_id_fkey":                     "user",
		"applications_llm_id_fkey":                           "llm",
		"application_function_groups_function_group_id_fkey": "function_group",
		"user_feedback_categories_category_id_fkey":          "category",
	}

	return fieldmap
}

func uniqueConstraints() map[string]string {
	fieldmap := map[string]string{
		"idx_applications_name": "name",
		"idx_templates_name":    "name",
		"idx_applications_slug": "slug",
		"idx_templates_slug":    "slug",
	}

	return fieldmap
}

func uniqueConstraintError(op string, err *pgconn.PgError) error {
	fieldmap := uniqueConstraints()
	name := fieldmap[err.ConstraintName]

	if name == "" {
		return wrapError(op, err)
	}

	return wrapError(op, &storage.DuplicateError{
		FieldName: name,
		Func:      op,
		Err:       err,
	})
}

func foreignKeyConstraintError(op string, err *pgconn.PgError) error {
	fkConstrains := foreignKeyConstraints()
	fieldName := fkConstrains[err.ConstraintName]

	if fieldName == "" {
		return wrapError(op, err)
	}

	return wrapError(op, &storage.ForeignKeyError{
		FieldName: fieldName,
		Func:      op,
		Err:       err,
	})
}

func parsePgConnError(op string, err *pgconn.PgError) error {
	state := err.SQLState()
	switch state {
	case "23505":
		return uniqueConstraintError(op, err)
	case "23503":
		return foreignKeyConstraintError(op, err)
	default:
		return wrapError(op, err)
	}
}

func wrapError(op string, err error) error {
	return fmt.Errorf("%s: %w", op, err)
}
