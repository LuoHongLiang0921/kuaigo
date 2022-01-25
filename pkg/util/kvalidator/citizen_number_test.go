// @Description

package kvalidator

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestIsCitizenNo18(t *testing.T) {
	//可用自己的身份证号码测试
	str := "111111111111111111"
	strByte := []byte(str)
	fmt.Println(IsCitizenNo18(&strByte))
}

func TestIsCitizenNo(t *testing.T) {
	//可用自己的身份证号码测试
	str := "451121199004236209"
	strByte := []byte(str)
	fmt.Println(IsCitizenNo(&strByte))
}

func TestGetCitizenNoInfo(t *testing.T) {
	type args struct {
		citizenNo []byte
	}
	formatTimeStr := "1990-04-23 00:00:00"
	formatTime, err := time.Parse("2006-01-02 15:04:05", formatTimeStr)
	if err == nil {
		fmt.Println(formatTime) //打印结果：2017-04-11 13:33:37 +0000 UTC
	}

	tests := []struct {
		name         string
		args         args
		wantErr      error
		wantBirthday time.Time
		wantSex      string
		wantAddress  string
	}{
		{"test", args{[]byte("451121199004236209")}, nil, formatTime, "女", "广西壮族自治区"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotErr, gotBirthday, gotSex, gotAddress := GetCitizenNoInfo(tt.args.citizenNo)
			if !reflect.DeepEqual(gotErr, tt.wantErr) {
				t.Errorf("GetCitizenNoInfo() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
			if !reflect.DeepEqual(gotBirthday, tt.wantBirthday) {
				t.Errorf("GetCitizenNoInfo() gotBirthday = %v, want %v", gotBirthday, tt.wantBirthday)
			}
			if gotSex != tt.wantSex {
				t.Errorf("GetCitizenNoInfo() gotSex = %v, want %v", gotSex, tt.wantSex)
			}
			if gotAddress != tt.wantAddress {
				t.Errorf("GetCitizenNoInfo() gotAddress = %v, want %v", gotAddress, tt.wantAddress)
			}
		})
	}
}
