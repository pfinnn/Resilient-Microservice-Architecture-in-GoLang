package entities

import "git.haw-hamburg.de/acm746/resilient-microservice/internal/pkg/util"

type Ore struct {
	CreatedAt    	util.Time
	Weight			int
	Quality			int
	Type 			string
}
