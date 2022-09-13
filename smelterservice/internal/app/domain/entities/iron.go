package entities

import "time"

type Iron struct {
	CreatedAt    	time.Time
	Weight			int
	Quality			int
	Type 			string
}

