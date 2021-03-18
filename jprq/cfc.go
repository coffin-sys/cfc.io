package jprq

type Cfc struct {
	baseHost string
	tunnels map[string]*Tunnel
}

func New(baseHost string) Cfc {
	return Cfc{
		baseHost: baseHost,
		tunnels: make(map[string]*Tunnel),
	}
}
