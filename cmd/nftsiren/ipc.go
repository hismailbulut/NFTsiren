package main

import (
	"bytes"
	"errors"
	"net"
	"net/rpc"
	"time"

	"nftsiren/pkg/log"
)

type IpcRequestArgs struct {
	MacAddress uint64
}

type IpcReplyArgs struct{}

type IpcServerHandler struct{}

func (h *IpcServerHandler) ShowWindow(args IpcRequestArgs, reply *IpcReplyArgs) error {
	if args.MacAddress != macAddress() {
		return errors.New("requests are only accepted on same machine")
	}
	// Tray must running to be able to show window
	for !TrayRunning() {
		// Sleep a bit in order to reduce cpu usage
		time.Sleep(time.Millisecond)
	}
	// Tray is running now, show window
	ShowWindowChan <- struct{}{}
	return nil
}

type IpcClient struct {
	conn   net.Conn
	client *rpc.Client
}

// Checks for a running instance of the program. Tries to connect to a server
// If there is a server, it sends SHOW and receives OK
// returns true when there is available instance and everything went good
// and the connection has to be closed and program will be terminated
func NewIpcClient(port string) (*IpcClient, error) {
	// NOTE: Timeout parameter may not be enough for tcp connection, but speeds up startup
	addr := "localhost:" + port
	conn, err := net.DialTimeout("tcp", addr, time.Second*5)
	if err != nil {
		return nil, err
	}
	return &IpcClient{
		conn:   conn,
		client: rpc.NewClient(conn),
	}, nil
}

func (client *IpcClient) ShowWindow() error {
	return client.client.Call("IpcServerHandler.ShowWindow",
		&IpcRequestArgs{MacAddress: macAddress()},
		&IpcReplyArgs{},
	)
}

func (client *IpcClient) Close() error {
	err := client.conn.Close()
	if err != nil {
		return err
	}
	return client.client.Close()
}

// Server receives SHOW and responds OK and then closes connection
type IpcServer struct {
	listener net.Listener
}

// Create a server and process incoming signals
func NewIpcServer() (*IpcServer, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, err
	}
	rpc.Register(&IpcServerHandler{})
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				if errors.Is(err, net.ErrClosed) {
					break
				}
				log.Warn().Println("Connection error:", err)
				continue
			}
			rpc.ServeConn(conn)
		}
	}()
	return &IpcServer{
		listener: listener,
	}, nil
}

func (server *IpcServer) Port() (string, error) {
	addr := server.listener.Addr().String()
	_, port, err := net.SplitHostPort(addr)
	if err != nil {
		return "", err
	}
	return port, nil
}

func (server *IpcServer) Close() error {
	return server.listener.Close()
}

func ShowExistingInstance(port string) error {
	client, err := NewIpcClient(port)
	if err != nil {
		return err
	}
	defer client.Close()
	err = client.ShowWindow()
	if err != nil {
		return err
	}
	return nil
}

// Returns the mac address of current hardware
func macAddress() uint64 {
	interfaces, err := net.Interfaces()
	if err != nil {
		return 0
	}
	for _, i := range interfaces {
		if i.Flags&net.FlagUp != 0 && !bytes.Equal(i.HardwareAddr, nil) {
			// Skip locally administered addresses
			if i.HardwareAddr[0]&2 == 2 {
				continue
			}
			var mac uint64
			for j, b := range i.HardwareAddr {
				if j >= 8 {
					break
				}
				mac <<= 8
				mac += uint64(b)
			}
			return mac
		}
	}
	return 0
}
