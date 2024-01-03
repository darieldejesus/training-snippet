package main

import (
	"testing"
	"time"

	"snippet.darieldejesus.com/internal/assert"
)

func TestHumanDate(t *testing.T) {
	tests := []struct {
		test     string
		tm       time.Time
		expected string
	}{
		{
			test:     "UTC",
			tm:       time.Date(2022, 12, 7, 10, 15, 0, 0, time.UTC),
			expected: "07 Dec 2022 at 10:15",
		},
		{
			test:     "Empty",
			tm:       time.Time{},
			expected: "",
		},
		{
			test:     "CET",
			tm:       time.Date(2022, 12, 7, 10, 15, 0, 0, time.FixedZone("CET", 1*60*60)),
			expected: "07 Dec 2022 at 09:15",
		},
	}

	for _, test := range tests {
		t.Run(test.test, func(t *testing.T) {
			hd := humanDate(test.tm)
			if hd != test.expected {

				assert.Equal(t, hd, test.expected)
			}
		})
	}
}
