package pcan

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"syscall"
	"time"
	"unsafe"
)

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
		return errors.New(fmt.Sprintf("invalid operating system. change compile option to match %v", runtime.GOOS))
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

// Initializes a basic plugNplay PCAN Channel
// Channel: The handle of a PCAN Channel
// baudRate: The speed for the communication (BTR0BTR1 code)
func InitializeBasic(handle TPCANHandle, baudRate TPCANBaudrate) (TPCANStatus, *TPCANBus, error) {
	LoadAPI()

	r, _, errno := pHandleInitialize.Call(uintptr(handle), uintptr(baudRate))
	err := syscallErr(errno)

	if TPCANStatus(r) != PCAN_ERROR_OK || err != nil {
		return TPCANStatus(r), nil, err
	}

	bus := TPCANBus{
		Handle:    handle,
		Baudrate:  baudRate,
		HWType:    PCAN_DEFAULT_HW_TYPE,
		IOPort:    PCAN_DEFAULT_IO_PORT,
		Interrupt: PCAN_DEFAULT_INTERRUPT}

	bus.initializeRecvEvent()

	return TPCANStatus(r), &bus, err
}

// Initializes a advanced PCAN Channel
// Channel: The handle of a PCAN Channel
// baudRate: The speed for the communication (BTR0BTR1 code)
// hwType: Non-PnP: The type of hardware and operation mode
// ioPort: Non-PnP: The I/O address for the parallel port
// interrupt: Non-PnP: Interrupt number of the parallel port
func Initialize(handle TPCANHandle, baudRate TPCANBaudrate, hwType TPCANType, ioPort uint32, interrupt uint16) (TPCANStatus, *TPCANBus, error) {
	LoadAPI()

	r, _, errno := pHandleInitialize.Call(uintptr(handle), uintptr(baudRate), uintptr(hwType), uintptr(ioPort), uintptr(interrupt))
	err := syscallErr(errno)

	if TPCANStatus(r) != PCAN_ERROR_OK || err != nil {
		return TPCANStatus(r), nil, err
	}

	bus := TPCANBus{
		Handle:    handle,
		Baudrate:  baudRate,
		HWType:    hwType,
		IOPort:    ioPort,
		Interrupt: interrupt}

	bus.initializeRecvEvent()

	return TPCANStatus(r), &bus, err
}

// Initializes a FD capable PCAN Channel
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
func InitializeFD(channel TPCANHandle, bitRateFD TPCANBitrateFD) (TPCANStatus, *TPCANBusFD, error) {
	LoadAPI()
	return PCAN_ERROR_UNKNOWN, nil, errors.New("Not implemented")
	// ret, _, errno := pHandleInitializeFD.Call(uintptr(channel), uintptr(unsafe.Pointer(&bitRateFD)))
	// return TPCANStatus(ret), syscallErr(errno)
}

// Uninitializes PCAN Channels initialized by CAN_Initialize
func (p *TPCANBus) Uninitialize() (TPCANStatus, error) {
	ret, _, errno := pHandleUninitialize.Call(uintptr(p.Handle))
	return TPCANStatus(ret), syscallErr(errno)
}

// Resets the receive and transmit queues of the PCAN Channel
func (p *TPCANBus) Reset() (TPCANStatus, error) {
	ret, _, errno := pHandleReset.Call(uintptr(p.Handle))
	return TPCANStatus(ret), syscallErr(errno)
}

// Gets the current status of a PCAN Channel
func (p *TPCANBus) GetStatus() (TPCANStatus, error) {
	ret, _, errno := pHandleGetStatus.Call(uintptr(p.Handle))
	return TPCANStatus(ret), syscallErr(errno)
}

// Reads a CAN message from the receive queue of a PCAN Channel
// Note: Does return nil if receive buffer is empty
func (p *TPCANBus) Read() (TPCANStatus, *TPCANMsg, *TPCANTimestamp, error) {
	var msg TPCANMsg
	var timeStamp TPCANTimestamp

	ret, _, errno := pHandleRead.Call(uintptr(p.Handle), uintptr(unsafe.Pointer(&msg)), uintptr(unsafe.Pointer(&timeStamp)))

	if TPCANStatus(ret) == PCAN_ERROR_QRCVEMPTY {
		return TPCANStatus(ret), nil, nil, syscallErr(errno)
	} else {
		return TPCANStatus(ret), &msg, &timeStamp, syscallErr(errno)
	}
}

// Reads a CAN message from the receive queue of a PCAN Channel with an timeout and only returns a valid messsage
// Note: Does return nil if receive buffer is empty or no message is read during timeout
// timeout: Timeout for receiving message from CAN bus in milliseconds (if set below zero, no timeout is set)
func (p *TPCANBus) ReadWithTimeout(timeout int) (TPCANStatus, *TPCANMsg, *TPCANTimestamp, error) {

	var ret = PCAN_ERROR_UNKNOWN
	var msg *TPCANMsg = nil
	var timeStamp *TPCANTimestamp = nil
	var err error = nil

	// timeout handling: a negative timeout sets timeout to infinity
	if timeout < 0 {
		timeout = syscall.INFINITE
	}
	var timeoutU32 = uint32(timeout)
	startTime := time.Now().UnixNano() / int64(time.Millisecond)
	endTime := startTime + int64(timeout)

	// receive message
	for msg == nil {
		ret, msg, timeStamp, err = p.Read()
		if ret == PCAN_ERROR_QRCVEMPTY {
			if hasEvents {
				val, errWait := syscall.WaitForSingleObject(p.recvEvent, timeoutU32)
				switch val {
				case syscall.WAIT_OBJECT_0:
					break
				case syscall.WAIT_FAILED:
					return ret, nil, nil, errWait
				case syscall.WAIT_TIMEOUT:
					return ret, nil, nil, errWait
				default:
					return ret, msg, timeStamp, errWait
				}
			} else {
				// timeout handling
				if time.Now().UnixNano()/int64(time.Millisecond) > endTime {
					return ret, nil, nil, err
				}
				time.Sleep(250 * time.Microsecond)
			}
		}
	}

	return ret, msg, timeStamp, err

}

// Reads from device buffer until it has no more messages stored with an optional message limit
// If limit is set to zero, no limit will will be used
func (p *TPCANBus) ReadFullBuffer(limit int) ([]TPCANMsg, []TPCANTimestamp, error) {

	var ret = PCAN_ERROR_UNKNOWN
	var msg *TPCANMsg = nil
	var timeStamp *TPCANTimestamp = nil
	var msgs []TPCANMsg
	var timestamps []TPCANTimestamp
	var err error = nil

	// read until buffer empty is returned
	for {
		ret, msg, timeStamp, err = p.Read()
		if ret == PCAN_ERROR_QRCVEMPTY {
			return msgs, timestamps, err
		} else {
			msgs = append(msgs, *msg)
			timestamps = append(timestamps, *timeStamp)
			if limit != 0 && len(msgs) >= int(limit) {
				return msgs, timestamps, err
			}
		}
	}
}

// Reads a CAN message from the receive queue of a FD capable PCAN Channel
func (p *TPCANBusFD) ReadFD() (TPCANStatus, TPCANMsgFD, TPCANTimestampFD, error) {
	var msgFD TPCANMsgFD
	var timeStampFD TPCANTimestampFD

	ret, _, errno := pHandleReadFD.Call(uintptr(p.Handle), uintptr(unsafe.Pointer(&msgFD)), uintptr(unsafe.Pointer(&timeStampFD)))
	return TPCANStatus(ret), msgFD, timeStampFD, syscallErr(errno)
}

// Transmits a CAN message
// msg: A Message struct with the message to be sent
func (p *TPCANBus) Write(msg *TPCANMsg) (TPCANStatus, error) {
	ret, _, errno := pHandleWrite.Call(uintptr(p.Handle), uintptr(unsafe.Pointer(msg)))
	return TPCANStatus(ret), syscallErr(errno)
}

// Transmits a CAN message over a FD capable PCAN Channel
// msgFD A MessageFD struct with the message to be sent
func (p *TPCANBusFD) WriteFD(msgFD *TPCANMsgFD) (TPCANStatus, error) {
	ret, _, errno := pHandleWriteFD.Call(uintptr(p.Handle), uintptr(unsafe.Pointer(msgFD)))
	return TPCANStatus(ret), syscallErr(errno)
}

// Configures the reception filter
// fromID: The lowest CAN ID to be received
// toID: The highest CAN ID to be received
// mode: Message type, Standard (11-bit identifier) or Extended (29-bit identifier)
func (p *TPCANBus) SetFilter(fromID TPCANMsgID, toID TPCANMsgID, mode TPCANMode) (TPCANStatus, error) {
	ret, _, errno := pHandleFilterMessages.Call(uintptr(p.Handle), uintptr(fromID), uintptr(toID), uintptr(mode))
	if TPCANStatus(ret) != PCAN_ERROR_OK {
		return TPCANStatus(ret), syscallErr(errno)
	}
	return p.SetParameter(PCAN_MESSAGE_FILTER, TPCANParameterValue(PCAN_FILTER_CLOSE)) // confirm filter
}

// Resets message filter set by SetFilter() function
func (p *TPCANBus) ResetFilter() (TPCANStatus, error) {
	return p.SetParameter(PCAN_MESSAGE_FILTER, TPCANParameterValue(PCAN_FILTER_OPEN))
}

// Retrieves a PCAN Channel value using a defined parameter value type
// param: The TPCANParameter parameter to get
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func (p *TPCANBus) GetParameter(param TPCANParameter) (TPCANStatus, TPCANParameterValue, error) {
	var val TPCANParameterValue
	ret, err := p.GetValue(param, unsafe.Pointer(&val), uint32(unsafe.Sizeof(val)))
	return TPCANStatus(ret), val, err
}

// Configures a PCAN Channel value using a defined parameter value type
// param: The TPCANParameter parameter to set
// value: Value of parameter
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func (p *TPCANBus) SetParameter(param TPCANParameter, val TPCANParameterValue) (TPCANStatus, error) {
	ret, err := p.SetValue(param, unsafe.Pointer(&val), uint32(unsafe.Sizeof(val)))
	return TPCANStatus(ret), err
}

// Retrieves a PCAN Channel value
// param: The TPCANParameter parameter to get
// Note: Parameters can be present or not according with the kind
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func (p *TPCANBus) GetValue(param TPCANParameter, buffer unsafe.Pointer, bufferSize uint32) (TPCANStatus, error) {
	ret, _, errno := pHandleGetValue.Call(uintptr(p.Handle), uintptr(param), uintptr(buffer), uintptr(bufferSize))
	return TPCANStatus(ret), syscallErr(errno)
}

// Configures a PCAN Channel value.
// Channel: The handle of a PCAN Channel
// param: The TPCANParameter parameter to set
// value: Value of parameter
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func (p *TPCANBus) SetValue(param TPCANParameter, buffer unsafe.Pointer, bufferSize uint32) (TPCANStatus, error) {
	ret, _, errno := pHandleSetValue.Call(uintptr(p.Handle), uintptr(param), uintptr(buffer), uintptr(bufferSize))
	return TPCANStatus(ret), syscallErr(errno)
}

// Allows or forbids receiving of status frames
// allowStatusFrames: Allows status frames if set to true
func (p *TPCANBus) SetAllowStatusFrames(allowStatusFrames bool) (TPCANStatus, error) {
	var conv = map[bool]TPCANParameterValue{false: PCAN_PARAMETER_OFF, true: PCAN_PARAMETER_ON}
	return p.SetParameter(PCAN_ALLOW_STATUS_FRAMES, conv[allowStatusFrames])
}

// Allows or forbids receiving of remote transmission request frames frames
// allowStatusFrames: Allows remote transmission request frames if set to true
func (p *TPCANBus) SetAllowRTRFrames(allowRTRFrames bool) (TPCANStatus, error) {
	var conv = map[bool]TPCANParameterValue{false: PCAN_PARAMETER_OFF, true: PCAN_PARAMETER_ON}
	return p.SetParameter(PCAN_ALLOW_RTR_FRAMES, conv[allowRTRFrames])
}

// Allows or forbids receiving of error frames
// allowStatusFrames: Allows error frames if set to true
func (p *TPCANBus) SetAllowErrorFrames(allowErrorFrames bool) (TPCANStatus, error) {
	var conv = map[bool]TPCANParameterValue{false: PCAN_PARAMETER_OFF, true: PCAN_PARAMETER_ON}
	return p.SetParameter(PCAN_ALLOW_ERROR_FRAMES, conv[allowErrorFrames])
}

// Allows or forbids receiving of echo frames
// allowStatusFrames: Allows echo frames if set to true
func (p *TPCANBus) SetAllowEchoFrames(allowEchoFrames bool) (TPCANStatus, error) {
	var conv = map[bool]TPCANParameterValue{false: PCAN_PARAMETER_OFF, true: PCAN_PARAMETER_ON}
	return p.SetParameter(PCAN_ALLOW_ECHO_FRAMES, conv[allowEchoFrames])
}

// Turn on or off flashing of the device's LED for physical identification purposes
func (p *TPCANBus) SetLEDState(ledState bool) (TPCANStatus, error) {
	var conv = map[bool]TPCANParameterValue{false: PCAN_PARAMETER_OFF, true: PCAN_PARAMETER_ON}
	return p.SetParameter(PCAN_CHANNEL_IDENTIFYING, conv[ledState])
}

// Returns the channel condition as a level for availablity
func (p *TPCANBus) GetChannelCondition() (TPCANStatus, TPCANCHannelCondition, error) {
	state, val, err := p.GetParameter(PCAN_CHANNEL_CONDITION)
	return state, TPCANCHannelCondition(val), err
}

// Starts recording a trace on given path with a max file size in MB
// maxFileSize: trace file is splitted in files with this maximum size of file in MB; set to zero to have a infinite large trace file (max is 100 MB)
// Note: A trace file only gets filled if the Recv() function is called!
func (p *TPCANBus) StartTrace(filePath string, maxFileSize uint32) (TPCANStatus, error) {
	if maxFileSize > MAX_TRACE_FILE_SIZE_ACCEPTED {
		return PCAN_ERROR_UNKNOWN, fmt.Errorf("maximum size of a trace file is %v MB", MAX_TRACE_FILE_SIZE_ACCEPTED)
	}

	// configure trace configuration (only file size is set, the other options are always used)
	cfg := TRACE_FILE_DATE | TRACE_FILE_TIME | TRACE_FILE_OVERWRITE
	if maxFileSize > 0 {
		cfg |= TRACE_FILE_SEGMENTED
	} else {
		cfg |= TRACE_FILE_SINGLE
	}
	state, err := p.SetParameter(PCAN_TRACE_CONFIGURE, TPCANParameterValue(cfg))
	if err != nil || state != PCAN_ERROR_OK {
		return state, err
	}
	if maxFileSize > 0 {
		state, err := p.SetValue(PCAN_TRACE_SIZE, unsafe.Pointer(&maxFileSize), 4)
		if err != nil || state != PCAN_ERROR_OK {
			return state, err
		}
	}

	// configure trace file path
	if len(filePath) > MAX_LENGHT_STRING_BUFFER {
		return PCAN_ERROR_UNKNOWN, fmt.Errorf("filepath exceeds max length of %v", MAX_LENGHT_STRING_BUFFER)
	}

	// convert path to a fixed buffer size as pcan wants it that way
	var buffer [MAX_LENGHT_STRING_BUFFER]byte
	for i := range filePath {
		buffer[i] = byte(filePath[i])
	}
	state, err = p.SetValue(PCAN_TRACE_LOCATION, unsafe.Pointer(&buffer), uint32(unsafe.Sizeof(buffer)))
	if err != nil || state != PCAN_ERROR_OK {
		return state, err
	}

	// start tracing
	state, err = p.SetParameter(PCAN_TRACE_STATUS, PCAN_PARAMETER_ON)
	return state, err
}

// Stops recording currently running trace
func (p *TPCANBus) StopTrace() (TPCANStatus, error) {
	return p.SetParameter(PCAN_TRACE_STATUS, PCAN_PARAMETER_OFF)
}

// prepare WaitForSingleObject implementation when waiting for CAN messages (currently only windows support)
func (p *TPCANBus) initializeRecvEvent() {
	p.recvEvent = 0
	if hasEvents {
		modkernel32, errLoad := syscall.LoadLibrary("kernel32.dll")
		procCreateEvent, errOpen := syscall.GetProcAddress(modkernel32, "CreateEventW")
		if errLoad == nil && errOpen == nil && procCreateEvent != 0 {
			r0, _, errno := syscall.SyscallN(procCreateEvent)
			if errno == 0 && r0 != 0 && syscall.Handle(r0) != syscall.InvalidHandle {
				p.recvEvent = syscall.Handle(r0)
				retVal, errVal := p.SetParameter(PCAN_RECEIVE_EVENT, TPCANParameterValue(r0))
				if retVal != PCAN_ERROR_OK || errVal != nil {
					hasEvents = false
					_ = syscall.CloseHandle(p.recvEvent)
					p.recvEvent = 0
				}
			}
		}
		// just for safety
		if p.recvEvent == 0 || p.recvEvent == syscall.InvalidHandle {
			hasEvents = false
		}
	}
}

// Uninitializes all PCAN Channels initialized by CAN_Initialize
func ShutdownAllHandles() error {
	return nil // TODO

	// state, err := Uninitialize(PCAN_NONEBUS)
	// return evalRetval(state, err)
}

// Gets information about all existing PCAN channels on a system in a single call, regardless of their current availability.
func AttachedChannelsCount() (uint32, error) {
	return 0, nil // TODO

	// var channelCount uint32
	// ret, err := GetValue(PCAN_NONEBUS, PCAN_ATTACHED_CHANNELS_COUNT, unsafe.Pointer(&channelCount), uint32(unsafe.Sizeof(channelCount)))
	// if err != nil {
	// 	return channelCount, err
	// }
	// return channelCount, getFormattedError(ret)
}

// Returns list of all existing PCAN channels on a system in a single call, regardless of their current availability
func AttachedChannels() ([]TPCANHandle, error) {
	return nil, nil // TODO

	//posChannels := [...]TPCANHandle{PCAN_USBBUS1, PCAN_USBBUS2, PCAN_USBBUS3, PCAN_USBBUS4,
	//	PCAN_USBBUS5, PCAN_USBBUS6, PCAN_USBBUS7, PCAN_USBBUS8,
	//	PCAN_USBBUS9, PCAN_USBBUS10, PCAN_USBBUS11, PCAN_USBBUS12,
	//	PCAN_USBBUS13, PCAN_USBBUS14, PCAN_USBBUS15, PCAN_USBBUS16}
	//attachedChannels := []TPCANHandle{}
	//
	//// iterate through channels and check for every channel if available
	//for i := range posChannels {
	//	state, cond, err := GetParameter(posChannels[i], PCAN_CHANNEL_CONDITION)
	//	if state != PCAN_ERROR_OK || err != nil {
	//		return nil, err
	//	}
	//	if cond == TPCANParameterValue(PCAN_CHANNEL_AVAILABLE) ||
	//		cond == TPCANParameterValue(PCAN_CHANNEL_OCCUPIED) ||
	//		cond == TPCANParameterValue(PCAN_CHANNEL_PCANVIEW) {
	//		attachedChannels = append(attachedChannels, posChannels[i])
	//	}
	//}
	//
	//return attachedChannels, nil
}

// Returns list of all existing PCAN channels on a system in a single call, regardless of their current availability
// TODO This function is not working correctly, as the given information does not matched connected hardware, use AttachedChannels instead
func AttachedChannels_Extended() ([]TPCANChannelInformation, error) {
	log.Fatalf("This function is not working correctly, as the given information does not matched connected hardware, use AttachedChannels instead!") // TODO
	return nil, nil                                                                                                                                   // TODO

	//count, err := AttachedChannelsCount()
	//if err != nil || count == 0 { // size calculation not possible with a slice len of 0
	//	return nil, err
	//}
	//
	//buf := make([]TPCANChannelInformation, count)
	//size := uintptr(len(buf)) * unsafe.Sizeof(buf[0])
	//state, err := GetValue(PCAN_NONEBUS, PCAN_ATTACHED_CHANNELS, unsafe.Pointer(&buf[0]), uint32(size))
	//
	//return buf, evalRetval(state, err)
}

// Returns a descriptive text of a given TPCANStatus error code, in any desired language
// err: A TPCANStatus error code
// language: Indicates a 'Primary language ID'
func GetErrorText(status TPCANStatus, language TPCANLanguage) (TPCANStatus, [MAX_LENGHT_STRING_BUFFER]byte, error) {
	var buffer [MAX_LENGHT_STRING_BUFFER]byte

	ret, _, errno := pHandleGetErrorText.Call(uintptr(status), uintptr(language), uintptr(unsafe.Pointer(&buffer)))
	return TPCANStatus(ret), buffer, syscallErr(errno)
}

// Finds a PCAN-Basic Channel that matches with the given parameters
// parameters: A comma separated string contained pairs of parameter-name/value to be matched within a PCAN-Basic Channel
// foundChannels: Buffer for returning the PCAN-Basic Channel when found
func LookUpChannel(deviceType string, deviceID string, controllerNumber string, ipAdress string) (TPCANStatus, TPCANHandle, error) {

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

	ret, _, errno := pHandleLookUpChannel.Call(uintptr(unsafe.Pointer(&sParameters)), uintptr(unsafe.Pointer(&foundChannel)))
	return TPCANStatus(ret), foundChannel, syscallErr(errno)
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
