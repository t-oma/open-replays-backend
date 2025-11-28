package domain

type Repository interface {
	GetAllRaces() ([]*Race, error)
	GetRaceByID(id string) (*Race, error)
	GetRaceSegmentsByRaceID(id string) ([]*RaceSegment, error)
	GetScheduleByRaceID(id string) (*Schedule, error)
	GetCalendarByRaceID(id string) (*Calendar, error)
}
