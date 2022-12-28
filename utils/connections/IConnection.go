package connections

type IConnection interface {
	Connect() interface{}
	Disconnect() error
	GetType() interface{}
}
