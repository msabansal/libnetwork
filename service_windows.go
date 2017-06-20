package libnetwork

import (
	"net"

	"github.com/Sirupsen/logrus"
)

func (n *network) addLBBackend(ip, vip net.IP, fwMark uint32, ingressPorts []*PortConfig) {
	logrus.Debugf("Adding lb backend %v %v portconfig %v", ip, vip, ingressPorts)
	// n.WalkEndpoints(func(e Endpoint) bool {
	// 	ep := e.(*endpoint)
	// 	if sb, ok := ep.getSandbox(); ok {
	// 		if !sb.isEndpointPopulated(ep) {
	// 			return false
	// 		}

	// 		var gwIP net.IP
	// 		if ep := sb.getGatewayEndpoint(); ep != nil {
	// 			gwIP = ep.Iface().Address().IP
	// 		}

	// 		sb.addLBBackend(ip, vip, fwMark, ingressPorts, ep.Iface().Address(), gwIP, n.ingress)
	// 	}

	// 	return false
	// })
}

func (n *network) rmLBBackend(ip, vip net.IP, fwMark uint32, ingressPorts []*PortConfig, rmService bool) {
	logrus.Debugf("Removing lb backend %v %v", ip, vip)
}

func (sb *sandbox) populateLoadbalancers(ep *endpoint) {
	logrus.Debugf("Populating lb for ep %v %v", ep)
}

func arrangeIngressFilterRule() {
}
