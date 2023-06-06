package natpmp

import (
	"log"
	"testing"
)

func TestInit(t *testing.T) {
	log.Println(gatewayIP.String())
	result, err := client.GetExternalAddress()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(result.ExternalIPAddress)
}
