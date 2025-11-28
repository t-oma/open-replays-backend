package domain

import "time"

type Driver struct {
	Name string
}

type Team struct {
	Name    string
	Drivers []*Driver
}

type Race struct {
	Name      string
	Teams     []*Team
	StartDate time.Time
	EndDate   time.Time
}

type RaceSegment struct {
	Name      string
	StartTime time.Time
	EndTime   time.Time
	Date      time.Time
}

type Schedule struct {
	Segments []*RaceSegment
}

type Calendar struct {
	Races []*Race
}
