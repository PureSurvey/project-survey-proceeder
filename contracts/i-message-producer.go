package contracts

type IMessageProducer interface {
	SendMessage(message []byte) error
	CloseConnection() error
}
