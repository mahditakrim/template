package transport

type Transport interface {
	Run() error
	Shutdown() error
}
