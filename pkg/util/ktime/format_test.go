// @Description

package ktime

import (
	"fmt"
	"testing"
	"time"
)

func TestTime_Format(t1 *testing.T) {
	type fields struct {
		Time time.Time
	}
	type args struct {
		format string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name:   "Y-m-d",
			fields: fields{Time: time.Now()},
			args:   args{format: "Y-m-d"},
			want:   "",
		},
		{
			name:   "Y-M-D",
			fields: fields{Time: time.Now()},
			args:   args{format: "Y-M-D"},
			want:   "",
		},
		{
			name:   "Y-m-d H:i:s",
			fields: fields{Time: time.Now()},
			args:   args{format: "Y-m-d H:i:s"},
			want:   "",
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Time{
				Time: tt.fields.Time,
			}
			if got := t.Format(tt.args.format); got != tt.want {
				fmt.Printf("Format() = %s\n", got)
			}
		})
	}
}

func TestTime_other(t1 *testing.T) {
	var time Time = Time{Time: time.Now()}

	fmt.Printf("Tomorrow= %s\n", time.Tomorrow().ToDateString())
	fmt.Printf("Yesterday= %s\n", time.Yesterday().ToDateString())

}