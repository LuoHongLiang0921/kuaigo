// @Description

package ktime

import (
	"testing"
	"time"
)


func TestTime_Compare(t1 *testing.T) {
	type fields struct {
		Time time.Time
	}
	type args struct {
		operator string
		tt       *Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "大于",
			fields:fields{Time: time.Now()},
			args: args{
				operator: ">",
				tt:       &Time{Parse("2020-08-05").Time},
			},
			want: false,
		},
		{
			name: "小于",
			fields:fields{Time: time.Now()},
			args: args{
				operator: "<",
				tt:       &Time{Parse("2022-08-05").Time},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Time{
				Time: tt.fields.Time,
			}
			if got := t.Compare(tt.args.operator, tt.args.tt); got != tt.want {
				t1.Errorf("Compare() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTime_Between
//  @Description  是否在两个时间之间(不包括这两个时间)
//  @Param t1
func TestTime_Between(t1 *testing.T) {
	type fields struct {
		Time time.Time
	}
	type args struct {
		start *Time
		end   *Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "在范围内",
			fields:fields{Time: time.Now()},
			args: args{
				start: &Time{Parse("2020-08-05").Time},
				end:   &Time{Parse("2022-08-05").Time},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Time{
				Time: tt.fields.Time,
			}
			if got := t.Between(tt.args.start, tt.args.end); got != tt.want {
				t1.Errorf("Between() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestTime_BetweenIncludedStartTime
//  @Description  是否在两个时间之间(包括开始时间)
//  @Param t1
func TestTime_BetweenIncludedStartTime(t1 *testing.T) {
	type fields struct {
		Time time.Time
	}
	type args struct {
		start *Time
		end   *Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "在范围内",
			fields:fields{Time: time.Now()},
			args: args{
				start: &Time{Parse("2020-08-05").Time},
				end:   &Time{Now().Time},
			},
			want: true,
		},
		{
			name: "不在范围内",
			fields:fields{Time: time.Now()},
			args: args{
				start: &Time{Parse("2020-08-05").Time},
				end:   (&Time{Now().Time}).Yesterday(),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := &Time{
				Time: tt.fields.Time,
			}
			if got := t.BetweenIncludedStartTime(tt.args.start, tt.args.end); got != tt.want {
				t1.Errorf("BetweenIncludedStartTime() = %v, want %v", got, tt.want)
			}
		})
	}
}