package main

import (
	"flag"
	"fmt"

	"github.com/morgadow/gopcan/pcan"
)

func example_cli() {

	channel := flag.String("channel", "PCAN_USBBUS1", "The communication channel, eg. 'PCAN_USBBUS1'")
	baudrate := flag.Int("baudrate", 500000, "The baud rate for communication, eg. '500000'")
	msgID := flag.Int("msg_id", 0x100, "The message ID, eg. '0x1252' (without the 0x)")
	msgData := flag.String("msg_data", "[0, 1, 2, 3, 4, 5, 6, 7, 8]", "The message data as a byte array, eg. [12, 32, 73, 92]. This has an valid default.")
	isExtended := flag.Bool("extended", false, "Whether the message is extended")

	flag.Parse()

	// convert data
	handle := StringToChannel(*channel)
	baud := StringToBaud(*baudrate)
	if handle == nil || baud == nil {
		fmt.Printf("Skipping CLI calls as no valid data given")
		return
	}

	// Convert messageData to a byte array
	byteArray := []byte(*msgData)
	data := [pcan.LENGTH_DATA_CAN_MESSAGE]byte{}
	copy(data[:], byteArray)
	dlc := len(byteArray)

	// Output the parsed values
	fmt.Printf("Parsed CLI data:\n")
	fmt.Printf("\tChannel: %s\n", *channel)
	fmt.Printf("\tBaudrate: %d\n", *baudrate)
	fmt.Printf("\tMessage ID: %d\n", *msgID)
	fmt.Printf("\tMessage Data: %v\n", data)
	fmt.Printf("\tMessage DLC: %v\n", dlc)
	fmt.Printf("\tIs Extended: %t\n", *isExtended)

	// call the api files
	status, bus, err := pcan.InitializeBasic(*handle, *baud)
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while creating PCAN bus: Status: %X, Error: %v\n", status, err)
		return
	}

	// send the message
	msg := pcan.TPCANMsg{ID: pcan.TPCANMsgID(*msgID), DLC: uint8(dlc), Data: data}
	status, err = bus.Write(&msg)
	if status != pcan.PCAN_ERROR_OK || err != nil {
		fmt.Printf("Error while sending message: Status: %X, Error: %v\n", status, err)
		return
	}

	// unitialize handle
	bus.Uninitialize() // returns error but still works, dont know why
}

func StringToChannel(channel string) *pcan.TPCANHandle {
	var handle pcan.TPCANHandle

	switch channel {
	case "PCAN_USBBUS1":
		handle = pcan.PCAN_USBBUS1
	case "PCAN_USBBUS2":
		handle = pcan.PCAN_USBBUS2
	case "PCAN_USBBUS3":
		handle = pcan.PCAN_USBBUS3
	case "PCAN_USBBUS4":
		handle = pcan.PCAN_USBBUS4
	case "PCAN_USBBUS5":
		handle = pcan.PCAN_USBBUS5
	case "PCAN_USBBUS6":
		handle = pcan.PCAN_USBBUS6
	default:
		return nil
	}

	return &handle
}

func StringToBaud(baudrate int) *pcan.TPCANBaudrate {
	var baud pcan.TPCANBaudrate

	switch baudrate {
	case 125000:
		baud = pcan.PCAN_BAUD_125K
	case 250000:
		baud = pcan.PCAN_BAUD_250K
	case 500000:
		baud = pcan.PCAN_BAUD_500K
	case 1000000:
		baud = pcan.PCAN_BAUD_1M
	default:
		return nil
	}

	return &baud
}
