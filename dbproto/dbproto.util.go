package dbproto

func (s *State) LatestTimestamp() int64 {
	tslen := len(s.AllTimestamp)
	if tslen == 0 {
		return 0
	}
	return s.AllTimestamp[tslen-1]
}
