package structure

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type FilterTokenActivities struct {
	BaseFilters
	InscriptionID *string
	ProjectID *string
	Types []int
}

func (f *FilterTokenActivities) CreateFilter(r *http.Request) error {
	inscriptionID := r.URL.Query().Get("inscription_id")
	projectID := r.URL.Query().Get("project_id")
	typesRaw := r.URL.Query().Get("types")
	types := []int{}
	if typesRaw != "" {
		typesStr := strings.Split(typesRaw, ",")
		for _, typeStr := range typesStr {
			typeInt, err := strconv.Atoi(typeStr)
			if err != nil {
				return errors.WithStack(err)
			}
			types = append(types, typeInt)
		}
	}
	f.Types = types
	if inscriptionID != "" {
		f.InscriptionID = &inscriptionID
	}
	if projectID != "" {
		f.ProjectID = &projectID
	}
	return nil
}
