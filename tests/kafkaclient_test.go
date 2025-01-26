package tests

import (
	"testing"
	"polling_websocket/pkg/infra/brokerclient"
	"polling_websocket/mocks"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

func TestKafkaClientImpl_Produce(t *testing.T) {
	type fields struct {
		producer *kafka.Producer
	}
	type args struct {
		topic string
		key   []byte
		value []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := &KafkaClientImpl{
				producer: tt.fields.producer,
			}
			if err := k.Produce(tt.args.topic, tt.args.key, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("KafkaClientImpl.Produce() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
