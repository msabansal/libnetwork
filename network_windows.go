// +build windows

package libnetwork

import (
	"runtime"

	"github.com/Microsoft/hcsshim"
	log "github.com/Sirupsen/logrus"
	"github.com/docker/libnetwork/drivers/windows"
)

func executeInCompartment(compartmentID uint32, x func()) {
	runtime.LockOSThread()

	if err := hcsshim.SetCurrentThreadCompartmentId(compartmentID); err != nil {
		log.Error(err)
	}
	defer func() {
		hcsshim.SetCurrentThreadCompartmentId(0)
		runtime.UnlockOSThread()
	}()

	x()
}

func (n *network) startResolver() {
	n.resolverOnce.Do(func() {
		log.Warnf("Launching DNS server for network", n.Name())
		options := n.Info().DriverOptions()
		hnsid := options[windows.HNSID]

		hnsresponse, err := hcsshim.HNSNetworkRequest("GET", hnsid, "")
		if err != nil {
			log.Errorf("Resolver Setup/Start failed for container %s, %q", n.Name(), err)
			return
		}

		if hnsresponse.DnsServerAddress != "" {
			n.resolver = NewResolver(nil, n)
			defer func() {
				if err != nil {
					n.resolver = nil
				}
			}()

			executeInCompartment(hnsresponse.DnsServerCompartment, n.resolver.SetupFunc(hnsresponse.DnsServerAddress, 53))
			if err = n.resolver.Start(); err != nil {
				log.Errorf("Resolver Setup/Start failed for container %s, %q", n.Name(), err)
			}
		}

	})
}
