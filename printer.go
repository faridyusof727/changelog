package main

// Printer defines the interface for formatting and outputting changelog data.
type Printer interface {
	// MapData processes and maps the provided tag information for output.
	MapData(tags []*TagInfo)
	// Print outputs the formatted changelog data.
	Print(current string)
}
