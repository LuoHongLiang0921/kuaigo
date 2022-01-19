package kstring

import (
	"fmt"
	"testing"
	"time"
)

func TestGenerateUUID(t *testing.T) {
	type args struct {
		seedTime time.Time
	}
	var tests []struct {
		name string
		args args
		want string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateUUID(tt.args.seedTime); got != tt.want {
				t.Errorf("GenerateUUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateID(t *testing.T) {
	var tests []struct {
		name string
		want string
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateID(); got != tt.want {
				t.Errorf("GenerateID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateUUID1(t *testing.T) {
	str := GenerateUUID(time.Now())
	str2 := GenerateID()
	fmt.Println(str)
	fmt.Println(str2)
}
