package processor

import (
	"EduCommentSync/internal/models"
	"reflect"
	"sort"
	"testing"
)

func TestGetWorkNamesUnique(t *testing.T) {
	tests := []struct {
		name     string
		comments []*models.StudentComment
		want     []string
	}{
		{
			name: "No duplicates",
			comments: []*models.StudentComment{
				{WorkName: "Homework 1"},
				{WorkName: "Homework 2"},
				{WorkName: "Homework 3"},
			},
			want: []string{"Homework 1", "Homework 2", "Homework 3"},
		},
		{
			name: "With duplicates",
			comments: []*models.StudentComment{
				{WorkName: "Homework 1"},
				{WorkName: "Homework 2"},
				{WorkName: "Homework 1"},
				{WorkName: "Homework 3"},
				{WorkName: "Homework 2"},
			},
			want: []string{"Homework 1", "Homework 2", "Homework 3"},
		},
		{
			name:     "Empty input",
			comments: []*models.StudentComment{},
			want:     []string{},
		},
		{
			name: "Single work name",
			comments: []*models.StudentComment{
				{WorkName: "Homework 1"},
				{WorkName: "Homework 1"},
				{WorkName: "Homework 1"},
			},
			want: []string{"Homework 1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetWorkNamesUnique(tt.comments)
			sort.Strings(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetWorkNamesUnique() = %v, want %v", got, tt.want)
			}
		})
	}
}
