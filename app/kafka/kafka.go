package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"microservice/app"
	"microservice/app/core"
	"path"
)

var k *KafkaService
var logger core.Logger

type KafkaService struct {
	producer      sarama.SyncProducer
	consumer      sarama.Consumer
	consumerGroup sarama.ConsumerGroup
}

func InitKafka(l core.Logger) error {

	logger = l

	enabled := viper.GetBool("kafka.enabled")
	if !enabled {
		return nil
	}

	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true

	brokers := viper.GetStringSlice("kafka.brokers")

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return err
	}

	group := viper.GetString("kafka.group")
	consumerGroup, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return err
	}

	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return err
	}

	k = &KafkaService{
		producer:      producer,
		consumer:      consumer,
		consumerGroup: consumerGroup,
	}

	return nil
}

type KafkaTopic[T any] struct {
	topic        string
	encoder      Encoder
	topicStorage *app.Storage
}

func Topics() ([]string, error) {
	return k.consumer.Topics()
}

func Topic[T any](topic string, encoder ...Encoder) (*KafkaTopic[T], error) {
	var enc Encoder
	var t interface{}
	t = new(T)
	switch t.(type) {
	case *string:
		enc = StringEncoder()
	case *[]byte:
		enc = ByteEncoder()
	default:
		enc = JsonEncoder()
	}
	if len(encoder) != 0 {
		enc = encoder[0]
	}

	// Storage for offsets
	storage, err := app.NewStorage(path.Join("offsets", topic))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create storage for kafka topic %s", topic)
	}

	return &KafkaTopic[T]{
		topic:        topic,
		encoder:      enc,
		topicStorage: storage,
	}, nil
}

func (t *KafkaTopic[T]) Produce(obj T) error {

	msg, err := t.encoder.Decode(obj)
	if err != nil {
		return err
	}

	partition := int32(0) // Partition = 0!
	message := &sarama.ProducerMessage{
		Topic:     t.topic,
		Partition: partition,
		Value:     sarama.ByteEncoder(msg),
	}

	_, offset, err := k.producer.SendMessage(message)
	if err != nil {
		return err
	}
	logger.Debug("Kafka message sent to topic %s (offset=%d)", t.topic, offset)
	return nil
}

func (t *KafkaTopic[T]) StartPolling() (chan *Message[T], error) {

	//partitionList, err := k.consumer.Partitions(t.topic)
	//if err != nil {
	//	return nil, err
	//}

	// Get last offset from storage
	initialOffset, err := t.topicStorage.GetInt64(t.topic)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot get initial offset of topic %s in storage", t.topic)
	}
	if initialOffset == nil {
		x := sarama.OffsetOldest
		initialOffset = &x
	}

	messages := make(chan *Message[T], 1)
	partition := int32(0) // Partition 0!
	pc, err := k.consumer.ConsumePartition(t.topic, partition, *initialOffset)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot consume broker partition %d", partition)
	}

	go func(pc sarama.PartitionConsumer) {
		logger.Info("KAFKA: starting polling messages: Topic <%s>, Offset: <%d>", t.topic, *initialOffset)
		for message := range pc.Messages() {
			msg := &Message[T]{
				Details: message,
			}
			err := t.encoder.Encode(message.Value, &msg.Value)
			if err != nil {
				//e have incorrect format - skip it
				logger.ErrorWrap(err, "cannot encode kafka message to receiver type")
				err := t.CommitOffset(msg)
				if err != nil {
					logger.ErrorWrap(err, "cannot commit offset after failed encoding in topic %s", t.topic)
				}
				continue
			}
			messages <- msg
		}
	}(pc)

	return messages, nil
}

func (t *KafkaTopic[T]) CommitOffset(msg *Message[T]) error {
	lastOffset, err := t.topicStorage.GetInt64(t.topic)
	if err != nil {
		return errors.Wrapf(err, "cannot get offset of topic %s in storage", t.topic)
	}
	if lastOffset == nil {
		firstOffset := int64(0)
		lastOffset = &firstOffset
	}
	newOffset := *lastOffset + 1

	// All last offsets before msg.Offset should be handled
	// If not then we have problems!
	if newOffset != msg.Details.Offset+1 {
		logger.Error("newOffset != msg.Offset+1 in topic %s", t.topic)
	}

	// Anyway we commit with message offset+1
	committedOffset := msg.Details.Offset + 1
	err = t.topicStorage.PutInt64(t.topic, committedOffset)
	if err != nil {
		return errors.Wrapf(err, "cannot put new offset %d of topic %s in storage", committedOffset, t.topic)
	}

	return nil
}
