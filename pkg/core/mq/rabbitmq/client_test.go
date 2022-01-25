package rabbitmq

import (
	"context"
	"encoding/json"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"sync"
	"sync/atomic"

	"github.com/stretchr/testify/assert"

	"log"
	"testing"
)

var count int64

type Num struct {
	N int
}

func TestClient_Consume_Publish(t *testing.T) {
	ctx := context.Background()
	client := Config{
		Address: "amqp://admin:admin@127.0.0.1:5672",
		Mode:    constant.ModeRabbitmq,
		RunType: constant.RunTypePublish + "|" + constant.RunTypeConsumer,
		ConsumeConfig: ConsumeConfig{
			Queue: "tabby.test",
		},
	}.Build(ctx)
	var wg sync.WaitGroup
	_ = client.RegisterHandler(ctx, func(ctx context.Context, msg *mq.Message) error {
		klog.WithContext(ctx).Infof("%+v", string(msg.Body))
		var numMsg Num
		_ = json.Unmarshal(msg.Body, &numMsg)
		atomic.AddInt64(&count, int64(numMsg.N))
		msg.Ack()
		wg.Done()
		return nil
	})
	done := make(chan struct{})
	finish := make(chan struct{})
	go func() {
		<-done
		_ = client.Consume(ctx)
	}()
	go func() {
		for i := 0; i < 1200; i++ {
			wg.Add(1)
			msgRaw, _ := json.Marshal(&Num{N: i})
			_, err := client.Publish(ctx, "media.test.exchange", &mq.Message{
				Header: map[string]string{"test": "testheader"},
				Body:   msgRaw,
			})
			if err == nil {
				log.Println(string(msgRaw))
			}
		}
		close(finish)
		close(done)

	}()
	<-finish
	wg.Wait()
	_ = client.Stop()
	assert.EqualValues(t, 719400, count)
}

func TestClient_Publish(t *testing.T) {
	ctx := context.Background()
	c := Config{
		Address:       "amqp://admin:admin@127.0.0.1:5672",
		Mode:          constant.ModeRabbitmq,
		RunType:       constant.RunTypePublish,
		PublishConfig: PublishConfig{Exchange: "media.test.exchange"},
	}
	client := c.Build(ctx)
	for i := 0; i < 12; i++ {
		msg, _ := json.Marshal(&Num{N: i})
		_, err := client.Publish(ctx, c.PublishConfig.Exchange, &mq.Message{
			Header: map[string]string{"test": "testheader"},
			Body:   msg,
		})
		if err == nil {
			log.Println(string(msg))
		}
	}
}

func TestClient_Consume(t *testing.T) {
	ctx := context.Background()
	client := Config{
		Address: "amqp://admin:admin@127.0.0.1:5672",
		Mode:    constant.ModeRabbitmq,
		RunType: constant.RunTypeConsumer,
		ConsumeConfig: ConsumeConfig{
			Queue: "tabby.test",
		},
	}.Build(ctx)
	_ = client.RegisterHandler(ctx, func(ctx context.Context, msg *mq.Message) error {
		klog.WithContext(ctx).Infof("%+v", string(msg.Body))
		var numMsg Num
		_ = json.Unmarshal(msg.Body, &numMsg)
		atomic.AddInt64(&count, int64(numMsg.N))
		msg.Ack()
		return nil
	})
	_ = client.Consume(ctx)
	_ = client.Stop()
}
