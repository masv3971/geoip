package model

// LoginEventsSortByTimestamp implements sort.Interface based on Timestamp field
type LoginEventsSortByTimestamp LoginEvents

// Len sort timestamp
func (s LoginEventsSortByTimestamp) Len() int { return len(s) }

// Less sort timestamp
func (s LoginEventsSortByTimestamp) Less(i, j int) bool {
	return s[i].Timestamp.After(s[j].Timestamp)
}

// Swap sorts timestamp
func (s LoginEventsSortByTimestamp) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// StatsOverviewDocsSortByOccurrences implements sort.Interface based on occurrences of loginEvents
type StatsOverviewDocsSortByOccurrences StatsOverviewDocs

// Len sort by occurrences of loginEvent
func (s StatsOverviewDocsSortByOccurrences) Len() int { return len(s) }

// Less sort by occurrences of loginEvent
func (s StatsOverviewDocsSortByOccurrences) Less(i, j int) bool {
	return s[i].NumbnerOfLoginEvents > s[j].NumbnerOfLoginEvents
}

// Swap sorts by occurrences of loginEvent
func (s StatsOverviewDocsSortByOccurrences) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
