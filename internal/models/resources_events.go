package models

type Resources struct {
	Money int
	Area  int // Total available area
	UsedArea int
	Power int
	Water int
}

type Event struct {
	Name        string
	Description string
	Effect      string
}
