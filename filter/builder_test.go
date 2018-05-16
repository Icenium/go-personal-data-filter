package filter

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type pd struct {
	MyProp string
	Email  string
}

func TestPersonalDataFilterBuilder(t *testing.T) {
	Convey("PersonalDataFilterBuilder", t, func() {
		testRegExpString := "test-regExp"
		testRegExp := regexp.MustCompile(`\-.*`)
		i := pd{MyProp: testRegExpString, Email: "not-personal"}

		Convey("WithMask", func() {
			Convey("Should set the filter mask correctly.", func() {
				mask := `¯\_(:|)_/¯`
				f, err := NewBuilder().WithMask(mask).Build()

				So(err, ShouldBeNil)

				res := f.RemovePersonalData(i)

				So(res, ShouldResemble, pd{MyProp: testRegExpString, Email: mask})
			})
			Convey("Should not set the mask if there is builder error.", func() {
				mask := `¯\_(:|)_/¯`
				b := NewBuilder()
				b.err = errPDPropsAndAdditionalPDProps
				b = b.WithMask(mask)
				So(b.mask, ShouldEqual, "")
			})
		})

		Convey("WithRegExp", func() {
			Convey("Should use the provided regular expression.", func() {
				f, err := NewBuilder().WithRegExp(testRegExp).Build()

				So(err, ShouldBeNil)

				res := f.RemovePersonalData(i)

				So(res, ShouldResemble, pd{MyProp: "test", Email: ""})
			})
			Convey("Should fail the build if it's used after WithAdditionalRegularExpression.", func() {
				_, err := NewBuilder().
					WithAdditionalRegularExpressions("some").
					WithRegExp(testRegExp).
					Build()

				So(err, ShouldBeError, errRegExpAndAdditionalRegExp)
			})
			Convey("Should not set reg exp if there is builder error.", func() {
				b := NewBuilder()
				b.err = errPDPropsAndAdditionalPDProps
				b = b.WithRegExp(testRegExp)
				So(b.regExp, ShouldBeNil)
			})
		})

		Convey("WithAdditionalRegularExpression", func() {
			Convey("Should use the provided regular expression.", func() {
				f, err := NewBuilder().WithAdditionalRegularExpressions(testRegExp.String()).Build()

				So(err, ShouldBeNil)

				res := f.RemovePersonalData(i)

				So(res, ShouldResemble, pd{MyProp: "test", Email: ""})
			})
			Convey("Should fail the build if it's used after WithRegExp.", func() {
				_, err := NewBuilder().
					WithRegExp(testRegExp).
					WithAdditionalRegularExpressions("some").
					Build()

				So(err, ShouldBeError, errRegExpAndAdditionalRegExp)
			})
			Convey("Should not add reg exp if there is builder error.", func() {
				b := NewBuilder()
				b.err = errPDPropsAndAdditionalPDProps
				b = b.WithAdditionalRegularExpressions("some")
				So(b.additionalRegExps, ShouldHaveLength, 0)
			})
		})

		Convey("WithPersonalDataProperties", func() {
			Convey("Should use the provided personal data properties.", func() {
				f, err := NewBuilder().WithPersonalDataProperties("myprop").Build()

				So(err, ShouldBeNil)

				res := f.RemovePersonalData(i)

				So(res, ShouldResemble, pd{MyProp: "", Email: "not-personal"})
			})
			Convey("Should fail the build if it's used after WithAdditionalPersonalDataProperties.", func() {
				_, err := NewBuilder().
					WithAdditionalPersonalDataProperties("some").
					WithPersonalDataProperties("prop").
					Build()

				So(err, ShouldBeError, errPDPropsAndAdditionalPDProps)
			})
			Convey("Should not set personal data properties if there is builder error.", func() {
				b := NewBuilder()
				b.err = errPDPropsAndAdditionalPDProps
				b = b.WithPersonalDataProperties("some")
				So(b.personalDataProperties, ShouldHaveLength, 0)
			})
		})

		Convey("WithAdditionalPersonalDataProperties", func() {
			Convey("Should use the provided personal data properties.", func() {
				f, err := NewBuilder().WithAdditionalPersonalDataProperties("myprop").Build()

				So(err, ShouldBeNil)

				res := f.RemovePersonalData(i)

				So(res, ShouldResemble, pd{MyProp: "", Email: ""})
			})
			Convey("Should fail the build if it's used after WithPersonalDataProperties.", func() {
				_, err := NewBuilder().
					WithPersonalDataProperties("prop").
					WithAdditionalPersonalDataProperties("some").
					Build()

				So(err, ShouldBeError, errPDPropsAndAdditionalPDProps)
			})
			Convey("Should not add personal data properties if there is builder error.", func() {
				b := NewBuilder()
				b.err = errPDPropsAndAdditionalPDProps
				b = b.WithAdditionalPersonalDataProperties("some")
				So(b.additionalPersonalDataProperties, ShouldHaveLength, 0)
			})
		})
	})
}
