package constants

const (
	LogFieldRequestID  = "reqID"
	LogFieldDestAddr   = "dstAddr"
	LogFieldLocalAddr  = "localAddr"
	LogFieldRemotePort = "remotePort"
	LogFieldProtocol   = "protocol"
)

const (
	ServerName = "krelay-server"
	ServerPort = 9527
)

const (
	UDPBufferSize = 65536 + 2
	TCPBufferSize = 32768
)

const PortForwardProtocolV1Name = "portforward.k8s.io"

const (
	ProtocolTCP = "tcp"
	ProtocolUDP = "udp"
)

const ServerImage = "eu.gcr.io/soundbadger-staging/krelay:0.0.1"

var ClientVersion = "unknown"
