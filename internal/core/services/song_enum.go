package services

type SongType struct {
	Opening    int
	Ending     int
	Soundtrack int
}

func NewSongType() SongType {
	return SongType{
		Opening:    1,
		Ending:     2,
		Soundtrack: 3,
	}
}
