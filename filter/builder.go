package filter

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	// Source: https://stackoverflow.com/a/46181/4922411
	emailRegExpTemplate = `(([^<>()\[\]\.,;:\s@"]+(\.[^<>()\[\]\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))`
	// Source: https://stackoverflow.com/a/11040993/4922411
	guidRegExpTemplate = `[{(]?[0-9A-F]{8}-([0-9A-F]{4}-){3}[0-9A-F]{12}[)}]?`
	// Source: https://stackoverflow.com/a/34529037/4922411
	ipV4RegExpTemplate = `(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
	// Source: https://stackoverflow.com/a/9221063/4922411
	// nolint[:lll]
	ipV6RegExpTemplate = `((([0-9A-Fa-f]{1,4}:){7}([0-9A-Fa-f]{1,4}|:))|(([0-9A-Fa-f]{1,4}:){6}(:[0-9A-Fa-f]{1,4}|((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){5}(((:[0-9A-Fa-f]{1,4}){1,2})|:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3})|:))|(([0-9A-Fa-f]{1,4}:){4}(((:[0-9A-Fa-f]{1,4}){1,3})|((:[0-9A-Fa-f]{1,4})?:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){3}(((:[0-9A-Fa-f]{1,4}){1,4})|((:[0-9A-Fa-f]{1,4}){0,2}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){2}(((:[0-9A-Fa-f]{1,4}){1,5})|((:[0-9A-Fa-f]{1,4}){0,3}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(([0-9A-Fa-f]{1,4}:){1}(((:[0-9A-Fa-f]{1,4}){1,6})|((:[0-9A-Fa-f]{1,4}){0,4}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:))|(:(((:[0-9A-Fa-f]{1,4}){1,7})|((:[0-9A-Fa-f]{1,4}){0,5}:((25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)(\.(25[0-5]|2[0-4]\d|1\d\d|[1-9]?\d)){3}))|:)))(%.+)?`
)

var (
	personalDataProperties = []string{"email", "useremail", "user", "username", "userid", "accountid", "account", "password", "pass", "pwd", "ip", "ipaddress"}

	errRegExpAndAdditionalRegExp   = errors.New("can't use WithAdditionalRegularExpressions and WithRegExp at the same time")
	errPDPropsAndAdditionalPDProps = errors.New("can't use WithPersonalDataProperties and WithAdditionalPersonalDataProperties at the same time")
)

// PersonalDataFilterBuilder builds personal data filter
// based on the provided configuration.
type PersonalDataFilterBuilder struct {
	mask                             string
	regExp                           *regexp.Regexp
	additionalRegExps                []string
	personalDataProperties           []string
	additionalPersonalDataProperties []string
	err                              error
}

// WithMask sets the mask string which will be used to replace the personal data.
func (b *PersonalDataFilterBuilder) WithMask(mask string) *PersonalDataFilterBuilder {
	if b.err != nil {
		return b
	}

	b.mask = mask
	return b
}

// WithRegExp sets the regular expression which will be used to search for personal data.
func (b *PersonalDataFilterBuilder) WithRegExp(regExp *regexp.Regexp) *PersonalDataFilterBuilder {
	if b.err != nil {
		return b
	}

	if len(b.additionalRegExps) > 0 {
		b.err = errRegExpAndAdditionalRegExp
		return b
	}

	b.regExp = regExp
	return b
}

// WithAdditionalRegularExpressions adds more regular expressions to the default ones for
// searching for personal data.
func (b *PersonalDataFilterBuilder) WithAdditionalRegularExpressions(additionalRegExps ...string) *PersonalDataFilterBuilder {
	if b.err != nil {
		return b
	}

	if b.regExp != nil {
		b.err = errRegExpAndAdditionalRegExp
		return b
	}

	b.additionalRegExps = additionalRegExps
	return b
}

// WithPersonalDataProperties sets the personal data properties which will be used when filtering structs and maps.
func (b *PersonalDataFilterBuilder) WithPersonalDataProperties(props ...string) *PersonalDataFilterBuilder {
	if b.err != nil {
		return b
	}

	if len(b.additionalPersonalDataProperties) > 0 {
		b.err = errPDPropsAndAdditionalPDProps
		return b
	}

	b.personalDataProperties = props
	return b
}

// WithAdditionalPersonalDataProperties sets the personal data properties which will be
// added to the default ones for filtering structs and maps.
func (b *PersonalDataFilterBuilder) WithAdditionalPersonalDataProperties(props ...string) *PersonalDataFilterBuilder {
	if b.err != nil {
		return b
	}

	if len(b.personalDataProperties) > 0 {
		b.err = errPDPropsAndAdditionalPDProps
		return b
	}

	b.additionalPersonalDataProperties = props
	return b
}

// Build creates new personal data filter from the provided configuration.
func (b *PersonalDataFilterBuilder) Build() (PersonalDataFilter, error) {
	if b.err != nil {
		return nil, b.err
	}

	res := new(personalDataFilter)

	// Handle mask config.
	res.mask = b.mask

	// Handle RegExp config.
	if b.regExp != nil {
		res.personalDataRegExp = b.regExp
	} else {
		// We need the []interface{} because the fmt.Sprintf does not work with []string.
		allRegularExpressions := []interface{}{}
		allRegularExpressions = append(allRegularExpressions,
			emailRegExpTemplate,
			guidRegExpTemplate,
			ipV4RegExpTemplate,
			ipV6RegExpTemplate,
		)
		for _, v := range b.additionalRegExps {
			allRegularExpressions = append(allRegularExpressions, v)
		}

		regExpStringTemplate := strings.TrimRight(strings.Repeat("(%s)|", len(allRegularExpressions)), "|")

		res.personalDataRegExp = regexp.MustCompile(fmt.Sprintf("(?i)"+regExpStringTemplate, allRegularExpressions...))
	}

	// Handle personal data properties config.
	if len(b.personalDataProperties) > 0 {
		res.personalDataProperties = b.personalDataProperties
	} else {
		res.personalDataProperties = append(personalDataProperties, b.additionalPersonalDataProperties...)
	}

	return res, nil
}

// NewBuilder creates new personal data filter builder.
func NewBuilder() *PersonalDataFilterBuilder {
	return new(PersonalDataFilterBuilder)
}
