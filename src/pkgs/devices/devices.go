package devices

type Device struct {
	MGMTAddress string
	Interfaces  []Interface
}

type Interface struct {
	Address string
	VRF     string
	VLAN    int
}
