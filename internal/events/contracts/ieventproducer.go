package contracts

type IEventProducer interface {
	SendMessage(message []byte) error
	CloseConnection() error
}
