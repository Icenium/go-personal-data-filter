package filter

// PersonalDataFilter is filter which takes care of removing personal data from all
// kinds of input.
type PersonalDataFilter interface {
	// RemovePersonalData removes the personal data from the provided input.
	RemovePersonalData(input interface{}) interface{}
}

// MatchReplacer is function which will be used to replace each match found by some
// of the registered regular expressions.
type MatchReplacer func(match string) (replaced string)

type filterTagConfig struct {
	NoFilter bool
}
