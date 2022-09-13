package entities

import "time"

type Sword struct {
	CreatedAt    	time.Time
	Weight			int
	Quality			int
	Sharpened		bool
	Type 			string
}
