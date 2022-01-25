package kgrpc

import (
	"context"
	"fmt"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"

	"net"

	"runtime"
	"strings"
	"time"
)

func defaultStreamServerInterceptor(ctx context.Context, logger *klog.Logger, slowQueryThresholdInMilli int64) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		var beg = time.Now()
		var fields = make([]klog.Field, 0, 8)
		var event = "normal"
		defer func() {
			if slowQueryThresholdInMilli > 0 {
				if int64(time.Since(beg))/1e6 > slowQueryThresholdInMilli {
					event = "slow"
				}
			}

			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}
				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, klog.FieldParams(stack))
				event = "recover"
			}

			fields = append(fields,
				klog.Any("grpc interceptor type", "unary"),
				klog.FieldMethod(info.FullMethod),
				klog.FieldCost(time.Since(beg)),
				klog.FieldName(event),
			)

			for key, val := range getPeer(stream.Context()) {
				fields = append(fields, klog.Any(key, val))
			}

			if err != nil {
				fields = append(fields, zap.String("err", err.Error()))
				logger.WithContext(ctx).Error("access", fields...)
				return
			}
			logger.WithContext(ctx).Info("access", fields...)
		}()
		return handler(srv, stream)
	}
}

func defaultUnaryServerInterceptor(ctx context.Context, logger *klog.Logger, slowQueryThresholdInMilli int64) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		var beg = time.Now()
		var fields = make([]klog.Field, 0, 8)
		var event = "normal"
		defer func() {
			if slowQueryThresholdInMilli > 0 {
				if int64(time.Since(beg))/1e6 > slowQueryThresholdInMilli {
					event = "slow"
				}
			}
			if rec := recover(); rec != nil {
				switch rec := rec.(type) {
				case error:
					err = rec
				default:
					err = fmt.Errorf("%v", rec)
				}

				stack := make([]byte, 4096)
				stack = stack[:runtime.Stack(stack, true)]
				fields = append(fields, klog.FieldStack(stack))
				event = "recover"
			}

			fields = append(fields,
				klog.Any("grpc interceptor type", "unary"),
				klog.FieldMethod(info.FullMethod),
				klog.FieldCost(time.Since(beg)),
				klog.FieldEvent(event),
			)

			for key, val := range getPeer(ctx) {
				fields = append(fields, klog.Any(key, val))
			}

			if err != nil {
				fields = append(fields, zap.String("err", err.Error()))
				logger.WithContext(ctx).Error("access", fields...)
				return
			}
			logger.WithContext(ctx).Info("access", fields...)
		}()
		return handler(ctx, req)
	}
}

func getClientIP(ctx context.Context) (string, error) {
	pr, ok := peer.FromContext(ctx)
	if !ok {
		return "", fmt.Errorf("[getClinetIP] invoke FromContext() failed")
	}
	if pr.Addr == net.Addr(nil) {
		return "", fmt.Errorf("[getClientIP] peer.Addr is nil")
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	return addSlice[0], nil
}

func getPeer(ctx context.Context) map[string]string {
	var peerMeta = make(map[string]string)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if val, ok := md["aid"]; ok {
			peerMeta["aid"] = strings.Join(val, ";")
		}
		var clientIP string
		if val, ok := md["client-ip"]; ok {
			clientIP = strings.Join(val, ";")
		} else {
			ip, err := getClientIP(ctx)
			if err == nil {
				clientIP = ip
			}
		}
		peerMeta["clientIP"] = clientIP
		if val, ok := md["client-host"]; ok {
			peerMeta["host"] = strings.Join(val, ";")
		}
	}
	return peerMeta

}
