package transport

type Transport interface {
	Run(string) error
	Shutdown() error
}
