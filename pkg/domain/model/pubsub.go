package model

type PubSubRequest struct {
	DeliveryAttempt int64         `json:"deliveryAttempt"`
	Message         PubSubMessage `json:"message"`
	Subscription    string        `json:"subscription"`
}

type PubSubMessage struct {
	Attributes  map[string]string `json:"attributes"`
	Data        []byte            `json:"data"`
	MessageID   string            `json:"message_id"`
	PublishTime string            `json:"publish_time"`
}
