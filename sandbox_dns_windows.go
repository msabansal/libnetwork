// +build windows

package libnetwork

// Stub implementations for DNS related functions

func (sb *sandbox) startResolver(bool) {
}

func (sb *sandbox) setupResolutionFiles() error {
	return nil
}

func (sb *sandbox) restorePath() {
}

func (sb *sandbox) updateDNS(ipv6Enabled bool) error {
	return nil
}
