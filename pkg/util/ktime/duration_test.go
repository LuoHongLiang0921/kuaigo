// @Description

package ktime

import (
	"testing"
	"time"
)

func TestDuration(t *testing.T) {
	type args struct {
		str string
	}
	var tests []struct {
		name string
		args args
		want time.Duration
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Duration(tt.args.str); got != tt.want {
				t.Errorf("Duration() = %v, want %v", got, tt.want)
			}
		})
	}
}
