package nats

type ConnectionOpts struct {
	NKey string `json:"nKey"`
}

// kafka-console-consumer --bootstrap-server localhost:9092 -topic test --consumer.config ./.kafka/client.properties
// kafka-console-producer --broker-list localhost:9092 -topic test --producer.config ./.kafka/client.properties
