package main

import (
	"log"

	"github.com/vladshub/tadiran_api/tapi"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	// Initialize with path to lirc socket
	t, err := tapi.NewTadiranAPI("TadiranCarrierAC", "/var/run/lirc/lircd")
	failOnError(err, "Problem initializing socket")

	err = t.SendCommand(19, tapi.COLD, tapi.TURBO)
	failOnError(err, "Failed sending message")
}
