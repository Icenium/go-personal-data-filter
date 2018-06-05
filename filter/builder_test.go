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

		Convey("SetMask", func() {
			Convey("Should set the filter mask correctly.", func() {
				mask := `¯\_(:|)_/¯`
				f, err := NewBuilder().SetMask(mask).Build()

				So(err, ShouldBeNil)

				res := f.RemovePersonalData(i)

				So(res, ShouldResemble, pd{MyProp: testRegExpString, Email: mask})
			})
			Convey("Should not set the mask if there is builder error.", func() {
				mask := `¯\_(:|)_/¯`
				b := NewBuilder()
				b.err = errPDPropsAndAdditionalPDProps
				b = b.SetMask(mask)
				So(b.mask, ShouldEqual, "")
			})
		})

		Convey("SetRegExp", func() {
			Convey("Should use the provided regular expression.", func() {
				f, err := NewBuilder().SetRegExp(testRegExp).Build()

				So(err, ShouldBeNil)

				res := f.RemovePersonalData(i)

				So(res, ShouldResemble, pd{MyProp: "test", Email: ""})
			})
			Convey("Should fail the build if it's used after WithAdditionalRegularExpression.", func() {
				_, err := NewBuilder().
					AddRegularExpressions("some").
					SetRegExp(testRegExp).
					Build()

				So(err, ShouldBeError, errRegExpAndAdditionalRegExp)
			})
			Convey("Should not set reg exp if there is builder error.", func() {
				b := NewBuilder()
				b.err = errPDPropsAndAdditionalPDProps
				b = b.SetRegExp(testRegExp)
				So(b.regExp, ShouldBeNil)
			})
		})

		Convey("AddRegularExpressions", func() {
			Convey("Should use the provided regular expression.", func() {
				f, err := NewBuilder().AddRegularExpressions(testRegExp.String()).Build()

				So(err, ShouldBeNil)

				res := f.RemovePersonalData(i)

				So(res, ShouldResemble, pd{MyProp: "test", Email: ""})
			})
			Convey("Should fail the build if it's used after WithRegExp.", func() {
				_, err := NewBuilder().
					SetRegExp(testRegExp).
					AddRegularExpressions("some").
					Build()

				So(err, ShouldBeError, errRegExpAndAdditionalRegExp)
			})
			Convey("Should not add reg exp if there is builder error.", func() {
				b := NewBuilder()
				b.err = errPDPropsAndAdditionalPDProps
				b = b.AddRegularExpressions("some")
				So(b.additionalRegExps, ShouldHaveLength, 0)
			})
		})

		Convey("SetPersonalDataProperties", func() {
			Convey("Should use the provided personal data properties.", func() {
				f, err := NewBuilder().SetPersonalDataProperties("myprop").Build()

				So(err, ShouldBeNil)

				res := f.RemovePersonalData(i)

				So(res, ShouldResemble, pd{MyProp: "", Email: "not-personal"})
			})
			Convey("Should fail the build if it's used after WithAdditionalPersonalDataProperties.", func() {
				_, err := NewBuilder().
					AddPersonalDataProperties("some").
					SetPersonalDataProperties("prop").
					Build()

				So(err, ShouldBeError, errPDPropsAndAdditionalPDProps)
			})
			Convey("Should not set personal data properties if there is builder error.", func() {
				b := NewBuilder()
				b.err = errPDPropsAndAdditionalPDProps
				b = b.SetPersonalDataProperties("some")
				So(b.personalDataProperties, ShouldHaveLength, 0)
			})
		})

		Convey("AddPersonalDataProperties", func() {
			Convey("Should use the provided personal data properties.", func() {
				f, err := NewBuilder().AddPersonalDataProperties("myprop").Build()

				So(err, ShouldBeNil)

				res := f.RemovePersonalData(i)

				So(res, ShouldResemble, pd{MyProp: "", Email: ""})
			})
			Convey("Should fail the build if it's used after WithPersonalDataProperties.", func() {
				_, err := NewBuilder().
					SetPersonalDataProperties("prop").
					AddPersonalDataProperties("some").
					Build()

				So(err, ShouldBeError, errPDPropsAndAdditionalPDProps)
			})
			Convey("Should not add personal data properties if there is builder error.", func() {
				b := NewBuilder()
				b.err = errPDPropsAndAdditionalPDProps
				b = b.AddPersonalDataProperties("some")
				So(b.additionalPersonalDataProperties, ShouldHaveLength, 0)
			})
		})

		Convey("UseDefaultMatchFilterFunc", func() {
			Convey("Should set the default match replacer.", func() {
				f, err := NewBuilder().UseDefaultMatchFilterFunc().Build()

				So(err, ShouldBeNil)

				res := f.RemovePersonalData("email@mail.com")

				So(res, ShouldResemble, "3d13579f08e876d2d2d94da15ea657fb39795dd2a59e3378c9e58c4f4b0d053b")
			})
		})

		Convey("SetMatchFilterFunc", func() {
			Convey("Should set the custom match replacer.", func() {
				expected := "expected"
				f, err := NewBuilder().SetMatchFilterFunc(func(string) string { return expected }).Build()

				So(err, ShouldBeNil)

				res := f.RemovePersonalData("email@mail.com")

				So(res, ShouldResemble, expected)
			})
			Convey("Should not add the custom match replacer if there is builder error.", func() {
				b := NewBuilder()
				b.err = errPDPropsAndAdditionalPDProps
				b = b.SetMatchFilterFunc(func(string) string { return "" })
				So(b.matchFilterFunc, ShouldBeNil)
			})
		})
	})
}
