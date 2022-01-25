package mqmanager

import (
	"context"
	"encoding/json"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/test"
	"log"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

var count int64

type Num struct {
	N int
}

func TestKafkaPublishConsume(t *testing.T) {
	err := test.InitTestForFile()
	if assert.NoError(t, err) {
		ctx := context.Background()
		mqs := Load(ctx)
		client := mqs["tkafka"]
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
				_, err := client.Publish(ctx, "test_kafka", &mq.Message{
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
}

func TestRabbitPublishConsume(t *testing.T) {
	err := test.InitTestForFile()
	if assert.NoError(t, err) {
		ctx := context.Background()
		mqs := Load(ctx)
		client := mqs["trabbit"]
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
}
