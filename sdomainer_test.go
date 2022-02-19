package sdomainer

import (
	"reflect"
	"testing"
)

func TestRun(t *testing.T) {
	tests := []struct {
		name   string
		domain string
		wd     string
		want   []string
	}{
		{"Valid one", "example.com", "wd.txt", []string{"www.example.com"}},
		{"Inexistent one", "shouldnotexistfr.xyx", "wd.txt", []string{}},
		{"Bad wordlist", "example.com", "xxx.txt", []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSDomainer(tt.domain, tt.wd, WithGoroutines(100))
			got, _ := s.Run()
			if len(tt.want) != len(got) && !reflect.DeepEqual(tt.want, got) {
				t.Errorf("expected: %v; got %v", tt.want, got)
			}
		})
	}
}
