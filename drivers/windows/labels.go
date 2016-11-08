package windows

const (
	// NetworkName label for bridge driver
	NetworkName = "com.docker.network.windowsshim.networkname"

	// HNSID of the discovered network
	HNSID = "com.docker.network.windowsshim.hnsid"

	// RoutingDomain of the network
	RoutingDomain = "com.docker.network.windowsshim.routingdomain"

	// Interface of the network
	Interface = "com.docker.network.windowsshim.interface"

	// QosPolicies of the endpoint
	QosPolicies = "com.docker.endpoint.windowsshim.qospolicies"

	// VLAN of the network
	VLAN = "com.docker.network.windowsshim.vlanid"

	// VSID of the network
	VSID = "com.docker.network.windowsshim.vsid"

	// DNSSuffix of the network
	DNSSuffix = "com.docker.network.windowsshim.dnssuffix"

	// DNSServers of the network
	DNSServers = "com.docker.network.windowsshim.dnsservers"

	// SourceMac of the network
	SourceMac = "com.docker.network.windowsshim.sourcemac"

	// EnableICC label
	EnableICC = "com.docker.network.bridge.enable_icc"

	// DisableDNS label
	DisableDNS = "com.docker.network.windowsshim.disable_dns"
)
