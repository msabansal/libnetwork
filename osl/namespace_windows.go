package osl

import "testing"

// var (
// 	// ErrNotImplemented is for platforms which don't implement sandbox
// 	ErrNotImplemented = errors.New("not implemented")
// )

// type winSandbox struct {
// 	key         string
// 	compartment *hcsshim.Compartment
// }

// func (sandbox *winSandbox) Key() string {
// 	return sandbox.key
// }

// func (sandbox *winSandbox) AddInterface(SrcName string, DstPrefix string, options ...IfaceOption) error {
// 	return ErrNotImplemented
// }

// func (sandbox *winSandbox) SetGateway(gw net.IP) error {
// 	return ErrNotImplemented
// }

// func (sandbox *winSandbox) SetGatewayIPv6(gw net.IP) error {
// 	return ErrNotImplemented
// }

// func (sandbox *winSandbox) UnsetGateway() error {
// 	return ErrNotImplemented
// }

// func (sandbox *winSandbox) UnsetGatewayIPv6() error {
// 	return ErrNotImplemented
// }

// func (sandbox *winSandbox) AddLoopbackAliasIP(ip *net.IPNet) error {
// 	return ErrNotImplemented
// }

// func (sandbox *winSandbox) RemoveLoopbackAliasIP(ip *net.IPNet) error {
// 	return ErrNotImplemented
// }

// func (sandbox *winSandbox) AddStaticRoute(*types.StaticRoute) error {
// 	return ErrNotImplemented
// }

// func (sandbox *winSandbox) RemoveStaticRoute(*types.StaticRoute) error {
// 	return ErrNotImplemented
// }

// func (sandbox *winSandbox) AddNeighbor(dstIP net.IP, dstMac net.HardwareAddr, force bool, option ...NeighOption) error {
// 	return ErrNotImplemented
// }

// func (sandbox *winSandbox) DeleteNeighbor(dstIP net.IP, dstMac net.HardwareAddr, osDelete bool) error {
// 	return ErrNotImplemented
// }

// func (sandbox *winSandbox) NeighborOptions() NeighborOptionSetter {
// 	return nil
// }

// func (sandbox *winSandbox) InterfaceOptions() IfaceOptionSetter {
// 	return nil
// }

// func (sandbox *winSandbox) InvokeFunc(func()) error {
// 	return ErrNotImplemented
// }

// func (sandbox *winSandbox) Info() Info

// func (sandbox *winSandbox) Destroy() error {
// 	if sandbox.compartment != nil {
// 		compartment.Delete()
// 		compartment = nil
// 	}
// 	return nil
// }

// func (sandbox *winSandbox) Restore(ifsopt map[string][]IfaceOption, routes []*types.StaticRoute, gw net.IP, gw6 net.IP) error {
// 	return ErrNotImplemented
// }

// GenerateKey generates a sandbox key based on the passed
// container id.
func GenerateKey(containerID string) string {
	return containerID
}

// NewSandbox provides a new sandbox instance created in an os specific way
// provided a key which uniquely identifies the sandbox
func NewSandbox(key string, osCreate, isRestore bool) (Sandbox, error) {
	// var compartment *hcsshim.Compartment
	// compartments, err = HNSListComparmentRequest()
	// for _, tp := range compartments {
	// 	if tp.Name == containerID {
	// 		compartment = tp
	// 	}
	// }

	// if compartment == nil {
	// 	compartment = &hcsshim.Compartment{
	// 		Name: Key,
	// 	}
	// 	compartment, err := compartment.Create()

	// 	if err == nil {
	// 		return winSandbox{
	// 			key:         key,
	// 			compartment: compartment,
	// 		}, nil
	// 	}
	// }
	return nil, nil
}

func GetSandboxForExternalKey(path string, key string) (Sandbox, error) {
	return nil, nil
}

// GC triggers garbage collection of namespace path right away
// and waits for it.
func GC() {
}

// InitOSContext initializes OS context while configuring network resources
func InitOSContext() func() {
	return func() {}
}

// SetupTestOSContext sets up a separate test  OS context in which tests will be executed.
func SetupTestOSContext(t *testing.T) func() {
	return func() {}
}

// SetBasePath sets the base url prefix for the ns path
func SetBasePath(path string) {
}
