package pcan

// Represents a PCAN message
type TPCANMsg struct {
	ID      TPCANMsgID                    // 11/29-bit message identifier
	MsgType TPCANMessageType              // Type of the message
	DLC     uint8                         // Data Length Code of the message (0..8)
	Data    [LENGTH_DATA_CAN_MESSAGE]byte // Data of the message (DATA[0]..DATA[7])
}

// Represents a timestamp of a received PCAN message
// Total Microseconds = micros + (1000ULL * millis) + (0x100000000ULL * 1000ULL * millis_overflow)
type TPCANTimestamp struct {
	Millis         uint32 // Base-value: milliseconds: 0.. 2^32-1
	MillisOverflow uint16 // Roll-arounds of millis
	Micros         uint16 // Microseconds: 0..999
}

// Represents a PCAN message from a FD capable hardware
type TPCANMsgFD struct {
	ID      TPCANMsgID
	MsgType TPCANMessageType
	DLC     uint8
	Data    [LENGTH_DATA_CANFD_MESSAGE]byte
}

// Describes an available PCAN channel
type TPCANChannelInformation struct {
	Channel          TPCANHandle                    // PCAN channel handle
	DeviceType       TPCANDevice                    // Kind of PCAN device
	ControllerNumber uint8                          // CAN-Controller number
	DeviceFeatures   uint32                         // Device capabilities flag (see FEATURE_*)
	DeviceName       [MAX_LENGTH_HARDWARE_NAME]rune // Device name
	DeviceID         uint32                         // Device number
	ChannelCondition TPCANCHannelCondition          // Availability status of a PCAN-Channel
}
