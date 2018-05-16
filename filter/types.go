package filter

// PersonalDataFilter is filter which takes care of removing personal data from all
// kinds of input.
type PersonalDataFilter interface {
	// RemovePersonalData removes the personal data from the provided input.
	RemovePersonalData(input interface{}) interface{}
}

type filterTagConfig struct {
	NoFilter bool
}
