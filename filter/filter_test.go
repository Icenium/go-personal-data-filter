package filter

import (
	"crypto/sha256"
	"fmt"
	"reflect"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type testCase struct {
	Input    string
	Expected string
}

var (
	testPersonalDataProperties = []string{"email", "useremail", "user", "username", "userid", "accountid", "account", "password", "pass", "pwd", "ip", "ipaddress"}
	filteredString             = "*****"
)

type nestedStruct struct {
	AccountID   string
	IP          string
	NotPersonal string
}

func TestPersonalDataFilter(t *testing.T) {
	notPersonalDataString := "not-personal"

	email := "qweasd12_=+3!#$%^&@ice.cold"
	guid := "1fec999a-7e81-4bce-8b32-1b6ddd144f1b"
	ip := "192.168.0.1"
	ipV6 := "fe80::f991:38d8:27e6:8b77"

	nStruct := nestedStruct{AccountID: "personal-data", IP: ip, NotPersonal: notPersonalDataString}
	filteredNStruct := nestedStruct{AccountID: filteredString, IP: filteredString, NotPersonal: notPersonalDataString}

	slice := []nestedStruct{nStruct, nStruct}
	filteredSlice := []nestedStruct{filteredNStruct, filteredNStruct}

	dynamicSlice := []interface{}{1, email, nStruct, notPersonalDataString}
	filteredDynamicSlice := []interface{}{1, filteredString, filteredNStruct, notPersonalDataString}

	array := [4]string{email, notPersonalDataString}
	filteredArray := [4]string{filteredString, notPersonalDataString}

	dynamicArray := [4]interface{}{1, email, nStruct, notPersonalDataString}
	filteredDynamicArray := [4]interface{}{1, filteredString, filteredNStruct, notPersonalDataString}

	mapInput := map[string]string{"email": notPersonalDataString, "Email": email, "notPersonal": notPersonalDataString}
	filteredMapInput := map[string]string{"email": filteredString, "Email": filteredString, "notPersonal": notPersonalDataString}

	dynamicMapInput := map[string]interface{}{
		"email":                     notPersonalDataString,
		"Email":                     email,
		"notEmail":                  notPersonalDataString,
		"someStruct":                nStruct,
		"someNumber":                1,
		"somePointerToStruct":       &nStruct,
		"someDynamicArray":          dynamicArray,
		"somePointerToDynamicArray": &dynamicArray,
		"someDynamicSlice":          dynamicSlice,
		"somePointerToDynamicSlice": &dynamicSlice,
		"someSlice":                 slice,
		"somePointerToSlice":        &slice,
		"someMap":                   &mapInput,
		"somePointerToMap":          &mapInput,
	}
	filteredDynamicMapInput := map[string]interface{}{
		"email":                     filteredString,
		"Email":                     filteredString,
		"notEmail":                  notPersonalDataString,
		"someStruct":                filteredNStruct,
		"someNumber":                1,
		"somePointerToStruct":       &filteredNStruct,
		"someDynamicArray":          filteredDynamicArray,
		"somePointerToDynamicArray": &filteredDynamicArray,
		"someDynamicSlice":          filteredDynamicSlice,
		"somePointerToDynamicSlice": &filteredDynamicSlice,
		"someSlice":                 filteredSlice,
		"somePointerToSlice":        &filteredSlice,
		"someMap":                   &filteredMapInput,
		"somePointerToMap":          &filteredMapInput,
	}

	intKeyMapInput := map[int]string{1: notPersonalDataString, 2: email, 3: notPersonalDataString}
	filteredIntKeyMapInput := map[int]string{1: notPersonalDataString, 2: filteredString, 3: notPersonalDataString}

	structKeyMap := map[nestedStruct]string{nStruct: email, filteredNStruct: notPersonalDataString}
	filteredStructKeyMap := map[nestedStruct]string{nStruct: filteredString, filteredNStruct: notPersonalDataString}

	Convey("RemovePersonalData", t, func() {
		filter, err := NewBuilder().
			SetMask(filteredString).
			Build()
		if err != nil {
			panic(err)
		}

		Convey("Strings", func() {
			Convey("Should filter string", func() {
				result := filter.RemovePersonalData(email)
				So(result, ShouldEqual, filteredString)
			})

			Convey("Should filter *string", func() {
				pointerResult := filter.RemovePersonalData(&email)
				So(*pointerResult.(*string), ShouldEqual, filteredString)
			})

			Convey("Should hide emails", func() {
				testCases := []testCase{
					createTestCase("%s", email),
					createTestCase("text %s text", email),
					createTestCase("%s text", email),
					createTestCase("text %s", email),
					createTestCase("%s %s", email, email),
					createTestCase("%s text %s", email, email),
					createTestCase("text %s %s", email, email),
					createTestCase("%s %s text", email, email),
				}

				checkTestCases(filter, testCases)
			})
			Convey("Should hide GUIDs", func() {
				testCases := []testCase{
					createTestCase("%s", guid),
					createTestCase("text%stext", guid),
					createTestCase("%stext", guid),
					createTestCase("text%s", guid),
					createTestCase("%s%s", guid, guid),
					createTestCase("%stext%s", guid, guid),
					createTestCase("text%s%s", guid, guid),
					createTestCase("%s%stext", guid, guid),
				}

				checkTestCases(filter, testCases)
			})
			Convey("Should hide v4 and v6 IPs", func() {
				testCases := []testCase{
					createTestCase("%s", ip),
					createTestCase("%s", ipV6),
					createTestCase("text%stext", ip),
					createTestCase("text%stext", ipV6),
					createTestCase("%stext", ip),
					createTestCase("%stext", ipV6),
					createTestCase("text%s", ip),
					createTestCase("text%s", ipV6),
					createTestCase("%s%s", ip, ip),
					createTestCase("%s%s", ipV6, ipV6),
				}

				checkTestCases(filter, testCases)
			})
			Convey("Should hide emails and GUIDs", func() {
				testCases := []testCase{
					createTestCase("%s%s", email, guid),
					createTestCase("text%stext %s", guid, email),
					createTestCase("%s %s %s %s", guid, email, guid, email),
					createTestCase("%s %stext%s %s", email, guid, guid, email),
				}

				checkTestCases(filter, testCases)
			})
			Convey("Should not mistake parts of longs strings for GUIDs", func() {
				checkTestCases(filter, []testCase{
					{Input: "487818704899480c907e2c0549664116", Expected: "487818704899480c907e2c0549664116"},
				})
			})
			Convey("Should hash the personal data when the default match replacer is provided.", func() {
				testCases := []testCase{
					createTestCaseWithReplacer(getHash, "%s%s", email, guid),
					createTestCaseWithReplacer(getHash, "text%stext %s", guid, email),
					createTestCaseWithReplacer(getHash, "%s %s %s %s", guid, email, guid, email),
					createTestCaseWithReplacer(getHash, "%s %stext%s %s", email, guid, guid, email),
					createTestCaseWithReplacer(getHash, "text%stext", ip),
					createTestCaseWithReplacer(getHash, "text%stext", ipV6),
				}

				hashFilter, _ := NewBuilder().UseDefaultMatchFilterFunc().Build()
				checkTestCases(hashFilter, testCases)
			})
		})

		Convey("Should filter slices", func() {
			result := filter.RemovePersonalData(slice)
			So(result, ShouldResemble, filteredSlice)

			pointerResult := filter.RemovePersonalData(&slice)
			So(pointerResult, ShouldResemble, &filteredSlice)

			result = filter.RemovePersonalData(dynamicSlice)
			So(result, ShouldResemble, filteredDynamicSlice)

			pointerResult = filter.RemovePersonalData(&dynamicSlice)
			So(pointerResult, ShouldResemble, &filteredDynamicSlice)
		})

		Convey("Should filter arrays", func() {
			result := filter.RemovePersonalData(array)
			So(result, ShouldResemble, filteredArray)

			pointerResult := filter.RemovePersonalData(&array)
			So(pointerResult, ShouldResemble, &filteredArray)

			result = filter.RemovePersonalData(dynamicArray)
			So(result, ShouldResemble, filteredDynamicArray)

			pointerResult = filter.RemovePersonalData(&dynamicArray)
			So(pointerResult, ShouldResemble, &filteredDynamicArray)
		})

		Convey("Should filter maps", func() {
			result := filter.RemovePersonalData(mapInput)
			So(result, ShouldResemble, filteredMapInput)

			pointerResult := filter.RemovePersonalData(&mapInput)
			So(pointerResult, ShouldResemble, &filteredMapInput)

			result = filter.RemovePersonalData(dynamicMapInput)
			So(result, ShouldResemble, filteredDynamicMapInput)

			pointerResult = filter.RemovePersonalData(&dynamicMapInput)
			So(pointerResult, ShouldResemble, &filteredDynamicMapInput)

			result = filter.RemovePersonalData(intKeyMapInput)
			So(result, ShouldResemble, filteredIntKeyMapInput)

			pointerResult = filter.RemovePersonalData(structKeyMap)
			So(pointerResult, ShouldResemble, filteredStructKeyMap)
		})

		Convey("Should filter structs", func() {
			type testStruct struct {
				Email               string
				NotPersonal         string
				Number              int
				Slice               []interface{}
				PointerEmail        *string
				PointerSlice        *[]interface{}
				Nested              nestedStruct
				PointerNested       *nestedStruct
				Array               [4]interface{}
				PointerArray        *[4]interface{}
				Map                 map[string]string
				PointerToMap        *map[string]string
				DynamicMap          map[string]interface{}
				PointerToDynamicMap *map[string]interface{}
			}

			input := testStruct{
				Email:               email,
				Nested:              nStruct,
				Number:              5,
				PointerEmail:        &email,
				PointerNested:       &nStruct,
				PointerSlice:        &dynamicSlice,
				Slice:               dynamicSlice,
				NotPersonal:         notPersonalDataString,
				Array:               dynamicArray,
				PointerArray:        &dynamicArray,
				Map:                 mapInput,
				PointerToMap:        &mapInput,
				DynamicMap:          dynamicMapInput,
				PointerToDynamicMap: &dynamicMapInput,
			}

			expected := testStruct{
				Email:               filteredString,
				Nested:              filteredNStruct,
				Number:              5,
				PointerEmail:        &filteredString,
				PointerNested:       &filteredNStruct,
				PointerSlice:        &filteredDynamicSlice,
				Slice:               filteredDynamicSlice,
				NotPersonal:         notPersonalDataString,
				Array:               filteredDynamicArray,
				PointerArray:        &filteredDynamicArray,
				Map:                 filteredMapInput,
				PointerToMap:        &filteredMapInput,
				DynamicMap:          filteredDynamicMapInput,
				PointerToDynamicMap: &filteredDynamicMapInput,
			}

			result := filter.RemovePersonalData(input)
			So(result, ShouldResemble, expected)

			result = filter.RemovePersonalData(&input)
			So(result, ShouldResemble, &expected)
		})

		Convey("Should handle nil input correctly", func() {
			result := filter.RemovePersonalData(nil)
			So(result, ShouldBeNil)
		})

		Convey("Should handle zero values correctly", func() {
			var arr [4]string
			var slc []string
			var m map[string]string
			var n int
			var str nestedStruct

			testItems := []interface{}{arr, slc, m, n, str}

			for _, testItem := range testItems {
				result := filter.RemovePersonalData(testItem)
				So(result, ShouldResemble, reflect.Zero(reflect.TypeOf(testItem)).Interface())
			}
		})

		Convey(fmt.Sprintf("Should filter %v properties", testPersonalDataProperties), func() {
			personalDataStructFields := []reflect.StructField{}
			for _, p := range testPersonalDataProperties {
				personalDataStructFields = append(personalDataStructFields, reflect.StructField{Name: strings.ToUpper(p), Type: reflect.TypeOf("")})
			}

			personalDataStruct := reflect.StructOf(personalDataStructFields)

			personalDataInstance := reflect.New(personalDataStruct)
			fillPersonalDataStruct(personalDataStruct, personalDataInstance, notPersonalDataString)

			expected := reflect.New(personalDataStruct)
			fillPersonalDataStruct(personalDataStruct, expected, filteredString)

			result := filter.RemovePersonalData(personalDataInstance.Interface())
			So(result, ShouldResemble, expected.Interface())
		})

		Convey("Should set zero values for unexported fields", func() {
			type unexported struct {
				un1 nestedStruct
				un2 int
				un3 map[string]string
			}

			input := unexported{un1: nStruct, un2: 5, un3: mapInput}

			result := filter.RemovePersonalData(input)
			So(result, ShouldResemble, unexported{})
		})

		Convey("Should handle ebedded types", func() {
			type Embedded struct {
				Email    string
				NotEmail string
			}

			type master struct {
				Embedded
				AccountID  string
				NotPrivate string
			}

			input := master{Embedded: Embedded{Email: email, NotEmail: notPersonalDataString}, AccountID: "id", NotPrivate: notPersonalDataString}
			expected := master{Embedded: Embedded{Email: filteredString, NotEmail: notPersonalDataString}, AccountID: filteredString, NotPrivate: notPersonalDataString}
			result := filter.RemovePersonalData(input)
			So(result, ShouldResemble, expected)
		})

		Convey("Should respect field tag configuration", func() {
			type tags struct {
				Email string `pdfilter:"nofilter"`
			}

			input := tags{Email: email}
			expected := tags{Email: email}
			result := filter.RemovePersonalData(input)
			So(result, ShouldResemble, expected)
		})
	})
}

func getHash(input string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(input)))
}

func fillPersonalDataStruct(t reflect.Type, str reflect.Value, value string) {
	for i := 0; i < t.NumField(); i++ {
		f := str.Elem().Field(i)
		f.SetString(value)
	}
}

func createTestCase(template string, args ...interface{}) testCase {
	expectedArgs := make([]interface{}, len(args))
	for i := range args {
		expectedArgs[i] = filteredString
	}

	return testCase{
		Input:    fmt.Sprintf(template, args...),
		Expected: fmt.Sprintf(template, expectedArgs...),
	}
}

func createTestCaseWithReplacer(replacer MatchFilterFunc, template string, args ...interface{}) testCase {
	expectedArgs := make([]interface{}, len(args))
	for i := range args {
		expectedArgs[i] = replacer(args[i].(string))
	}

	return testCase{
		Input:    fmt.Sprintf(template, args...),
		Expected: fmt.Sprintf(template, expectedArgs...),
	}
}

func checkTestCases(filter PersonalDataFilter, testCases []testCase) {
	for _, tc := range testCases {
		res := filter.RemovePersonalData(tc.Input)
		So(res, ShouldEqual, tc.Expected)
	}
}
