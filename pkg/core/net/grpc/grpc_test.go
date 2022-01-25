package grpc

import (
	"testing"
)

func TestPools_Connect(t *testing.T) {
	type args struct {
		serviceName string
		addr        string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "test1",
			args: args{
				serviceName: "service-collect",
				addr:        "127.0.0.1:6664",
			},
		},
		//{
		//	name: "test2",
		//	args: args{
		//		serviceName: "service-collect-2",
		//		addr:        "http://127.0.0.1:6664,",
		//	},
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Gets()

			got, err := p.Connect(tt.args.serviceName, tt.args.addr)

			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			t.Logf("%+v", got)
		})
	}
}

func TestPool_health(t *testing.T) {
	tests := []struct {
		name        string
		serviceName string
		addr        string
	}{
		{
			name:        "test1",
			serviceName: "service-collect",
			addr:        "10.10.9.199:8061,http://10.10.9.199:8060",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Gets()
			pool, err := p.Connect(tt.serviceName, tt.addr)
			if err != nil {
				t.Error(err)
				t.FailNow()
			}
			t.Log(pool)
		})
	}
	select {}
}
