package tapi

import (
	"fmt"
	"log"

	"github.com/chbmuc/lirc"
)

// NewTadiranAPI initializes a new lircd api for tadiran.
func NewTadiranAPI(name, socket string) (tapi *TadiranAPI, err error) {
	ir, err := lirc.Init(socket)
	if err != nil {
		log.Println("Could not initialize new lircd socket at", socket)
		return
	}

	tapi = &TadiranAPI{
		lircRouter: ir,
		remoteName: name,
	}
	return
}

// TadiranAPI defines all of the functionality that you can use our API with lirc
type TadiranAPI struct {
	lircRouter         *lirc.Router
	remoteName         string
	currentFanSpeed    fanSpeed
	currentHeatState   heatState
	currentState       state
	currentTemperature int
}

type state bool

const (
	// ON indicates that the AC is on
	ON state = true
	// OFF indicates that the AC is off
	OFF state = false
)

type heatState string

const (
	// HOT Defines the hot state
	HOT heatState = "HOT"
	// COLD Defines the cold state
	COLD heatState = "COLD"
)

type fanSpeed string

const (
	// ECO Defines the echo fan speed
	ECO fanSpeed = "ECO"
	// ONE Defines the 1 fan speed
	ONE fanSpeed = "FAN_1"
	// TWO Defines the 2 fan speed
	TWO fanSpeed = "FAN_2"
	// THREE Defines the 3 fan speed
	THREE fanSpeed = "FAN_3"
	// TURBO Defines the turbo fan speed
	TURBO fanSpeed = "TURBO"
)

// Off turns off the AC
func (tapi *TadiranAPI) Off() error {
	tapi.currentState = OFF
	return tapi.send("POWER_OFF")
}

// SendCommand wraps the structure needed to send a command
func (tapi *TadiranAPI) SendCommand(temperature int, newState heatState, newFanSpeed fanSpeed) error {
	if temperature < 16 || temperature > 32 {
		return fmt.Errorf("Temperature can be only between 16-32 and not %d", temperature)
	}
	tapi.currentFanSpeed = newFanSpeed
	tapi.currentHeatState = newState
	tapi.currentTemperature = temperature
	tapi.currentState = ON
	command := fmt.Sprintf("ON_%s_%s_%d", newFanSpeed, newState, temperature)
	return tapi.send(command)
}

func (tapi *TadiranAPI) send(command string) error {
	fullCommand := fmt.Sprintf("%s %s", tapi.remoteName, command)
	log.Println("Sending command", fullCommand)
	err := tapi.lircRouter.Send(fullCommand) //"TadiranCarrierAC POWER_OFF")
	if err != nil {
		log.Println("Error Sending Command", fullCommand)
		log.Println(err)
	}
	return err
}
