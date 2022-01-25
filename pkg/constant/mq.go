// @Description mq 常量

package constant

const (
	ConfigRootKey = "mq"
)

const (
	// ModeKafka kafka 模式
	ModeKafka = "kafka"
	// ModeRabbitmq rabbitmq 模式
	ModeRabbitmq = "rabbitmq"
	// ModeRocketmq rocketmq 模式
	ModeRocketmq = "rocketmq"
)

const (
	RunTypePublish  = "publish"
	RunTypeConsumer = "consumer"

	RunTypePublishConsumer = RunTypePublish + "|" + RunTypeConsumer
)
