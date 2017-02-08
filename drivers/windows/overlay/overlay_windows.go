package overlay

//go:generate protoc -I.:../../Godeps/_workspace/src/github.com/gogo/protobuf  --gogo_out=import_path=github.com/docker/libnetwork/drivers/overlay,Mgogoproto/gogo.proto=github.com/gogo/protobuf/gogoproto:. overlay.proto

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"

	"github.com/Microsoft/hcsshim"
	"github.com/Sirupsen/logrus"
	"github.com/docker/libnetwork/datastore"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/driverapi"
	"github.com/docker/libnetwork/netlabel"
	"github.com/docker/libnetwork/types"
	"github.com/hashicorp/serf/serf"
)

const (
	networkType  = "overlay"
	vethPrefix   = "veth"
	vethLen      = 7
	secureOption = "encrypted"
)

type driver struct {
	eventCh          chan serf.Event
	notifyCh         chan ovNotify
	exitCh           chan chan struct{}
	bindAddress      string
	advertiseAddress string
	neighIP          string
	config           map[string]interface{}
	serfInstance     *serf.Serf
	networks         networkTable
	store            datastore.DataStore
	localStore       datastore.DataStore
	once             sync.Once
	joinOnce         sync.Once
	sync.Mutex
}

// Init registers a new instance of overlay driver
func Init(dc driverapi.DriverCallback, config map[string]interface{}) error {
	c := driverapi.Capability{
		DataScope: datastore.GlobalScope,
	}

	d := &driver{
		networks: networkTable{},
		config:   config,
	}

	if data, ok := config[netlabel.GlobalKVClient]; ok {
		var err error
		dsc, ok := data.(discoverapi.DatastoreConfigData)
		if !ok {
			return types.InternalErrorf("incorrect data in datastore configuration: %v", data)
		}
		d.store, err = datastore.NewDataStoreFromConfig(dsc)
		if err != nil {
			return types.InternalErrorf("failed to initialize data store: %v", err)
		}
	}

	if data, ok := config[netlabel.LocalKVClient]; ok {
		var err error
		dsc, ok := data.(discoverapi.DatastoreConfigData)
		if !ok {
			return types.InternalErrorf("incorrect data in datastore configuration: %v", data)
		}
		d.localStore, err = datastore.NewDataStoreFromConfig(dsc)
		if err != nil {
			return types.InternalErrorf("failed to initialize local data store: %v", err)
		}
	}

	d.restoreHNSNetworks()

	return dc.RegisterDriver(networkType, d, c)
}

func (d *driver) restoreHNSNetworks() error {
	logrus.Infof("Restoring existing overlay networks from HNS into docker")

	hnsresponse, err := hcsshim.HNSListNetworkRequest("GET", "", "")
	if err != nil {
		return err
	}

	for _, v := range hnsresponse {
		if v.Type != networkType {
			continue
		}

		logrus.Infof("Restoring overlay network: %s", v.Name)
		n := d.convertToOverlayNetwork(&v)
		d.addNetwork(n)

		n.restoreNetworkEndpoints()
	}

	return nil
}

func (d *driver) convertToOverlayNetwork(v *hcsshim.HNSNetwork) *network {
	n := &network{
		id:              v.Name,
		hnsId:           v.Id,
		driver:          d,
		endpoints:       endpointTable{},
		subnets:         []*subnet{},
		providerAddress: v.ManagementIP,
	}

	for _, hnsSubnet := range v.Subnets {
		vsidPolicy := &hcsshim.VsidPolicy{}
		for _, policy := range hnsSubnet.Policies {
			if err := json.Unmarshal([]byte(policy), &vsidPolicy); err == nil && vsidPolicy.Type == "VSID" {
				break
			}
		}

		gwIP := net.ParseIP(hnsSubnet.GatewayAddress)
		localsubnet := &subnet{
			vni:  uint32(vsidPolicy.VSID),
			gwIP: &gwIP,
		}

		_, subnetIP, err := net.ParseCIDR(hnsSubnet.AddressPrefix)

		if err != nil {
			logrus.Errorf("Error parsing subnet address %s ", hnsSubnet.AddressPrefix)
			continue
		}

		localsubnet.subnetIP = subnetIP

		n.subnets = append(n.subnets, localsubnet)
	}

	return n
}

func (n *network) restoreNetworkEndpoints() error {
	logrus.Infof("Restoring endpoints for overlay network: %s", n.id)

	hnsresponse, err := hcsshim.HNSListEndpointRequest("GET", "", "")
	if err != nil {
		return err
	}

	for _, endpoint := range hnsresponse {
		if endpoint.VirtualNetwork != n.hnsId {
			continue
		}

		ep := n.convertToOverlayEndpoint(&endpoint)

		if ep != nil {
			logrus.Debugf("Restored endpoint:%s Remote:%t", ep.id, ep.remote)
			n.addEndpoint(ep)
		}
	}

	return nil
}

func (n *network) convertToOverlayEndpoint(v *hcsshim.HNSEndpoint) *endpoint {
	ep := &endpoint{
		id:        v.Name,
		profileId: v.Id,
		nid:       n.id,
		remote:    v.IsRemoteEndpoint,
	}

	mac, err := net.ParseMAC(v.MacAddress)

	if err != nil {
		return nil
	}

	ep.mac = mac
	ep.addr = &net.IPNet{
		IP:   v.IPAddress,
		Mask: net.CIDRMask(32, 32),
	}

	return ep
}

// Fini cleans up the driver resources
func Fini(drv driverapi.Driver) {
	d := drv.(*driver)

	if d.exitCh != nil {
		waitCh := make(chan struct{})

		d.exitCh <- waitCh

		<-waitCh
	}
}

func (d *driver) configure() error {
	if d.store == nil {
		return nil
	}

	return nil
}

func (d *driver) Type() string {
	return networkType
}

func (d *driver) IsBuiltIn() bool {
	return true
}

func validateSelf(node string) error {
	advIP := net.ParseIP(node)
	if advIP == nil {
		return fmt.Errorf("invalid self address (%s)", node)
	}

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return fmt.Errorf("Unable to get interface addresses %v", err)
	}
	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err == nil && ip.Equal(advIP) {
			return nil
		}
	}
	return fmt.Errorf("Multi-Host overlay networking requires cluster-advertise(%s) to be configured with a local ip-address that is reachable within the cluster", advIP.String())
}

func (d *driver) nodeJoin(advertiseAddress, bindAddress string, self bool) {
	if self && !d.isSerfAlive() {
		if err := validateSelf(advertiseAddress); err != nil {
			logrus.Errorf("%s", err.Error())
		}

		d.Lock()
		d.advertiseAddress = advertiseAddress
		d.bindAddress = bindAddress
		d.Unlock()

		// If there is no cluster store there is no need to start serf.
		if d.store != nil {
			err := d.serfInit()
			if err != nil {
				logrus.Errorf("initializing serf instance failed: %v", err)
				return
			}
		}
	}

	d.Lock()
	if !self {
		d.neighIP = advertiseAddress
	}
	neighIP := d.neighIP
	d.Unlock()

	if d.serfInstance != nil && neighIP != "" {
		var err error
		d.joinOnce.Do(func() {
			err = d.serfJoin(neighIP)
			if err == nil {
				d.pushLocalDb()
			}
		})
		if err != nil {
			logrus.Errorf("joining serf neighbor %s failed: %v", advertiseAddress, err)
			d.Lock()
			d.joinOnce = sync.Once{}
			d.Unlock()
			return
		}
	}
}

func (d *driver) pushLocalEndpointEvent(action, nid, eid string) {
	n := d.network(nid)
	if n == nil {
		logrus.Debugf("Error pushing local endpoint event for network %s", nid)
		return
	}
	ep := n.endpoint(eid)
	if ep == nil {
		logrus.Debugf("Error pushing local endpoint event for ep %s / %s", nid, eid)
		return
	}

	if !d.isSerfAlive() {
		return
	}
	d.notifyCh <- ovNotify{
		action: action,
		nw:     n,
		ep:     ep,
	}
}

// DiscoverNew is a notification for a new discovery event, such as a new node joining a cluster
func (d *driver) DiscoverNew(dType discoverapi.DiscoveryType, data interface{}) error {

	var err error
	switch dType {
	case discoverapi.NodeDiscovery:
		nodeData, ok := data.(discoverapi.NodeDiscoveryData)
		if !ok || nodeData.Address == "" {
			return fmt.Errorf("invalid discovery data")
		}
		d.nodeJoin(nodeData.Address, nodeData.BindAddress, nodeData.Self)
	case discoverapi.DatastoreConfig:
		if d.store != nil {
			return types.ForbiddenErrorf("cannot accept datastore configuration: Overlay driver has a datastore configured already")
		}
		dsc, ok := data.(discoverapi.DatastoreConfigData)
		if !ok {
			return types.InternalErrorf("incorrect data in datastore configuration: %v", data)
		}
		d.store, err = datastore.NewDataStoreFromConfig(dsc)
		if err != nil {
			return types.InternalErrorf("failed to initialize data store: %v", err)
		}
	default:
	}
	return nil
}

// DiscoverDelete is a notification for a discovery delete event, such as a node leaving a cluster
func (d *driver) DiscoverDelete(dType discoverapi.DiscoveryType, data interface{}) error {
	return nil
}
