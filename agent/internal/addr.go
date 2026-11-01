package internal

type Addr struct{}

func (Addr) Network() string { return "ssh" }
func (Addr) String() string  { return "ssh" }
