package learning

import (
	"fmt"
	"net"
	"os"
	"strings"
	"testing"
)

func TestInterfaceAddresses(t *testing.T) {

	if os.Getenv("PRINT-NETWORKS") == "" {
		t.Skip("Set PRINT-NETWORKS to enable this test and printing available interfaces.")
	}

	addresses, err := net.InterfaceAddrs()
	if err != nil {
		t.Error(err)
		return
	}
	for _, a := range addresses {
		fmt.Printf("Network: %s. Address: %s.\n", a.Network(), a.String())
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		t.Error(err)
		return
	}
	ifaceNames := []string{}
	for _, iface := range interfaces {
		ifaceNames = append(ifaceNames, iface.Name)
	}
	fmt.Println("Interfaces:", strings.Join(ifaceNames, ", "))

}
