package amqp

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/Azure/go-amqp/internal/encoding"
	"github.com/Azure/go-amqp/internal/frames"
	"github.com/Azure/go-amqp/internal/mocks"
	"github.com/stretchr/testify/require"
)

func TestSenderMethodsNoSend(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeUnsettled)
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	const (
		linkAddr   = "addr1"
		linkName   = "test1"
		maxMsgSize = uint64(4096)
	)
	snd, err := session.NewSender(LinkTargetAddress(linkAddr), LinkName(linkName), LinkMaxMessageSize(maxMsgSize))
	require.NoError(t, err)
	require.NotNil(t, snd)
	require.Equal(t, linkAddr, snd.Address())
	require.Equal(t, linkName, snd.LinkName())
	require.Equal(t, maxMsgSize, snd.MaxMessageSize())
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	require.NoError(t, snd.Close(ctx))
	cancel()
	require.NoError(t, client.Close())
}

func TestSenderSendOnClosed(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeUnsettled)
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender()
	require.NoError(t, err)
	require.NotNil(t, snd)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	require.NoError(t, snd.Close(ctx))
	cancel()
	// sending on a closed sender returns ErrLinkClosed
	if err = snd.Send(context.Background(), NewMessage([]byte("failed"))); !errors.Is(err, ErrLinkClosed) {
		t.Fatalf("unexpected error %T", err)
	}
	require.NoError(t, client.Close())
}

func TestSenderSendOnSessionClosed(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeUnsettled)
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender()
	require.NoError(t, err)
	require.NotNil(t, snd)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	require.NoError(t, session.Close(ctx))
	cancel()
	// sending on a closed sender returns ErrLinkClosed
	if err = snd.Send(context.Background(), NewMessage([]byte("failed"))); !errors.Is(err, ErrSessionClosed) {
		t.Fatalf("unexpected error %T", err)
	}
	require.NoError(t, client.Close())
}

func TestSenderSendOnConnClosed(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeUnsettled)
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender()
	require.NoError(t, err)
	require.NotNil(t, snd)

	require.NoError(t, client.Close())
	// sending on a closed sender returns ErrLinkClosed
	if err = snd.Send(context.Background(), NewMessage([]byte("failed"))); !errors.Is(err, ErrConnClosed) {
		t.Fatalf("unexpected error %T", err)
	}
	require.NoError(t, client.Close())
}

func TestSenderSendOnDetached(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeUnsettled)
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender()
	require.NoError(t, err)
	require.NotNil(t, snd)
	// initiate a server-side detach
	const (
		errcon  = "detaching"
		errdesc = "server side detach"
	)
	b, err := mocks.PerformDetach(0, 0, &Error{Condition: errcon, Description: errdesc})
	require.NoError(t, err)
	netConn.SendFrame(b)
	// sending on a detached link returns a DetachError
	err = snd.Send(context.Background(), NewMessage([]byte("failed")))
	var de *DetachError
	if !errors.As(err, &de) {
		t.Fatalf("unexpected error type %T", err)
	}
	require.Equal(t, encoding.ErrorCondition(errcon), de.RemoteError.Condition)
	require.Equal(t, errdesc, de.RemoteError.Description)
	require.NoError(t, client.Close())
}

func TestSenderAttachError(t *testing.T) {
	detachAck := make(chan bool)
	var enqueueFrames func(string)
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			enqueueFrames(tt.Name)
			return nil, nil
		case *frames.PerformDetach:
			// we don't need to respond to the ack
			detachAck <- true
			return nil, nil
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)
	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)

	const (
		errcon  = "cantattach"
		errdesc = "server side error"
	)

	enqueueFrames = func(n string) {
		// send an invalid attach response
		b, err := mocks.EncodeFrame(mocks.FrameAMQP, 0, &frames.PerformAttach{
			Name: n,
			Role: encoding.RoleReceiver,
		})
		require.NoError(t, err)
		netConn.SendFrame(b)
		// now follow up with a detach frame
		b, err = mocks.EncodeFrame(mocks.FrameAMQP, 0, &frames.PerformDetach{
			Error: &encoding.Error{
				Condition:   errcon,
				Description: errdesc,
			},
		})
		require.NoError(t, err)
		netConn.SendFrame(b)
	}
	snd, err := session.NewSender()
	var de *Error
	if !errors.As(err, &de) {
		t.Fatalf("unexpected error type %T", err)
	}
	require.Equal(t, encoding.ErrorCondition(errcon), de.Condition)
	require.Equal(t, errdesc, de.Description)
	require.Nil(t, snd)
	require.Equal(t, true, <-detachAck)
	require.NoError(t, client.Close())
}

func TestSenderSendMismatchedModes(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeUnsettled)
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender(LinkSenderSettle(encoding.ModeSettled))
	require.Error(t, err)
	require.Equal(t, "amqp: sender settlement mode \"settled\" requested, received \"unsettled\" from server", err.Error())
	require.Nil(t, snd)
	require.NoError(t, client.Close())
}

func TestSenderSendSuccess(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeUnsettled)
		case *frames.PerformTransfer:
			if tt.More {
				return nil, errors.New("didn't expect more to be true")
			}
			if tt.Settled {
				return nil, errors.New("didn't expect message to be settled")
			}
			if tt.MessageFormat == nil {
				return nil, errors.New("unexpected nil MessageFormat")
			}
			if !reflect.DeepEqual([]byte{0, 83, 117, 160, 4, 116, 101, 115, 116}, tt.Payload) {
				return nil, fmt.Errorf("unexpected payload %v", tt.Payload)
			}
			return mocks.PerformDisposition(encoding.RoleReceiver, 0, *tt.DeliveryID, &encoding.StateAccepted{})
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender()
	require.NoError(t, err)

	sendInitialFlowFrame(t, netConn, 0, 100)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	require.NoError(t, snd.Send(ctx, NewMessage([]byte("test"))))
	cancel()

	require.NoError(t, client.Close())
}

func TestSenderSendSettled(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeSettled)
		case *frames.PerformTransfer:
			if tt.More {
				return nil, errors.New("didn't expect more to be true")
			}
			if !tt.Settled {
				return nil, errors.New("expected message to be settled")
			}
			if !reflect.DeepEqual([]byte{0, 83, 117, 160, 4, 116, 101, 115, 116}, tt.Payload) {
				return nil, fmt.Errorf("unexpected payload %v", tt.Payload)
			}
			return nil, nil
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender(LinkSenderSettle(ModeSettled))
	require.NoError(t, err)

	sendInitialFlowFrame(t, netConn, 0, 100)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	require.NoError(t, snd.Send(ctx, NewMessage([]byte("test"))))
	cancel()

	require.NoError(t, client.Close())
}

func TestSenderSendRejected(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeUnsettled)
		case *frames.PerformTransfer:
			return mocks.PerformDisposition(encoding.RoleReceiver, 0, *tt.DeliveryID, &encoding.StateRejected{
				Error: &Error{
					Condition:   "rejected",
					Description: "didn't like it",
				},
			})
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender()
	require.NoError(t, err)

	sendInitialFlowFrame(t, netConn, 0, 100)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = snd.Send(ctx, NewMessage([]byte("test")))
	cancel()
	var asErr *Error
	if !errors.As(err, &asErr) {
		t.Fatalf("unexpected error type %T", err)
	}
	require.Equal(t, encoding.ErrorCondition("rejected"), asErr.Condition)

	require.NoError(t, client.Close())
}

func TestSenderSendDetached(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeUnsettled)
		case *frames.PerformTransfer:
			return mocks.PerformDetach(0, 0, &Error{
				Condition:   "detached",
				Description: "server exploded",
			})
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender()
	require.NoError(t, err)

	sendInitialFlowFrame(t, netConn, 0, 100)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	err = snd.Send(ctx, NewMessage([]byte("test")))
	cancel()
	var asErr *DetachError
	if !errors.As(err, &asErr) {
		t.Fatalf("unexpected error type %T", err)
	}
	require.Equal(t, encoding.ErrorCondition("detached"), asErr.RemoteError.Condition)

	require.NoError(t, client.Close())
}

func TestSenderSendTimeout(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeUnsettled)
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender()
	require.NoError(t, err)

	// no credits have been issued so the send will time out
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	require.Error(t, snd.Send(ctx, NewMessage([]byte("test"))))
	cancel()

	require.NoError(t, client.Close())
}

func TestSenderSendMsgTooBig(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			mode := encoding.ModeUnsettled
			return mocks.EncodeFrame(mocks.FrameAMQP, 0, &frames.PerformAttach{
				Name:   tt.Name,
				Handle: 0,
				Role:   encoding.RoleReceiver,
				Target: &frames.Target{
					Address:      "test",
					Durable:      encoding.DurabilityNone,
					ExpiryPolicy: encoding.ExpirySessionEnd,
				},
				SenderSettleMode: &mode,
				MaxMessageSize:   16, // really small messages only
			})
		case *frames.PerformTransfer:
			return mocks.PerformDisposition(encoding.RoleReceiver, 0, *tt.DeliveryID, &encoding.StateAccepted{})
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender()
	require.NoError(t, err)

	sendInitialFlowFrame(t, netConn, 0, 100)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	require.Error(t, snd.Send(ctx, NewMessage([]byte("test message that's too big"))))
	cancel()

	require.NoError(t, client.Close())
}

func TestSenderSendTagTooBig(t *testing.T) {
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.PerformOpen("container")
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeUnsettled)
		case *frames.PerformTransfer:
			return mocks.PerformDisposition(encoding.RoleReceiver, 0, *tt.DeliveryID, &encoding.StateAccepted{})
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender()
	require.NoError(t, err)

	sendInitialFlowFrame(t, netConn, 0, 100)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	msg := NewMessage([]byte("test"))
	// make the tag larger than max allowed of 32
	msg.DeliveryTag = make([]byte, 33)
	require.Error(t, snd.Send(ctx, msg))
	cancel()

	require.NoError(t, client.Close())
}

func TestSenderSendMultiTransfer(t *testing.T) {
	var deliveryID uint32
	transferCount := 0
	const maxReceiverFrameSize = 128
	responder := func(req frames.FrameBody) ([]byte, error) {
		switch tt := req.(type) {
		case *mocks.AMQPProto:
			return []byte{'A', 'M', 'Q', 'P', 0, 1, 0, 0}, nil
		case *frames.PerformOpen:
			return mocks.EncodeFrame(mocks.FrameAMQP, 0, &frames.PerformOpen{
				ChannelMax:   65535,
				ContainerID:  "container",
				IdleTimeout:  time.Minute,
				MaxFrameSize: maxReceiverFrameSize, // really small max frame size
			})
		case *frames.PerformBegin:
			return mocks.PerformBegin(0)
		case *frames.PerformEnd:
			return mocks.PerformEnd(0, nil)
		case *frames.PerformAttach:
			return mocks.SenderAttach(0, tt.Name, 0, encoding.ModeUnsettled)
		case *frames.PerformTransfer:
			if tt.DeliveryID != nil {
				// deliveryID is only sent on the first transfer frame for multi-frame transfers
				if transferCount != 0 {
					return nil, fmt.Errorf("unexpected DeliveryID for frame number %d", transferCount)
				}
				deliveryID = *tt.DeliveryID
			}
			if tt.MessageFormat != nil && transferCount != 0 {
				// MessageFormat is only sent on the first transfer frame for multi-frame transfers
				return nil, fmt.Errorf("unexpected MessageFormat for frame number %d", transferCount)
			} else if tt.MessageFormat == nil && transferCount == 0 {
				return nil, errors.New("unexpected nil MessageFormat")
			}
			if tt.More {
				transferCount++
				return nil, nil
			}
			return mocks.PerformDisposition(encoding.RoleReceiver, 0, deliveryID, &encoding.StateAccepted{})
		case *frames.PerformDetach:
			return mocks.PerformDetach(0, 0, nil)
		default:
			return nil, fmt.Errorf("unhandled frame %T", req)
		}
	}
	netConn := mocks.NewNetConn(responder)

	client, err := New(netConn)
	require.NoError(t, err)

	session, err := client.NewSession()
	require.NoError(t, err)
	snd, err := session.NewSender()
	require.NoError(t, err)

	sendInitialFlowFrame(t, netConn, 0, 100)

	ctx, cancel := context.WithTimeout(context.Background(), 100000*time.Millisecond)
	payload := make([]byte, maxReceiverFrameSize*4)
	for i := 0; i < maxReceiverFrameSize*4; i++ {
		payload[i] = byte(i % 256)
	}
	require.NoError(t, snd.Send(ctx, NewMessage(payload)))
	cancel()

	// split up into 8 transfers due to transfer frame header size
	require.Equal(t, 8, transferCount)

	require.NoError(t, client.Close())
}
