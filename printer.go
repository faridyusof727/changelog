package main

type Printer interface {
	MapData(tags []*TagInfo)
	Print()
}
