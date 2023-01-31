package antivpn

type Source interface {
	GetResult(ip string) bool
}
