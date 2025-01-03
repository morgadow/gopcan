package pcan

import (
	"errors"
	"fmt"
	"runtime"
	"syscall"
	"unsafe"
)

/* This file is the windows specific implementation for handling the PCAN driver. */

// PCAN Bus interface
type TPCANBus struct {
	Handle    TPCANHandle
	Baudrate  TPCANBaudrate // only set if not a FD channel
	HWType    TPCANType     // only for non plug´n´play devices and currently not used
	IOPort    uint32        // only for non plug´n´play devices and currently not used
	Interrupt uint16        // only for non plug´n´play devices and currently not used
	recvEvent syscall.Handle
}

// PCAN Bus interface for CANFD channels
type TPCANBusFD struct {
	Handle    TPCANHandle
	BitrateFD TPCANBitrateFD // only set if a FD channel
	// TODO fill with FD parameter and other necessary stuff
}

// api procedures
var (
	pcanAPIHandle         *syscall.DLL  = nil // procedure handle for PCAN driver
	pHandleInitialize     *syscall.Proc = nil
	pHandleInitializeFD   *syscall.Proc = nil
	pHandleUninitialize   *syscall.Proc = nil
	pHandleReset          *syscall.Proc = nil
	pHandleGetStatus      *syscall.Proc = nil
	pHandleRead           *syscall.Proc = nil
	pHandleReadFD         *syscall.Proc = nil
	pHandleWrite          *syscall.Proc = nil
	pHandleWriteFD        *syscall.Proc = nil
	pHandleFilterMessages *syscall.Proc = nil
	pHandleGetValue       *syscall.Proc = nil
	pHandleSetValue       *syscall.Proc = nil
	pHandleGetErrorText   *syscall.Proc = nil
	pHandleLookUpChannel  *syscall.Proc = nil

	apiLoaded bool = false // indicates if the api was loaded already, set by LoadApi() and unset by UnloadApi()
	hasEvents bool = false
)

// Loads PCAN API (.ddl) file
func LoadAPI() error {
	var err error = nil

	if apiLoaded {
		return nil
	}

	// evaluate operating system and architecture and select driver file
	if runtime.GOOS != "windows" {
		return fmt.Errorf("invalid operating system. change compile option to match %v", runtime.GOOS)
	}

	pcanAPIHandle, err = syscall.LoadDLL("PCANBasic.dll")
	if err != nil || pcanAPIHandle == nil {
		return err
	}

	pHandleInitialize, _ = pcanAPIHandle.FindProc("CAN_Initialize")
	pHandleInitializeFD, _ = pcanAPIHandle.FindProc("CAN_InitializeFD")
	pHandleUninitialize, _ = pcanAPIHandle.FindProc("CAN_Uninitialize")
	pHandleReset, _ = pcanAPIHandle.FindProc("CAN_Reset")
	pHandleGetStatus, _ = pcanAPIHandle.FindProc("CAN_GetStatus")
	pHandleRead, _ = pcanAPIHandle.FindProc("CAN_Read")
	pHandleReadFD, _ = pcanAPIHandle.FindProc("CAN_ReadFD")
	pHandleWrite, _ = pcanAPIHandle.FindProc("CAN_Write")
	pHandleWriteFD, _ = pcanAPIHandle.FindProc("CAN_WriteFD")
	pHandleFilterMessages, _ = pcanAPIHandle.FindProc("CAN_FilterMessages")
	pHandleGetValue, _ = pcanAPIHandle.FindProc("CAN_GetValue")
	pHandleSetValue, _ = pcanAPIHandle.FindProc("CAN_SetValue")
	pHandleGetErrorText, _ = pcanAPIHandle.FindProc("CAN_GetErrorText")
	pHandleLookUpChannel, _ = pcanAPIHandle.FindProc("CAN_LookUpChannel")

	apiLoaded = pHandleInitialize != nil && pHandleInitializeFD != nil && pHandleReset != nil && pHandleGetStatus != nil &&
		pHandleRead != nil && pHandleReadFD != nil && pHandleWrite != nil && pHandleWriteFD != nil && pHandleFilterMessages != nil && pHandleGetValue != nil &&
		pHandleSetValue != nil && pHandleGetErrorText != nil && pHandleLookUpChannel != nil && pHandleUninitialize != nil

	if !apiLoaded {
		return errors.New("could not load pointers to pcan functions")
	}
	return nil
}

// Unloads PCAN API (.ddl) file
func UnloadAPI() error {

	// reset pointers
	pHandleInitialize = nil
	pHandleInitializeFD = nil
	pHandleUninitialize = nil
	pHandleReset = nil
	pHandleGetStatus = nil
	pHandleRead = nil
	pHandleReadFD = nil
	pHandleWrite = nil
	pHandleWriteFD = nil
	pHandleFilterMessages = nil
	pHandleGetValue = nil
	pHandleSetValue = nil
	pHandleGetErrorText = nil
	pHandleLookUpChannel = nil
	pHandleUninitialize = nil
	apiLoaded = false

	err := pcanAPIHandle.Release()
	return err
}

// API call to iInitializes a basic plugNplay PCAN Channel
// Channel: The handle of a PCAN Channel
// baudRate: The speed for the communication (BTR0BTR1 code)
func APIInitializeBasic(handle TPCANHandle, baudRate TPCANBaudrate) (TPCANStatus, error) {
	r, _, errno := pHandleInitialize.Call(uintptr(handle), uintptr(baudRate))
	return TPCANStatus(r), syscallErr(errno)
}

// API call to initializes a advanced PCAN Channel
// Channel: The handle of a PCAN Channel
// baudRate: The speed for the communication (BTR0BTR1 code)
// hwType: Non-PnP: The type of hardware and operation mode
// ioPort: Non-PnP: The I/O address for the parallel port
// interrupt: Non-PnP: Interrupt number of the parallel port
func APIInitialize(handle TPCANHandle, baudRate TPCANBaudrate, hwType TPCANType, ioPort uint32, interrupt uint16) (TPCANStatus, error) {
	r, _, errno := pHandleInitialize.Call(uintptr(handle), uintptr(baudRate), uintptr(hwType), uintptr(ioPort), uintptr(interrupt))
	return TPCANStatus(r), syscallErr(errno)
}

// API call to initializes a FD capable PCAN Channel
// Channel: The handle of a PCAN Channel
// bitRateFD: The speed for the communication (FD bit rate string)
// Note:
// Bit rate string must follow the following construction rules:
//   - parameter and values must be separated by '='
//   - Couples of Parameter/value must be separated by ','
//   - Following Parameter must be filled out: f_clock, data_brp, data_sjw, data_tseg1, data_tseg2,
//     nom_brp, nom_sjw, nom_tseg1, nom_tseg2.
//   - Following Parameters are optional (not used yet): data_ssp_offset, nom_sam
//   - Example: f_clock=80000000,nom_brp=10,nom_tseg1=5,nom_tseg2=2,nom_sjw=1,data_brp=4,data_tseg1=7,data_tseg2=2,data_sjw=1
func APIInitializeFD(handle TPCANHandle, bitRateFD TPCANBitrateFD) (TPCANStatus, error) {
	r, _, errno := pHandleInitializeFD.Call(uintptr(handle), uintptr(unsafe.Pointer(&bitRateFD)))
	return TPCANStatus(r), syscallErr(errno)
}

// API call to uninitializes PCAN Channels initialized by CAN_Initialize
func APIUninitialize(handle TPCANHandle) (TPCANStatus, error) {
	r, _, errno := pHandleUninitialize.Call(uintptr(handle))
	return TPCANStatus(r), syscallErr(errno)
}

// API call to reset the receive and transmit queues of the PCAN Channel
func APIReset(handle TPCANHandle) (TPCANStatus, error) {
	r, _, errno := pHandleReset.Call(uintptr(handle))
	return TPCANStatus(r), syscallErr(errno)
}

// API call to get the current status of a PCAN Channel
func APIGetStatus(handle TPCANHandle) (TPCANStatus, error) {
	r, _, errno := pHandleGetStatus.Call(uintptr(handle))
	return TPCANStatus(r), syscallErr(errno)
}

// API call to read a CAN message from the receive queue of a PCAN Channel
// Note: Does return nil if receive buffer is empty
func APIRead(handle TPCANHandle) (TPCANStatus, TPCANMsg, TPCANTimestamp, error) {
	var msg TPCANMsg
	var timestamp TPCANTimestamp

	r, _, errno := pHandleRead.Call(uintptr(handle), uintptr(unsafe.Pointer(&msg)), uintptr(unsafe.Pointer(&timestamp)))
	return TPCANStatus(r), msg, timestamp, syscallErr(errno)
}

// API call to read a CAN message from the receive queue of a FD capable PCAN Channel
func APIReadFD(handle TPCANHandle) (TPCANStatus, TPCANMsgFD, TPCANTimestampFD, error) {
	var msg TPCANMsgFD
	var timestamp TPCANTimestampFD

	r, _, errno := pHandleReadFD.Call(uintptr(handle), uintptr(unsafe.Pointer(&msg)), uintptr(unsafe.Pointer(&timestamp)))
	return TPCANStatus(r), msg, timestamp, syscallErr(errno)
}

// API call to transmits a CAN message
// msg: A Message struct with the message to be sent
func APIWrite(handle TPCANHandle, msg *TPCANMsg) (TPCANStatus, error) {
	r, _, errno := pHandleWrite.Call(uintptr(handle), uintptr(unsafe.Pointer(msg)))
	return TPCANStatus(r), syscallErr(errno)
}

// API call to transmit a CAN message over a FD capable PCAN Channel
// msgFD A MessageFD struct with the message to be sent
func APIWriteFD(handle TPCANHandle, msg *TPCANMsgFD) (TPCANStatus, error) {
	r, _, errno := pHandleWriteFD.Call(uintptr(handle), uintptr(unsafe.Pointer(msg)))
	return TPCANStatus(r), syscallErr(errno)
}

// API call to retrieve a PCAN Channel value
// param: The TPCANParameter parameter to get
// Note: Parameters can be present or not according with the kind
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func APIGetValue(handle TPCANHandle, param TPCANParameter, buffer unsafe.Pointer, bufferSize uint32) (TPCANStatus, error) {
	r, _, errno := pHandleGetValue.Call(uintptr(handle), uintptr(param), uintptr(buffer), uintptr(bufferSize))
	return TPCANStatus(r), syscallErr(errno)
}

// API call to configure a PCAN Channel value.
// handle: The handle of a PCAN Channel
// param: The TPCANParameter parameter to set
// value: Value of parameter
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func APISetValue(handle TPCANHandle, param TPCANParameter, buffer unsafe.Pointer, bufferSize uint32) (TPCANStatus, error) {
	r, _, errno := pHandleSetValue.Call(uintptr(handle), uintptr(param), uintptr(buffer), uintptr(bufferSize))
	return TPCANStatus(r), syscallErr(errno)
}

// API call to configure the reception filter
// fromID: The lowest CAN ID to be received
// toID: The highest CAN ID to be received
// mode: Message type, Standard (11-bit identifier) or Extended (29-bit identifier)
func APISetFilter(handle TPCANHandle, fromID TPCANMsgID, toID TPCANMsgID, mode TPCANMode) (TPCANStatus, error) {
	r, _, errno := pHandleFilterMessages.Call(uintptr(handle), uintptr(fromID), uintptr(toID), uintptr(mode))
	return TPCANStatus(r), syscallErr(errno)
}

// API call to return a descriptive text of a given TPCANStatus error code, in any desired language
// err: A TPCANStatus error code
// language: Indicates a 'Primary language ID'
func APIGetErrorText(status TPCANStatus, language TPCANLanguage) (TPCANStatus, [MAX_LENGHT_STRING_BUFFER]byte, error) {
	var buffer [MAX_LENGHT_STRING_BUFFER]byte

	r, _, errno := pHandleGetErrorText.Call(uintptr(status), uintptr(language), uintptr(unsafe.Pointer(&buffer)))
	return TPCANStatus(r), buffer, syscallErr(errno)
}

// API call to find a PCAN-Basic Channel that matches with the given parameters
// parameters: A comma separated string contained pairs of parameter-name/value to be matched within a PCAN-Basic Channel
// foundChannels: Buffer for returning the PCAN-Basic Channel when found
func APILookUpChannel(deviceType string, deviceID string, controllerNumber string, ipAdress string) (TPCANStatus, TPCANHandle, error) {

	var sParameters string = ""
	var foundChannel TPCANHandle

	// merge search parameter
	if deviceType != "" {
		sParameters += string(LOOKUP_DEVICE_TYPE) + "=" + deviceType
	}

	if deviceID != "" {
		if sParameters != "" {
			sParameters += ", "
		}
		sParameters += string(LOOKUP_DEVICE_ID) + "=" + deviceID
	}
	if controllerNumber != "" {
		if sParameters != "" {
			sParameters += ", "
		}
		sParameters += string(LOOKUP_CONTROLLER_NUMBER) + "=" + controllerNumber
	}
	if ipAdress != "" {
		if sParameters != "" {
			sParameters += ", "
		}
		sParameters += string(LOOKUP_IP_ADDRESS) + "=" + ipAdress
	}

	r, _, errno := pHandleLookUpChannel.Call(uintptr(unsafe.Pointer(&sParameters)), uintptr(unsafe.Pointer(&foundChannel)))
	return TPCANStatus(r), foundChannel, syscallErr(errno)
}

// helper function to handle syscall return value
func syscallErr(err error) error {
	if err != nil {
		errno := err.(syscall.Errno)
		if errno != 0 {

			// suppress this warning as this is set by PCAN api
			if errno == syscall.ERROR_INSUFFICIENT_BUFFER {
				fmt.Printf("pcan api warning: %v\n", errno)
				return nil
			}

			return errors.New(errno.Error())
		}
	}
	return nil
}
