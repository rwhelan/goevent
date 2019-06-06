# goevent

Overly simple system for eventing in Go.

TODO:
- Add tests
- Add Godoc
- Profit?

### Example

```go
package main

import (
	"fmt"
	"time"

	event "github.com/rwhelan/goevent"
)

// Example event to emit
type msg struct {
	eventType string
}

// Satisfy Event interface
func (m msg) Name() string {
	return m.eventType
}

func newMessage(t string) msg {
	return msg{eventType: t}
}

func main() {
	shortbus := event.NewEventBus()

	listenerOne, err := shortbus.Subscribe("One", event.FilterAll())
	if err != nil {
		panic(err)
	}

	listenerTwo, err := shortbus.Subscribe("Two", event.FilterEveryOther())
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			select {
			case m := <-listenerOne:
				fmt.Printf("Recv Event On Listener #1: %+v\n", m)
			case m := <-listenerTwo:
				fmt.Printf("Recv Event On Listener #2: %+v\n", m)
			}

		}
	}()

	shortbus.Publish(newMessage("Msg1"))
	shortbus.Publish(newMessage("Msg2"))
	shortbus.Publish(newMessage("Msg3"))
	shortbus.Publish(newMessage("Msg4"))
	shortbus.Publish(newMessage("Msg5"))
	shortbus.Publish(newMessage("Msg6"))

	time.Sleep(time.Millisecond * 20)
}
```

The above example should produce: 

(Notice the filter func on listener #2 at work)
```go
Recv Event On Listener #1: {eventType:Msg1}
Recv Event On Listener #2: {eventType:Msg3}
Recv Event On Listener #1: {eventType:Msg2}
Recv Event On Listener #1: {eventType:Msg3}
Recv Event On Listener #1: {eventType:Msg5}
Recv Event On Listener #2: {eventType:Msg1}
Recv Event On Listener #1: {eventType:Msg4}
Recv Event On Listener #2: {eventType:Msg5}
Recv Event On Listener #1: {eventType:Msg6}
```

