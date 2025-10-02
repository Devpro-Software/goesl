// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import (
	"context"
	"fmt"
	"net"
)

type OutboundServer struct {
	net.Listener

	Addr  string `json:"address"`
	Proto string

	Conns chan *SocketConnection
}

func (s *OutboundServer) Start() error {
	Notice("Starting Freeswitch Outbound Server @ (address: %s) ...", s.Addr)

	var err error

	s.Listener, err = net.Listen(s.Proto, s.Addr)
	if err != nil {
		Error(ECouldNotStartListener, err)
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			Warning("Waiting for incoming connections ...")

			c, err := s.Accept()
			if err != nil {
				Error(EListenerConnection, err)
				cancel()
				break
			}

			ctx := context.Background()
			ctx, cancel := context.WithCancel(ctx)
			conn := &SocketConnection{
				Conn:   c,
				err:    make(chan error),
				m:      make(chan *Message),
				ctx:    ctx,
				cancel: cancel,
			}

			Notice("Got new connection from: %s", conn.OriginatorAddr())

			go conn.Handle()

			s.Conns <- conn
		}
	}()

	<-ctx.Done()
	s.Close()

	return err
}

func NewOutboundServer(addr string) (*OutboundServer, error) {
	if addr == "" {
		return nil, fmt.Errorf(EInvalidServerAddr, addr)
	}

	server := OutboundServer{
		Addr:  addr,
		Proto: "tcp",
		Conns: make(chan *SocketConnection),
	}

	return &server, nil
}
