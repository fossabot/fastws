package ws

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"
)

func TestMessageSerialization(t *testing.T) {
	resetCounter()

	var tests = []struct {
		msg        Message // in
		serialized []byte  // out
	}{
		{ // 0
			msg: Message{
				Namespace: "default",
				isConnect: true,
			},
			serialized: []byte("0;default;;0;1;0;"),
		},
		{ // 1
			msg: Message{
				Namespace: "default",
				Body:      []byte("some id"),
				isConnect: true,
			},
			serialized: []byte("0;default;;0;1;0;some id"),
		},
		{ // 2
			msg: Message{
				Namespace:    "default",
				isDisconnect: true,
			},
			serialized: []byte("0;default;;0;0;1;"),
		},
		{ // 3
			msg: Message{
				Namespace: "default",
				Event:     "chat",
				Body:      []byte("text"),
			},
			serialized: []byte("0;default;chat;0;0;0;text"),
		},
		{ // 4
			msg: Message{
				Namespace: "default",
				Event:     "chat",
				Err:       fmt.Errorf("error message"),
				isError:   true,
			},
			serialized: []byte("0;default;chat;1;0;0;error message"),
		},
		{ // 5
			msg: Message{
				Namespace: "default",
				Event:     "chat",
				Body:      []byte("a body with many ; delimeters; like that;"),
			},
			serialized: []byte("0;default;chat;0;0;0;a body with many ; delimeters; like that;"),
		},
		{ // 6
			msg: Message{
				Namespace: "",
				Event:     "chat",
				Err:       fmt.Errorf("an error message with many ; delimeters; like that;"),
				isError:   true,
			},
			serialized: []byte("0;;chat;1;0;0;an error message with many ; delimeters; like that;"),
		},
		{ // 7
			msg: Message{
				Namespace: "default",
				Event:     "chat",
				Body:      []byte("body"),
				wait:      incrementCounter(),
			},
			serialized: []byte("1;default;chat;0;0;0;body"),
		},
	}

	for i, tt := range tests {
		got := serializeMessage(nil, tt.msg)
		if !bytes.Equal(got, tt.serialized) {
			t.Fatalf("[%d] serialize: expected %s but got %s", i, tt.serialized, got)
		}

		msg := deserializeMessage(nil, got)
		if !reflect.DeepEqual(msg, tt.msg) {
			t.Fatalf("[%d] deserialize: expected\n%#+v but got\n%#+v", i, tt.msg, msg)
		}
	}

	msg := deserializeMessage(nil, []byte("default;chat;"))
	if !msg.isInvalid {
		t.Fatalf("expected message to be invalid but it seems that it is a valid one")
	}
}
