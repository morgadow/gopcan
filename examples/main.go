package main

import (
	"fmt"

	"github.com/morgadow/gopcan/pcan"
)

func main() {

	// connect to new channel
	status, bus, err := pcan.InitializeBasic(pcan.PCAN_USBBUS1, pcan.PCAN_BAUD_500K)
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while creating PCAN bus: Status: %X, Error: %v\n", status, err)
		return
	}

	// reset hardware
	status, err = bus.Reset()
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while resetting PCAN bus: Status: %X, Error: %v\n", status, err)
		return
	}

	// check hardware status
	status, err = bus.GetStatus()
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while resetting PCAN bus: Status: %X, Error: %v\n", status, err)
		return
	}

	// setup allow several frame types
	status, err = bus.SetAllowEchoFrames(true)
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while setting parameter: Status: %X, Error: %v\n", status, err)
		return
	}
	status, err = bus.SetAllowErrorFrames(true)
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while setting parameter: Status: %X, Error: %v\n", status, err)
		return
	}
	status, err = bus.SetAllowRTRFrames(false)
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while setting parameter: Status: %X, Error: %v\n", status, err)
		return
	}
	status, err = bus.SetAllowStatusFrames(true)
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while setting parameter: Status: %X, Error: %v\n", status, err)
		return
	}

	// send standard message
	txMsg := pcan.TPCANMsg{ID: 0x123, Data: [pcan.LENGTH_DATA_CAN_MESSAGE]byte{1, 2, 3, 4, 5, 6, 7, 8}, DLC: 8, MsgType: pcan.PCAN_MESSAGE_STANDARD}
	status, err = bus.Write(&txMsg)
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while sending message: Status: %X, Error: %v\n", status, err)
		return
	}

	// send extended message
	txMsg = pcan.TPCANMsg{ID: 0x12345, Data: [pcan.LENGTH_DATA_CAN_MESSAGE]byte{1, 2, 3, 4, 5, 6, 7, 8}, DLC: 8, MsgType: pcan.PCAN_MESSAGE_EXTENDED}
	status, err = bus.Write(&txMsg)
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while sending message: Status: %X, Error: %v\n", status, err)
		return
	}

	// read message from bus
	status, rxmsg, txtimestamp, err := bus.Read()
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while reading message: Status: %X, Error: %v\n", status, err)
		return
	}
	if rxmsg != nil {
		fmt.Printf("Received message 0x%X with type %v and data %v at: %v:%v:%v\n", rxmsg.ID, rxmsg.MsgType, rxmsg.Data, txtimestamp.Millis, txtimestamp.MillisOverflow, txtimestamp.Micros)
	} else {
		fmt.Printf("Did not receive a message\n")
	}

	// read message in timeout from bus
	status, rxmsg, txtimestamp, err = bus.ReadWithTimeout(500)
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while reading message: Status: %X, Error: %v\n", status, err)
		return
	}
	if rxmsg != nil {
		fmt.Printf("Received message 0x%X with type %v and data %v at: %v:%v:%v\n", rxmsg.ID, rxmsg.MsgType, rxmsg.Data, txtimestamp.Millis, txtimestamp.MillisOverflow, txtimestamp.Micros)
	} else {
		fmt.Printf("Did not receive a message in timeout of 500ms\n")
	}

}