package pcan

type (
	TPCANLanguage         uint16              // Represents a language chosen for the error messages
	TPCANHandle           uint16              // Represents a PCAN hardware channel handle
	TPCANStatus           uint32              // Represents a PCAN status/error code
	TPCANDevice           uint8               // Represents a PCAN device
	TPCANParameter        uint8               // Represents a PCAN parameter to be read or set
	TPCANParameterValue   uint32              // Represents a PCAN parameter value
	TPCANFilterValue      TPCANParameterValue // Represents a PCAN filter parameter value
	TPCANCHannelCondition TPCANParameterValue // Represents a PCAN channel condition value
	TPCANFunctionValue    TPCANParameterValue // Represents a PCAN function parameter value
	TPCANTraceFileValue   TPCANParameterValue // Represents a PCAN trace file parameter value
	TPCANFeatureValue     TPCANParameterValue // Represents a PCAN feature parameter value
	TPCANStatusValue      TPCANParameterValue // Represents a PCAN status parameter value
	TPCANMessageType      uint8               // Represents the type of a PCAN message
	TPCANMode             uint8               // Represents a PCAN filter mode
	TPCANBaudrate         uint16              // Represents a PCAN Baud rate register value (BTR0/BTR1 register values for the CAN controller)
	TPCANType             uint8               // Represents the type of PCAN hardware to be initialized
	TPCANLookupParameter  string              // LookUp Parameters
	TPCANBRParameter      string              // Represents the configuration for a CAN bit rate (Example: f_clock=80000000,nom_brp=10,nom_tseg1=5,nom_tseg2=2,nom_sjw=1,data_brp=4,data_tseg1=7,data_tseg2=2,data_sjw=1)
	TPCANBitrateFD        string              // Represents a PCAN-FD bit rate string
	TPCANMsgID            uint32              // 11/29-bit message identifier
	TPCANTimestampFD      uint64              // Represents a timestamp of a received PCAN FD message
)

// constants for driver handling
const (
	LENGTH_DATA_CAN_MESSAGE      = 8                        // maximum amount of bytes in an PCAN CAN message
	LENGTH_DATA_CANFD_MESSAGE    = 64                       // maximum amount of bytes in can CAN FD message
	MAX_LENGTH_HARDWARE_NAME     = 33                       // Maximum length of the name of a device: 32 characters + terminator
	MAX_LENGHT_STRING_BUFFER     = 256                      // Maximum length of any string buffer sent or received from pcan dll
	MAX_LENGTH_VERSION_STRING    = MAX_LENGHT_STRING_BUFFER // Maximum length of a version string: 255 characters + terminator
	MAX_TRACE_FILE_SIZE_ACCEPTED = 100                      // Maximum size of a trace file in MB

	PCAN_DEFAULT_HW_TYPE   TPCANType = PCAN_TYPE_ISA // Default hardware type for a plug-n-play channel
	PCAN_DEFAULT_IO_PORT   uint32    = 0x02A0        // Default IO port for a plug-n-play channel
	PCAN_DEFAULT_INTERRUPT uint16    = 11            // Default interrupt id for a plug-n-play channel
)

// Represents a language chosen for the error messages
const (
	LanguageNeutral TPCANLanguage = 0x00
	LanguageGerman  TPCANLanguage = 0x07
	LanguageEnglish TPCANLanguage = 0x09
	LanguageItalian TPCANLanguage = 0x10
	LanguageSpanish TPCANLanguage = 0x0A
	LanguageFrench  TPCANLanguage = 0x0C
)

// Represents a PCAN hardware channel handle
const (
	PCAN_NONEBUS = TPCANHandle(0x00) // Undefined/default value for a PCAN bus

	PCAN_ISABUS1 = TPCANHandle(0x21) // PCAN-ISA interface, Channel 1
	PCAN_ISABUS2 = TPCANHandle(0x22) // PCAN-ISA interface, Channel 2
	PCAN_ISABUS3 = TPCANHandle(0x23) // PCAN-ISA interface, Channel 3
	PCAN_ISABUS4 = TPCANHandle(0x24) // PCAN-ISA interface, Channel 4
	PCAN_ISABUS5 = TPCANHandle(0x25) // PCAN-ISA interface, Channel 5
	PCAN_ISABUS6 = TPCANHandle(0x26) // PCAN-ISA interface, Channel 6
	PCAN_ISABUS7 = TPCANHandle(0x27) // PCAN-ISA interface, Channel 7
	PCAN_ISABUS8 = TPCANHandle(0x28) // PCAN-ISA interface, Channel 8

	PCAN_DNGBUS1 = TPCANHandle(0x31) // PCAN-Dongle/LPT interface, Channel 1

	PCAN_PCIBUS1  = TPCANHandle(0x41)  // PCAN-PCI interface, Channel 1
	PCAN_PCIBUS2  = TPCANHandle(0x42)  // PCAN-PCI interface, Channel 2
	PCAN_PCIBUS3  = TPCANHandle(0x43)  // PCAN-PCI interface, Channel 3
	PCAN_PCIBUS4  = TPCANHandle(0x44)  // PCAN-PCI interface, Channel 4
	PCAN_PCIBUS5  = TPCANHandle(0x45)  // PCAN-PCI interface, Channel 5
	PCAN_PCIBUS6  = TPCANHandle(0x46)  // PCAN-PCI interface, Channel 6
	PCAN_PCIBUS7  = TPCANHandle(0x47)  // PCAN-PCI interface, Channel 7
	PCAN_PCIBUS8  = TPCANHandle(0x48)  // PCAN-PCI interface, Channel 8
	PCAN_PCIBUS9  = TPCANHandle(0x409) // PCAN-PCI interface, Channel 9
	PCAN_PCIBUS10 = TPCANHandle(0x40A) // PCAN-PCI interface, Channel 10
	PCAN_PCIBUS11 = TPCANHandle(0x40B) // PCAN-PCI interface, Channel 11
	PCAN_PCIBUS12 = TPCANHandle(0x40C) // PCAN-PCI interface, Channel 12
	PCAN_PCIBUS13 = TPCANHandle(0x40D) // PCAN-PCI interface, Channel 13
	PCAN_PCIBUS14 = TPCANHandle(0x40E) // PCAN-PCI interface, Channel 14
	PCAN_PCIBUS15 = TPCANHandle(0x40F) // PCAN-PCI interface, Channel 15
	PCAN_PCIBUS16 = TPCANHandle(0x410) // PCAN-PCI interface, Channel 16

	PCAN_USBBUS1  = TPCANHandle(0x51)  // PCAN-USB interface, Channel 1
	PCAN_USBBUS2  = TPCANHandle(0x52)  // PCAN-USB interface, Channel 2
	PCAN_USBBUS3  = TPCANHandle(0x53)  // PCAN-USB interface, Channel 3
	PCAN_USBBUS4  = TPCANHandle(0x54)  // PCAN-USB interface, Channel 4
	PCAN_USBBUS5  = TPCANHandle(0x55)  // PCAN-USB interface, Channel 5
	PCAN_USBBUS6  = TPCANHandle(0x56)  // PCAN-USB interface, Channel 6
	PCAN_USBBUS7  = TPCANHandle(0x57)  // PCAN-USB interface, Channel 7
	PCAN_USBBUS8  = TPCANHandle(0x58)  // PCAN-USB interface, Channel 8
	PCAN_USBBUS9  = TPCANHandle(0x509) // PCAN-USB interface, Channel 9
	PCAN_USBBUS10 = TPCANHandle(0x50A) // PCAN-USB interface, Channel 10
	PCAN_USBBUS11 = TPCANHandle(0x50B) // PCAN-USB interface, Channel 11
	PCAN_USBBUS12 = TPCANHandle(0x50C) // PCAN-USB interface, Channel 12
	PCAN_USBBUS13 = TPCANHandle(0x50D) // PCAN-USB interface, Channel 13
	PCAN_USBBUS14 = TPCANHandle(0x50E) // PCAN-USB interface, Channel 14
	PCAN_USBBUS15 = TPCANHandle(0x50F) // PCAN-USB interface, Channel 15
	PCAN_USBBUS16 = TPCANHandle(0x510) // PCAN-USB interface, Channel 16

	PCAN_PCCBUS1 = TPCANHandle(0x61) // PCAN-PC Card interface, Channel 1
	PCAN_PCCBUS2 = TPCANHandle(0x62) // PCAN-PC Card interface, Channel 2

	PCAN_LANBUS1  = TPCANHandle(0x801) // PCAN-LAN interface, Channel 1
	PCAN_LANBUS2  = TPCANHandle(0x802) // PCAN-LAN interface, Channel 2
	PCAN_LANBUS3  = TPCANHandle(0x803) // PCAN-LAN interface, Channel 3
	PCAN_LANBUS4  = TPCANHandle(0x804) // PCAN-LAN interface, Channel 4
	PCAN_LANBUS5  = TPCANHandle(0x805) // PCAN-LAN interface, Channel 5
	PCAN_LANBUS6  = TPCANHandle(0x806) // PCAN-LAN interface, Channel 6
	PCAN_LANBUS7  = TPCANHandle(0x807) // PCAN-LAN interface, Channel 7
	PCAN_LANBUS8  = TPCANHandle(0x808) // PCAN-LAN interface, Channel 8
	PCAN_LANBUS9  = TPCANHandle(0x809) // PCAN-LAN interface, Channel 9
	PCAN_LANBUS10 = TPCANHandle(0x80A) // PCAN-LAN interface, Channel 10
	PCAN_LANBUS11 = TPCANHandle(0x80B) // PCAN-LAN interface, Channel 11
	PCAN_LANBUS12 = TPCANHandle(0x80C) // PCAN-LAN interface, Channel 12
	PCAN_LANBUS13 = TPCANHandle(0x80D) // PCAN-LAN interface, Channel 13
	PCAN_LANBUS14 = TPCANHandle(0x80E) // PCAN-LAN interface, Channel 14
	PCAN_LANBUS15 = TPCANHandle(0x80F) // PCAN-LAN interface, Channel 15
	PCAN_LANBUS16 = TPCANHandle(0x810) // PCAN-LAN interface, Channel 16
)

// Represents a PCAN status/error code
const (
	PCAN_ERROR_OK           = TPCANStatus(0x00000)                                                                                            // No error
	PCAN_ERROR_XMTFULL      = TPCANStatus(0x00001)                                                                                            // Transmit buffer in CAN controller is full
	PCAN_ERROR_OVERRUN      = TPCANStatus(0x00002)                                                                                            // CAN controller was read too late
	PCAN_ERROR_BUSLIGHT     = TPCANStatus(0x00004)                                                                                            // Bus error: an error counter reached the 'light' limit
	PCAN_ERROR_BUSHEAVY     = TPCANStatus(0x00008)                                                                                            // Bus error: an error counter reached the 'heavy' limit
	PCAN_ERROR_BUSWARNING   = PCAN_ERROR_BUSHEAVY                                                                                             // Bus error: an error counter reached the 'warning' limit
	PCAN_ERROR_BUSPASSIVE   = TPCANStatus(0x40000)                                                                                            // Bus error: the CAN controller is error passive
	PCAN_ERROR_BUSOFF       = TPCANStatus(0x00010)                                                                                            // Bus error: the CAN controller is in bus-off state
	PCAN_ERROR_ANYBUSERR    = (PCAN_ERROR_BUSWARNING | PCAN_ERROR_BUSLIGHT | PCAN_ERROR_BUSHEAVY | PCAN_ERROR_BUSOFF | PCAN_ERROR_BUSPASSIVE) // Mask for all bus errors
	PCAN_ERROR_QRCVEMPTY    = TPCANStatus(0x00020)                                                                                            // Receive queue is empty
	PCAN_ERROR_QOVERRUN     = TPCANStatus(0x00040)                                                                                            // Receive queue was read too late
	PCAN_ERROR_QXMTFULL     = TPCANStatus(0x00080)                                                                                            // Transmit queue is full
	PCAN_ERROR_REGTEST      = TPCANStatus(0x00100)                                                                                            // Test of the CAN controller hardware registers failed (no hardware found)
	PCAN_ERROR_NODRIVER     = TPCANStatus(0x00200)                                                                                            // Driver not loaded
	PCAN_ERROR_HWINUSE      = TPCANStatus(0x00400)                                                                                            // Hardware already in use by a Net
	PCAN_ERROR_NETINUSE     = TPCANStatus(0x00800)                                                                                            // A Client is already connected to the Net
	PCAN_ERROR_ILLHW        = TPCANStatus(0x01400)                                                                                            // Hardware handle is invalid
	PCAN_ERROR_ILLNET       = TPCANStatus(0x01800)                                                                                            // Net handle is invalid
	PCAN_ERROR_ILLCLIENT    = TPCANStatus(0x01C00)                                                                                            // Client handle is invalid
	PCAN_ERROR_ILLHANDLE    = (PCAN_ERROR_ILLHW | PCAN_ERROR_ILLNET | PCAN_ERROR_ILLCLIENT)                                                   // Mask for all handle errors
	PCAN_ERROR_RESOURCE     = TPCANStatus(0x02000)                                                                                            // Resource (FIFO, Client, timeout) cannot be created
	PCAN_ERROR_ILLPARAMTYPE = TPCANStatus(0x04000)                                                                                            // Invalid parameter
	PCAN_ERROR_ILLPARAMVAL  = TPCANStatus(0x08000)                                                                                            // Invalid parameter value
	PCAN_ERROR_UNKNOWN      = TPCANStatus(0x10000)                                                                                            // Unknown error
	PCAN_ERROR_ILLDATA      = TPCANStatus(0x20000)                                                                                            // Invalid data, function, or action.
	PCAN_ERROR_ILLMODE      = TPCANStatus(0x80000)                                                                                            // Driver object state is wrong for the attempted operation
	PCAN_ERROR_CAUTION      = TPCANStatus(0x2000000)                                                                                          // An operation was successfully carried out, however, irregularities were registered
	PCAN_ERROR_INITIALIZE   = TPCANStatus(0x4000000)                                                                                          // Channel is not initialized
	PCAN_ERROR_ILLOPERATION = TPCANStatus(0x8000000)                                                                                          // Invalid operation
)

// Represents a PCAN device
const (
	PCAN_NONE    = TPCANDevice(0x0) // Undefined, unknown or not selected PCAN device value
	PCAN_PEAKCAN = TPCANDevice(0x1) // PCAN Non-PnP devices. NOT USED WITHIN PCAN-Basic API
	PCAN_ISA     = TPCANDevice(0x2) // PCAN-ISA, PCAN-PC/104, and PCAN-PC/104-Plus
	PCAN_DNG     = TPCANDevice(0x3) // PCAN-Dongle
	PCAN_PCI     = TPCANDevice(0x4) // PCAN-PCI, PCAN-cPCI, PCAN-miniPCI, and PCAN-PCI Express
	PCAN_USB     = TPCANDevice(0x5) // PCAN-USB and PCAN-USB Pro
	PCAN_PCC     = TPCANDevice(0x6) // PCAN-PC Card
	PCAN_VIRTUAL = TPCANDevice(0x7) // PCAN Virtual hardware. NOT USED WITHIN PCAN-Basic API
	PCAN_LAN     = TPCANDevice(0x8) // PCAN Gateway devices
)

// Represents a PCAN parameter to be read or set
const (
	PCAN_DEVICE_ID                = TPCANParameter(1)  // Device identifier parameter
	PCAN_DEVICE_NUMBER            = PCAN_DEVICE_ID     // DEPRECATED parameter. Use PCAN_DEVICE_ID instead
	PCAN_5VOLTS_POWER             = TPCANParameter(2)  // 5-Volt power parameter
	PCAN_RECEIVE_EVENT            = TPCANParameter(3)  // PCAN receive event handler parameter
	PCAN_MESSAGE_FILTER           = TPCANParameter(4)  // PCAN message filter parameter
	PCAN_API_VERSION              = TPCANParameter(5)  // PCAN-Basic API version parameter
	PCAN_CHANNEL_VERSION          = TPCANParameter(6)  // PCAN device channel version parameter
	PCAN_BUSOFF_AUTORESET         = TPCANParameter(7)  // PCAN Reset-On-Busoff parameter
	PCAN_LISTEN_ONLY              = TPCANParameter(8)  // PCAN Listen-Only parameter
	PCAN_LOG_LOCATION             = TPCANParameter(9)  // Directory path for log files
	PCAN_LOG_STATUS               = TPCANParameter(10) // Debug-Log activation status
	PCAN_LOG_CONFIGURE            = TPCANParameter(11) // Configuration of the debugged information (LOG_FUNCTION_***)
	PCAN_LOG_TEXT                 = TPCANParameter(12) // Custom insertion of text into the log file
	PCAN_CHANNEL_CONDITION        = TPCANParameter(13) // Availability status of a PCAN-Channel
	PCAN_HARDWARE_NAME            = TPCANParameter(14) // PCAN hardware name parameter
	PCAN_RECEIVE_STATUS           = TPCANParameter(15) // Message reception status of a PCAN-Channel
	PCAN_CONTROLLER_NUMBER        = TPCANParameter(16) // CAN-Controller number of a PCAN-Channel
	PCAN_TRACE_LOCATION           = TPCANParameter(17) // Directory path for PCAN trace files
	PCAN_TRACE_STATUS             = TPCANParameter(18) // CAN tracing activation status
	PCAN_TRACE_SIZE               = TPCANParameter(19) // Configuration of the maximum file size of a CAN trace
	PCAN_TRACE_CONFIGURE          = TPCANParameter(20) // Configuration of the trace file storing mode (TRACE_FILE_***)
	PCAN_CHANNEL_IDENTIFYING      = TPCANParameter(21) // Physical identification of a USB based PCAN-Channel by blinking its associated LED
	PCAN_CHANNEL_FEATURES         = TPCANParameter(22) // Capabilities of a PCAN device (FEATURE_***)
	PCAN_BITRATE_ADAPTING         = TPCANParameter(23) // Using of an existing bit rate (PCAN-View connected to a channel)
	PCAN_BITRATE_INFO             = TPCANParameter(24) // Configured bit rate as Btr0Btr1 value
	PCAN_BITRATE_INFO_FD          = TPCANParameter(25) // Configured bit rate as TPCANBitrateFD string
	PCAN_BUSSPEED_NOMINAL         = TPCANParameter(26) // Configured nominal CAN Bus speed as Bits per seconds
	PCAN_BUSSPEED_DATA            = TPCANParameter(27) // Configured CAN data speed as Bits per seconds
	PCAN_IP_ADDRESS               = TPCANParameter(28) // Remote address of a LAN channel as string in IPv4 format
	PCAN_LAN_SERVICE_STATUS       = TPCANParameter(29) // Status of the Virtual PCAN-Gateway Service
	PCAN_ALLOW_STATUS_FRAMES      = TPCANParameter(30) // Status messages reception status within a PCAN-Channel
	PCAN_ALLOW_RTR_FRAMES         = TPCANParameter(31) // RTR messages reception status within a PCAN-Channel
	PCAN_ALLOW_ERROR_FRAMES       = TPCANParameter(32) // Error messages reception status within a PCAN-Channel
	PCAN_INTERFRAME_DELAY         = TPCANParameter(33) // Delay, in microseconds, between sending frames
	PCAN_ACCEPTANCE_FILTER_11BIT  = TPCANParameter(34) // Filter over code and mask patterns for 11-Bit messages
	PCAN_ACCEPTANCE_FILTER_29BIT  = TPCANParameter(35) // Filter over code and mask patterns for 29-Bit messages
	PCAN_IO_DIGITAL_CONFIGURATION = TPCANParameter(36) // Output mode of 32 digital I/O pin of a PCAN-USB Chip. 1: Output-Active 0 : Output Inactive
	PCAN_IO_DIGITAL_VALUE         = TPCANParameter(37) // Value assigned to a 32 digital I/O pins of a PCAN-USB Chip
	PCAN_IO_DIGITAL_SET           = TPCANParameter(38) // Value assigned to a 32 digital I/O pins of a PCAN-USB Chip - Multiple digital I/O pins to 1 = High
	PCAN_IO_DIGITAL_CLEAR         = TPCANParameter(39) // Clear multiple digital I/O pins to 0
	PCAN_IO_ANALOG_VALUE          = TPCANParameter(40) // Get value of a single analog input pin
	PCAN_FIRMWARE_VERSION         = TPCANParameter(41) // Get the version of the firmware used by the device associated with a PCAN-Channel
	PCAN_ATTACHED_CHANNELS_COUNT  = TPCANParameter(42) // Get the amount of PCAN channels attached to a system
	PCAN_ATTACHED_CHANNELS        = TPCANParameter(43) // Get information about PCAN channels attached to a system
	PCAN_ALLOW_ECHO_FRAMES        = TPCANParameter(44) // Echo messages reception status within a PCAN-Channel
	PCAN_DEVICE_PART_NUMBER       = TPCANParameter(45) // Get the part number associated to a device
)

// PCAN parameter values
const (
	PCAN_PARAMETER_OFF = TPCANParameterValue(0x00) // The PCAN parameter is not set (inactive)
	PCAN_PARAMETER_ON  = TPCANParameterValue(0x01) // The PCAN parameter is set (active)
)

const (
	PCAN_FILTER_CLOSE  = TPCANFilterValue(0x00) // The PCAN filter is closed. No messages will be received
	PCAN_FILTER_OPEN   = TPCANFilterValue(0x01) // The PCAN filter is fully opened. All messages will be received
	PCAN_FILTER_CUSTOM = TPCANFilterValue(0x02) // The PCAN filter is custom configured. Only registered messages will be received
)

const (
	PCAN_CHANNEL_UNAVAILABLE = TPCANCHannelCondition(0x00)                    // The PCAN-Channel handle is illegal, or its associated hardware is not available
	PCAN_CHANNEL_AVAILABLE   = TPCANCHannelCondition(0x01)                    // The PCAN-Channel handle is available to be connected (Plug&Play Hardware: it means furthermore that the hardware is plugged-in)
	PCAN_CHANNEL_OCCUPIED    = TPCANCHannelCondition(0x02)                    // The PCAN-Channel handle is valid, and is already being used
	PCAN_CHANNEL_PCANVIEW    = PCAN_CHANNEL_AVAILABLE | PCAN_CHANNEL_OCCUPIED // The PCAN-Channel handle is already being used by a PCAN-View application, but is available to connect
)

const (
	LOG_FUNCTION_DEFAULT    = TPCANFunctionValue(0x00)   // Logs system exceptions / errors
	LOG_FUNCTION_ENTRY      = TPCANFunctionValue(0x01)   // Logs the entries to the PCAN-Basic API functions
	LOG_FUNCTION_PARAMETERS = TPCANFunctionValue(0x02)   // Logs the parameters passed to the PCAN-Basic API functions
	LOG_FUNCTION_LEAVE      = TPCANFunctionValue(0x04)   // Logs the exits from the PCAN-Basic API functions
	LOG_FUNCTION_WRITE      = TPCANFunctionValue(0x08)   // Logs the CAN messages passed to the CAN_Write function
	LOG_FUNCTION_READ       = TPCANFunctionValue(0x10)   // Logs the CAN messages received within the CAN_Read function
	LOG_FUNCTION_ALL        = TPCANFunctionValue(0xFFFF) // Logs all possible information within the PCAN-Basic API functions
)

const (
	TRACE_FILE_SINGLE    = TPCANTraceFileValue(0x00) // A single file is written until it size reaches PAN_TRACE_SIZE
	TRACE_FILE_SEGMENTED = TPCANTraceFileValue(0x01) // Traced data is distributed in several files with size PAN_TRACE_SIZE
	TRACE_FILE_DATE      = TPCANTraceFileValue(0x02) // Includes the date into the name of the trace file
	TRACE_FILE_TIME      = TPCANTraceFileValue(0x04) // Includes the start time into the name of the trace file
	TRACE_FILE_OVERWRITE = TPCANTraceFileValue(0x80) // Causes the overwriting of available traces (same name)
)

const (
	FEATURE_FD_CAPABLE    = TPCANFeatureValue(0x01) // Device supports flexible data-rate (CAN-FD)
	FEATURE_DELAY_CAPABLE = TPCANFeatureValue(0x02) // Device supports a delay between sending frames (FPGA based USB devices)
	FEATURE_IO_CAPABLE    = TPCANFeatureValue(0x04) // Device supports I/O functionality for electronic circuits (USB-Chip devices)
)

const (
	SERVICE_STATUS_STOPPED = TPCANStatusValue(0x01) // The service is not running
	SERVICE_STATUS_RUNNING = TPCANStatusValue(0x04) // The service is running
)

// Represents the type of a PCAN message
const (
	PCAN_MESSAGE_STANDARD = TPCANMessageType(0x00) // The PCAN message is a CAN Standard Frame (11-bit identifier)
	PCAN_MESSAGE_RTR      = TPCANMessageType(0x01) // The PCAN message is a CAN Remote-Transfer-Request Frame
	PCAN_MESSAGE_EXTENDED = TPCANMessageType(0x02) // The PCAN message is a CAN Extended Frame (29-bit identifier)
	PCAN_MESSAGE_FD       = TPCANMessageType(0x04) // The PCAN message represents a FD frame in terms of CiA Specs
	PCAN_MESSAGE_BRS      = TPCANMessageType(0x08) // The PCAN message represents a FD bit rate switch (CAN data at a higher bit rate)
	PCAN_MESSAGE_ESI      = TPCANMessageType(0x10) // The PCAN message represents a FD error state indicator(CAN FD transmitter was error active)
	PCAN_MESSAGE_ECHO     = TPCANMessageType(0x20) // The PCAN message represents an echo CAN Frame. This message is set if PCAN_PARAMETER_ALLOW_ECHO is enabled and a sent message is read back
	PCAN_MESSAGE_ERRFRAME = TPCANMessageType(0x40) // The PCAN message represents an error frame
	PCAN_MESSAGE_STATUS   = TPCANMessageType(0x80) // The PCAN message represents a PCAN status message
)

// Represents a PCAN filter mode
const (
	PCAN_MODE_STANDARD = TPCANMode(PCAN_MESSAGE_STANDARD) // Mode is Standard (11-bit identifier)
	PCAN_MODE_EXTENDED = TPCANMode(PCAN_MESSAGE_EXTENDED) // Mode is Extended (29-bit identifier)
)

// Represents a PCAN Baud rate register value (BTR0/BTR1 register values for the CAN controller)
const (
	PCAN_BAUD_1M   = TPCANBaudrate(0x0014) // 1 MBit/s
	PCAN_BAUD_800K = TPCANBaudrate(0x0016) // 800 KBit/s
	PCAN_BAUD_500K = TPCANBaudrate(0x001C) // 500 kBit/s
	PCAN_BAUD_250K = TPCANBaudrate(0x011C) // 250 kBit/s
	PCAN_BAUD_125K = TPCANBaudrate(0x031C) // 125 kBit/s
	PCAN_BAUD_100K = TPCANBaudrate(0x432F) // 100 kBit/s
	PCAN_BAUD_95K  = TPCANBaudrate(0xC34E) // 95,238 KBit/s
	PCAN_BAUD_83K  = TPCANBaudrate(0x852B) // 83,333 KBit/s
	PCAN_BAUD_50K  = TPCANBaudrate(0x472F) // 50 kBit/s
	PCAN_BAUD_47K  = TPCANBaudrate(0x1414) // 47,619 KBit/s
	PCAN_BAUD_33K  = TPCANBaudrate(0x8B2F) // 33,333 KBit/s
	PCAN_BAUD_20K  = TPCANBaudrate(0x532F) // 20 kBit/s
	PCAN_BAUD_10K  = TPCANBaudrate(0x672F) // 10 kBit/s
	PCAN_BAUD_5K   = TPCANBaudrate(0x7F7F) // 5 kBit/s
)

// Represents the type of PCAN (Non-PnP) hardware to be initialized
const (
	PCAN_TYPE_ISA         = TPCANType(0x01) // PCAN-ISA 82C200
	PCAN_TYPE_ISA_SJA     = TPCANType(0x09) // PCAN-ISA SJA1000
	PCAN_TYPE_ISA_PHYTEC  = TPCANType(0x04) // PHYTEC ISA
	PCAN_TYPE_DNG         = TPCANType(0x02) // PCAN-Dongle 82C200
	PCAN_TYPE_DNG_EPP     = TPCANType(0x03) // PCAN-Dongle EPP 82C200
	PCAN_TYPE_DNG_SJA     = TPCANType(0x05) // PCAN-Dongle SJA1000
	PCAN_TYPE_DNG_SJA_EPP = TPCANType(0x06) // PCAN-Dongle EPP SJA1000
)

// LookUp Parameters
const (
	LOOKUP_DEVICE_TYPE       = TPCANLookupParameter("devicetype")       // Lookup channel by Device type (see PCAN devices e.g. PCAN_USB)
	LOOKUP_DEVICE_ID         = TPCANLookupParameter("deviceid")         // Lookup channel by device id
	LOOKUP_CONTROLLER_NUMBER = TPCANLookupParameter("controllernumber") // Lookup channel by CAN controller 0-based index
	LOOKUP_IP_ADDRESS        = TPCANLookupParameter("ipaddress")        // Lookup channel by IP address (LAN channels only)
)

// Represents the configuration for a CAN bit rate
// Note:
//   - Each parameter and its value must be separated with a '='.
//   - Each pair of parameter/value must be separated using ','.
//
// Example:
//
//	f_clock=80000000,nom_brp=10,nom_tseg1=5,nom_tseg2=2,nom_sjw=1,data_brp=4,data_tseg1=7,data_tseg2=2,data_sjw=1
const (
	PCAN_BR_CLOCK       = TPCANBRParameter("f_clock")
	PCAN_BR_CLOCK_MHZ   = TPCANBRParameter("f_clock_mhz")
	PCAN_BR_NOM_BRP     = TPCANBRParameter("nom_brp")
	PCAN_BR_NOM_TSEG1   = TPCANBRParameter("nom_tseg1")
	PCAN_BR_NOM_TSEG2   = TPCANBRParameter("nom_tseg2")
	PCAN_BR_NOM_SJW     = TPCANBRParameter("nom_sjw")
	PCAN_BR_NOM_SAMPLE  = TPCANBRParameter("nom_sam")
	PCAN_BR_DATA_BRP    = TPCANBRParameter("data_brp")
	PCAN_BR_DATA_TSEG1  = TPCANBRParameter("data_tseg1")
	PCAN_BR_DATA_TSEG2  = TPCANBRParameter("data_tseg2")
	PCAN_BR_DATA_SJW    = TPCANBRParameter("data_sjw")
	PCAN_BR_DATA_SAMPLE = TPCANBRParameter("data_ssp_offset")
)
