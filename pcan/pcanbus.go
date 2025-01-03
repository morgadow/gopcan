package pcan

import (
	"errors"
	"fmt"
	"log"
	"syscall"
	"time"
	"unsafe"
)

/* Generic bus implementation laoding the os specific files which are hanlding the api calls.  */

// Initializes a basic plugNplay PCAN Channel
// Channel: The handle of a PCAN Channel
// baudRate: The speed for the communication (BTR0BTR1 code)
func InitializeBasic(handle TPCANHandle, baudRate TPCANBaudrate) (TPCANStatus, *TPCANBus, error) {
	LoadAPI()

	status, err := APIInitializeBasic(handle, baudRate)
	if status != PCAN_ERROR_OK || err != nil {
		return status, nil, err
	}

	bus := TPCANBus{
		Handle:    handle,
		Baudrate:  baudRate,
		HWType:    PCAN_DEFAULT_HW_TYPE,
		IOPort:    PCAN_DEFAULT_IO_PORT,
		Interrupt: PCAN_DEFAULT_INTERRUPT}

	bus.initializeRecvEvent()

	return status, &bus, err
}

// Initializes a advanced PCAN Channel
// Channel: The handle of a PCAN Channel
// baudRate: The speed for the communication (BTR0BTR1 code)
// hwType: Non-PnP: The type of hardware and operation mode
// ioPort: Non-PnP: The I/O address for the parallel port
// interrupt: Non-PnP: Interrupt number of the parallel port
func Initialize(handle TPCANHandle, baudRate TPCANBaudrate, hwType TPCANType, ioPort uint32, interrupt uint16) (TPCANStatus, *TPCANBus, error) {
	LoadAPI()

	status, err := APIInitialize(handle, baudRate, hwType, ioPort, interrupt)
	if status != PCAN_ERROR_OK || err != nil {
		return status, nil, err
	}

	bus := TPCANBus{
		Handle:    handle,
		Baudrate:  baudRate,
		HWType:    hwType,
		IOPort:    ioPort,
		Interrupt: interrupt}

	bus.initializeRecvEvent()

	return status, &bus, err
}

// Initializes a FD capable PCAN Channel
// handle: The handle of a PCAN Channel
// bitRateFD: The speed for the communication (FD bit rate string)
// Note:
// Bit rate string must follow the following construction rules:
//   - parameter and values must be separated by '='
//   - Couples of Parameter/value must be separated by ','
//   - Following Parameter must be filled out: f_clock, data_brp, data_sjw, data_tseg1, data_tseg2,
//     nom_brp, nom_sjw, nom_tseg1, nom_tseg2.
//   - Following Parameters are optional (not used yet): data_ssp_offset, nom_sam
//   - Example: f_clock=80000000,nom_brp=10,nom_tseg1=5,nom_tseg2=2,nom_sjw=1,data_brp=4,data_tseg1=7,data_tseg2=2,data_sjw=1
func InitializeFD(handle TPCANHandle, bitRateFD TPCANBitrateFD) (TPCANStatus, *TPCANBusFD, error) {
	LoadAPI()

	status, err := APIInitializeFD(handle, bitRateFD)
	if status != PCAN_ERROR_OK || err != nil {
		return status, nil, err
	}

	return PCAN_ERROR_UNKNOWN, nil, errors.New("not implemented") // TODO
}

// Uninitializes PCAN Channels initialized by CAN_Initialize
func (p *TPCANBus) Uninitialize() (TPCANStatus, error) {
	return APIUninitialize(p.Handle)
}

// Resets the receive and transmit queues of the PCAN Channel
func (p *TPCANBus) Reset() (TPCANStatus, error) {
	return APIReset(p.Handle)
}

// Gets the current status of a PCAN Channel
func (p *TPCANBus) GetStatus() (TPCANStatus, error) {
	return APIGetStatus(p.Handle)
}

// Reads a CAN message from the receive queue of a PCAN Channel
// Note: Does return nil if receive buffer is empty
func (p *TPCANBus) Read() (TPCANStatus, *TPCANMsg, *TPCANTimestamp, error) {
	status, msg, timestamp, err := APIRead(p.Handle)
	if status == PCAN_ERROR_QRCVEMPTY {
		return status, nil, nil, err
	} else {
		return status, &msg, &timestamp, err
	}
}

// Reads a CAN message from the receive queue of a PCAN Channel with an timeout and only returns a valid messsage
// Note: Does return nil if receive buffer is empty or no message is read during timeout
// timeout: Timeout for receiving message from CAN bus in milliseconds (if set below zero, no timeout is set)
func (p *TPCANBus) ReadWithTimeout(timeout int) (TPCANStatus, *TPCANMsg, *TPCANTimestamp, error) {

	var ret = PCAN_ERROR_UNKNOWN
	var msg *TPCANMsg = nil
	var timestamp *TPCANTimestamp = nil
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
		ret, msg, timestamp, err = p.Read()
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
					return ret, msg, timestamp, errWait
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

	return ret, msg, timestamp, err

}

// Reads from device buffer until it has no more messages stored with an optional message limit
// If limit is set to zero, no limit will will be used
func (p *TPCANBus) ReadFullBuffer(limit int) ([]TPCANMsg, []TPCANTimestamp, error) {

	var ret = PCAN_ERROR_UNKNOWN
	var msg *TPCANMsg = nil
	var timestamp *TPCANTimestamp = nil
	var msgs []TPCANMsg
	var timestamps []TPCANTimestamp
	var err error = nil

	// read until buffer empty is returned
	for {
		ret, msg, timestamp, err = p.Read()
		if ret == PCAN_ERROR_QRCVEMPTY {
			return msgs, timestamps, err
		} else {
			msgs = append(msgs, *msg)
			timestamps = append(timestamps, *timestamp)
			if limit != 0 && len(msgs) >= int(limit) {
				return msgs, timestamps, err
			}
		}
	}
}

// Reads a CAN message from the receive queue of a FD capable PCAN Channel
func (p *TPCANBusFD) ReadFD() (TPCANStatus, *TPCANMsgFD, *TPCANTimestampFD, error) {
	status, msg, timestamp, err := APIReadFD(p.Handle)
	if status == PCAN_ERROR_QRCVEMPTY {
		return status, nil, nil, err
	} else {
		return status, &msg, &timestamp, err
	}
}

// Transmits a CAN message
// msg: A Message struct with the message to be sent
func (p *TPCANBus) Write(msg *TPCANMsg) (TPCANStatus, error) {
	return APIWrite(p.Handle, msg)
}

// Transmits a CAN message over a FD capable PCAN Channel
// msgFD A MessageFD struct with the message to be sent
func (p *TPCANBusFD) WriteFD(msg *TPCANMsgFD) (TPCANStatus, error) {
	return APIWriteFD(p.Handle, msg)
}

// Configures the reception filter
// fromID: The lowest CAN ID to be received
// toID: The highest CAN ID to be received
// mode: Message type, Standard (11-bit identifier) or Extended (29-bit identifier)
func (p *TPCANBus) SetFilter(fromID TPCANMsgID, toID TPCANMsgID, mode TPCANMode) (TPCANStatus, error) {
	status, err := APISetFilter(p.Handle, fromID, toID, mode)
	if status != PCAN_ERROR_OK {
		return status, err
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
	status, err := p.GetValue(param, unsafe.Pointer(&val), uint32(unsafe.Sizeof(val)))
	return status, val, err
}

// Configures a PCAN Channel value using a defined parameter value type
// param: The TPCANParameter parameter to set
// value: Value of parameter
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func (p *TPCANBus) SetParameter(param TPCANParameter, val TPCANParameterValue) (TPCANStatus, error) {
	status, err := p.SetValue(param, unsafe.Pointer(&val), uint32(unsafe.Sizeof(val)))
	return status, err
}

// Retrieves a PCAN Channel value
// param: The TPCANParameter parameter to get
// Note: Parameters can be present or not according with the kind
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func (p *TPCANBus) GetValue(param TPCANParameter, buffer unsafe.Pointer, bufferSize uint32) (TPCANStatus, error) {
	return APIGetValue(p.Handle, param, buffer, bufferSize)
}

// Configures a PCAN Channel value.
// Channel: The handle of a PCAN Channel
// param: The TPCANParameter parameter to set
// value: Value of parameter
// Note: Parameters can be present or not according with the kind of Hardware (PCAN Channel) being used.
// If a parameter is not available, a PCAN_ERROR_ILLPARAMTYPE error will be returned
func (p *TPCANBus) SetValue(param TPCANParameter, buffer unsafe.Pointer, bufferSize uint32) (TPCANStatus, error) {
	return APISetValue(p.Handle, param, buffer, bufferSize)
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
func ShutdownAllHandles() (TPCANStatus, error) {
	return APIUninitialize(PCAN_NONEBUS)
}

// Gets information about all existing PCAN channels on a system in a single call, regardless of their current availability.
func AttachedChannelsCount() (TPCANStatus, uint32, error) {
	var channelCount uint32

	status, err := APIGetValue(PCAN_NONEBUS, PCAN_ATTACHED_CHANNELS_COUNT, unsafe.Pointer(&channelCount), uint32(unsafe.Sizeof(channelCount)))
	if err != nil {
		return status, channelCount, err
	}
	return status, channelCount, err
}

// Returns list of all existing PCAN channels on a system in a single call, regardless of their current availability
func AttachedChannels() ([]TPCANHandle, error) {
	posChannels := [...]TPCANHandle{PCAN_USBBUS1, PCAN_USBBUS2, PCAN_USBBUS3, PCAN_USBBUS4,
		PCAN_USBBUS5, PCAN_USBBUS6, PCAN_USBBUS7, PCAN_USBBUS8,
		PCAN_USBBUS9, PCAN_USBBUS10, PCAN_USBBUS11, PCAN_USBBUS12,
		PCAN_USBBUS13, PCAN_USBBUS14, PCAN_USBBUS15, PCAN_USBBUS16}
	attachedChannels := []TPCANHandle{}

	// iterate through channels and check for every channel if available
	var cond TPCANParameterValue
	for i := range posChannels {
		state, err := APIGetValue(posChannels[i], PCAN_CHANNEL_CONDITION, unsafe.Pointer(&cond), uint32(unsafe.Sizeof(cond)))
		if state != PCAN_ERROR_OK || err != nil {
			return nil, err
		}
		if cond == TPCANParameterValue(PCAN_CHANNEL_AVAILABLE) ||
			cond == TPCANParameterValue(PCAN_CHANNEL_OCCUPIED) ||
			cond == TPCANParameterValue(PCAN_CHANNEL_PCANVIEW) {
			attachedChannels = append(attachedChannels, posChannels[i])
		}
	}

	return attachedChannels, nil
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

// Finds a PCAN-Basic Channel that matches with the given parameters
// parameters: A comma separated string contained pairs of parameter-name/value to be matched within a PCAN-Basic Channel
// foundChannels: Buffer for returning the PCAN-Basic Channel when found
func LookUpChannel(deviceType string, deviceID string, controllerNumber string, ipAdress string) (TPCANStatus, TPCANHandle, error) {
	return APILookUpChannel(deviceType, deviceID, controllerNumber, ipAdress)
}
