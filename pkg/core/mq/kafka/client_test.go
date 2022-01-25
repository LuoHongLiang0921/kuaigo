package kafka

import (
	"context"
	"encoding/json"
	"github.com/LuoHongLiang0921/kuaigo/pkg/constant"
	"github.com/LuoHongLiang0921/kuaigo/pkg/core/mq"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/klog"
	"github.com/LuoHongLiang0921/kuaigo/pkg/util/ktime"
	"github.com/stretchr/testify/assert"
	"log"
	"sync"
	"sync/atomic"
	"testing"
)

var count int64

type Num struct {
	N int
}

func TestConsumePublish(t *testing.T) {
	ctx := context.Background()
	c := Config{
		Mode:    constant.ModeKafka,
		RunType: constant.RunTypePublish + "|" + constant.RunTypeConsumer,
		Brokers: []string{"127.0.0.1:9092"},
		Version: "",
		ConsumerConfig: &ConsumerConfig{
			Topic:          []string{"test_kafka"},
			GroupID:        "test_kafka_group",
			OffsetsInitial: OffsetOldest,
			Assinor:        AssinorRoundRobin,
		},
		PublishConfig: &PublishConfig{
			PublishTimeout: ktime.Duration("5s"),
			Async:          true,
			RequiredAcks:   WaitForLocal,
			Backoff:        ktime.Duration("2s"),
			Partitioner:    HashPartitioner,
		},
	}
	client := c.Build(ctx)
	var wg sync.WaitGroup
	client.logger = klog.KuaigoLogger
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

func TestClient_Publish(t *testing.T) {
	ctx := context.Background()
	c := Config{
		Mode:    constant.ModeKafka,
		RunType: constant.RunTypePublish,
		Brokers: []string{"127.0.0.1:9092"},
		PublishConfig: &PublishConfig{
			PublishTimeout: 0,
			Async:          true,
			RequiredAcks:   0,
			Backoff:        0,
			Partitioner:    "",
		},
	}
	client := c.Build(ctx)
	for i := 0; i < 1200; i++ {
		msgRaw, _ := json.Marshal(&Num{N: i})
		_, err := client.Publish(ctx, "test_kafka", &mq.Message{
			Header: map[string]string{"test": "testheader"},
			Body:   msgRaw,
		})
		if err == nil {
			log.Println(string(msgRaw))
		}
	}
}
