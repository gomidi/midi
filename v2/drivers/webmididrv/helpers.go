package webmididrv

import (
	"syscall/js"

	"gitlab.com/gomidi/midi/v2/drivers"
)

func log(s string) {
	jsConsole := js.Global().Get("console")

	if !jsConsole.Truthy() {
		return
	}

	jsConsole.Call("log", js.ValueOf(s))
}

type inPorts []drivers.In

func (i inPorts) Len() int {
	return len(i)
}

func (i inPorts) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func (i inPorts) Less(a, b int) bool {
	return i[a].Number() < i[b].Number()
}

type outPorts []drivers.Out

func (i outPorts) Len() int {
	return len(i)
}

func (i outPorts) Swap(a, b int) {
	i[a], i[b] = i[b], i[a]
}

func (i outPorts) Less(a, b int) bool {
	return i[a].Number() < i[b].Number()
}
