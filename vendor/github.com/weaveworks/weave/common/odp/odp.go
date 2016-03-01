package odp

import (
	"fmt"
	"net"
	"syscall"

	"github.com/dilgerma/go-odp/odp"
)

// ODP admin functionality

func CreateDatapath(dpname string) (err error, supported bool) {
	dpif, err := odp.NewDpif()
	if err != nil {
		if odp.IsKernelLacksODPError(err) {
			return nil, false
		}

		return err, true
	}

	defer dpif.Close()

	dp, err := dpif.CreateDatapath(dpname)
	if err != nil && !odp.IsDatapathNameAlreadyExistsError(err) {
		return err, true
	}

	// Pick an ephemeral port number to use in probing for vxlan
	// support.
	udpconn, err := net.ListenUDP("udp4", nil)
	if err != nil {
		return err, true
	}

	// we leave the UDP socket open, so creating a vxlan vport on
	// the same port number should fail.  But that's fine: It's
	// still sufficient to probe for support.
	portno := uint16(udpconn.LocalAddr().(*net.UDPAddr).Port)
	vpid, err := dp.CreateVport(odp.NewVxlanVportSpec(
		fmt.Sprintf("vxlan-%d", portno), portno))
	if nlerr, ok := err.(odp.NetlinkError); ok {
		if syscall.Errno(nlerr) == syscall.EAFNOSUPPORT {
			dp.Delete()
			return fmt.Errorf("kernel does not have Open vSwitch VXLAN support"), false
		}
	}

	if err == nil {
		dp.DeleteVport(vpid)
	}

	udpconn.Close()
	return nil, true
}

func DeleteDatapath(dpname string) error {
	dpif, err := odp.NewDpif()
	if err != nil {
		return err
	}

	defer dpif.Close()

	dp, err := dpif.LookupDatapath(dpname)
	if err != nil {
		if odp.IsNoSuchDatapathError(err) {
			return nil
		}

		return err
	}

	return dp.Delete()
}

func AddDatapathInterface(dpname string, ifname string) error {
	dpif, err := odp.NewDpif()
	if err != nil {
		return err
	}

	defer dpif.Close()

	dp, err := dpif.LookupDatapath(dpname)
	if err != nil {
		return err
	}

	_, err = dp.CreateVport(odp.NewNetdevVportSpec(ifname))
	return err
}
