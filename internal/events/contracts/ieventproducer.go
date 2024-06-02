package contracts

type IEventProducer interface {
	Init() error
	AsyncSendMessage(message []byte)
	CloseConnection()
}
