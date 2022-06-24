package main

import (
	"testing"
	"time"

	"snippetbox.bimasenaputra/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		name string
		input time.Time
		expected string
	} {
		{
			name: "UTC",
			input: time.Date(2022, 3, 17, 10, 15, 0, 0, time.UTC),
			expected: "17 Mar 2022 at 10:15",
		},
		{
			name: "Empty",
			input: time.Time{},
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func (t *testing.T)  {
			actual := humanDate(test.input)
			assert.Equal(t, actual, test.expected)
		})
	}
}