package delivery

type IDelivery interface {
	StartServer()
}

type AddedServer struct {
	Server IDelivery
	Enabled bool
}