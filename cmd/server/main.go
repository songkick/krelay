package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"time"

	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	"github.com/knight42/krelay/pkg/constants"
	"github.com/knight42/krelay/pkg/xnet"
)

type options struct {
	connectTimeout time.Duration
}

func (o *options) run(ctx context.Context) error {
	tcpListener, err := net.Listen(constants.ProtocolTCP, fmt.Sprintf("0.0.0.0:%d", constants.ServerPort))
	if err != nil {
		return err
	}
	defer tcpListener.Close()

	dialer := net.Dialer{Timeout: o.connectTimeout}

	for {
		c, err := tcpListener.Accept()
		if err != nil {
			klog.ErrorS(err, "Fail to accept connection")
			continue
		}
		go handleConn(ctx, c.(*net.TCPConn), &dialer)
	}
}

func handleConn(ctx context.Context, c *net.TCPConn, dialer *net.Dialer) {
	defer c.Close()

	hdr := xnet.Header{}
	err := hdr.FromReader(c)
	if err != nil {
		klog.ErrorS(err, "Fail to read header")
		return
	}

	dstAddr := xnet.JoinHostPort(hdr.Addr.String(), hdr.Port)

	switch hdr.Protocol {
	case xnet.ProtocolTCP:
		upstreamConn, err := dialer.DialContext(ctx, constants.ProtocolTCP, dstAddr)
		if err != nil {
			klog.ErrorS(err, "Fail to create tcp connection", constants.LogFieldRequestID, hdr.RequestID.String(), constants.LogFieldDestAddr, dstAddr)
			return
		}
		klog.InfoS("Start proxy tcp request", constants.LogFieldRequestID, hdr.RequestID.String(), constants.LogFieldDestAddr, dstAddr)
		xnet.ProxyTCP(hdr.RequestID.String(), c, upstreamConn.(*net.TCPConn))

	case xnet.ProtocolUDP:
		upstreamConn, err := dialer.DialContext(ctx, constants.ProtocolUDP, dstAddr)
		if err != nil {
			klog.ErrorS(err, "Fail to create udp connection", constants.LogFieldRequestID, hdr.RequestID.String(), constants.LogFieldDestAddr, dstAddr)
			return
		}
		klog.InfoS("Start proxy udp request", constants.LogFieldRequestID, hdr.RequestID.String(), constants.LogFieldDestAddr, dstAddr)
		udpConn := &xnet.UDPConn{UDPConn: upstreamConn.(*net.UDPConn)}
		xnet.ProxyUDP(hdr.RequestID.String(), c, udpConn)

	default:
		klog.InfoS("Unknown protocol", constants.LogFieldRequestID, hdr.RequestID.String(), constants.LogFieldDestAddr, dstAddr, constants.LogFieldProtocol, hdr.Protocol)
	}
}

func main() {
	klog.InitFlags(nil)
	o := options{}
	c := cobra.Command{
		Use: constants.ServerName,
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			return o.run(context.Background())
		},
		SilenceUsage: true,
	}
	c.Flags().AddGoFlagSet(flag.CommandLine)
	c.Flags().DurationVar(&o.connectTimeout, "connect-timeout", time.Second*10, "Timeout for connecting to upstream")
	_ = c.Execute()
}
