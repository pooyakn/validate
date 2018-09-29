package validate

import (
	"encoding/json"
	"github.com/gookit/filter"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"unicode/utf8"
)

// Basic regular expressions for validating strings.(it is from package "asaskevich/govalidator")
const (
	Email             string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	UUID3             string = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	UUID4             string = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	UUID5             string = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	UUID              string = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
	Int               string = "^(?:[-+]?(?:0|[1-9][0-9]*))$"
	Float             string = "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
	Hexadecimal       string = "^[0-9a-fA-F]+$"
	RGBColor          string = "^rgb\\(\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*,\\s*(0|[1-9]\\d?|1\\d\\d?|2[0-4]\\d|25[0-5])\\s*\\)$"
	MultiByte         string = "[^\x00-\x7F]"
	FullWidth         string = "[^\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
	HalfWidth         string = "[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]"
	Base64            string = "^(?:[A-Za-z0-9+\\/]{4})*(?:[A-Za-z0-9+\\/]{2}==|[A-Za-z0-9+\\/]{3}=|[A-Za-z0-9+\\/]{4})$"
	Latitude          string = "^[-+]?([1-8]?\\d(\\.\\d+)?|90(\\.0+)?)$"
	Longitude         string = "^[-+]?(180(\\.0+)?|((1[0-7]\\d)|([1-9]?\\d))(\\.\\d+)?)$"
	DNSName           string = `^([a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62}){1}(\.[a-zA-Z0-9_]{1}[a-zA-Z0-9_-]{0,62})*[\._]?$`
	IP                string = `(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`
	URLSchema         string = `((ftp|tcp|udp|wss?|https?):\/\/)`
	URLUsername       string = `(\S+(:\S*)?@)`
	URLPath           string = `((\/|\?|#)[^\s]*)`
	URLPort           string = `(:(\d{1,5}))`
	URLIP             string = `([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))`
	URLSubdomain      string = `((www\.)|([a-zA-Z0-9]+([-_\.]?[a-zA-Z0-9])*[a-zA-Z0-9]\.[a-zA-Z0-9]+))`
	URL                      = `^` + URLSchema + `?` + URLUsername + `?` + `((` + URLIP + `|(\[` + IP + `\])|(([a-zA-Z0-9]([a-zA-Z0-9-_]+)?[a-zA-Z0-9]([-\.][a-zA-Z0-9]+)*)|(` + URLSubdomain + `?))?(([a-zA-Z\x{00a1}-\x{ffff}0-9]+-?-?)*[a-zA-Z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-zA-Z\x{00a1}-\x{ffff}]{1,}))?))\.?` + URLPort + `?` + URLPath + `?$`
	WinPath           string = `^[a-zA-Z]:\\(?:[^\\/:*?"<>|\r\n]+\\)*[^\\/:*?"<>|\r\n]*$`
	UnixPath          string = `^(/[^/\x00]*)+/?$`
	Semver            string = "^v?(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)\\.(?:0|[1-9]\\d*)(-(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*)(\\.(0|[1-9]\\d*|\\d*[a-zA-Z-][0-9a-zA-Z-]*))*)?(\\+[0-9a-zA-Z-]+(\\.[0-9a-zA-Z-]+)*)?$"
	hasLowerCase      string = ".*[[:lower:]]"
	hasUpperCase      string = ".*[[:upper:]]"
	hasWhitespace     string = ".*[[:space:]]"
	hasWhitespaceOnly string = "^[[:space:]]+$"
)

// some string regexp. (it is from package "asaskevich/govalidator")
var (
	rxUser              = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+/=?^_`{|}~.-]+$")
	rxHostname          = regexp.MustCompile("^[^\\s]+\\.[^\\s]+$")
	rxUserDot           = regexp.MustCompile("(^[.]{1})|([.]{1}$)|([.]{2,})")
	rxEmail             = regexp.MustCompile(Email)
	rxISBN10            = regexp.MustCompile("^(?:[0-9]{9}X|[0-9]{10})$")
	rxISBN13            = regexp.MustCompile("^(?:[0-9]{13})$")
	rxUUID3             = regexp.MustCompile(UUID3)
	rxUUID4             = regexp.MustCompile(UUID4)
	rxUUID5             = regexp.MustCompile(UUID5)
	rxUUID              = regexp.MustCompile(UUID)
	rxAlpha             = regexp.MustCompile("^[a-zA-Z]+$")
	rxAlphaNum          = regexp.MustCompile("^[a-zA-Z0-9]+$")
	rxNumber            = regexp.MustCompile("^[0-9]+$")
	rxInt               = regexp.MustCompile(Int)
	rxFloat             = regexp.MustCompile(Float)
	rxHexadecimal       = regexp.MustCompile("^[0-9a-fA-F]+$")
	rxHexColor          = regexp.MustCompile("^#?([0-9a-fA-F]{3}|[0-9a-fA-F]{6})$")
	rxRGBColor          = regexp.MustCompile(RGBColor)
	rxASCII             = regexp.MustCompile("^[\x00-\x7F]+$")
	rxPrintableASCII    = regexp.MustCompile("^[\x20-\x7E]+$")
	rxMultiByte         = regexp.MustCompile("[^\x00-\x7F]")
	rxFullWidth         = regexp.MustCompile(FullWidth)
	rxHalfWidth         = regexp.MustCompile(HalfWidth)
	rxBase64            = regexp.MustCompile(Base64)
	rxDataURI           = regexp.MustCompile("^data:.+\\/(.+);base64$")
	rxLatitude          = regexp.MustCompile(Latitude)
	rxLongitude         = regexp.MustCompile(Longitude)
	rxDNSName           = regexp.MustCompile(DNSName)
	rxIP                = regexp.MustCompile(IP)
	rxURL               = regexp.MustCompile(URL)
	rxSSN               = regexp.MustCompile(`^\d{3}[- ]?\d{2}[- ]?\d{4}$`)
	rxWinPath           = regexp.MustCompile(WinPath)
	rxUnixPath          = regexp.MustCompile(UnixPath)
	rxSemver            = regexp.MustCompile(Semver)
	rxHasLowerCase      = regexp.MustCompile(hasLowerCase)
	rxHasUpperCase      = regexp.MustCompile(hasUpperCase)
	rxHasWhitespace     = regexp.MustCompile(hasWhitespace)
	rxHasWhitespaceOnly = regexp.MustCompile(hasWhitespaceOnly)
)

// some validator alias name
var validatorAliases = map[string]string{
	// alias -> real name
	"in":    "enum",
	"num":   "number",
	"range": "between",
	// type
	"int":     "isInt",
	"uint":    "isUint",
	"bool":    "isBool",
	"float":   "isFloat",
	"map":     "isMap",
	"ints":    "isInts", // []int
	"str":     "isString",
	"string":  "isString",
	"strings": "isStrings", // []string
	"arr":     "isArray",
	"array":   "isArray",
	"slice":   "isSlice",
	// val
	"regex": "regexp",
	"eq":    "isEqual",
	"equal": "isEqual",
	"intEq": "intEqual",
	"ne":    "notEqual",
	"notEq": "notEqual",
	// int
	"lte": "max",
	"gte": "min",
	// len
	"len":      "Length",
	"lenEq":    "Length",
	"lengthEq": "Length",
	"minLen":   "minLength",
	"maxLen":   "maxLength",
	"minSize":  "minLength",
	"maxSize":  "maxLength",
	// string rune length
	"strLength":  "stringLength",
	"runeLength": "stringLength",
	// string
	"ip":        "isIP",
	"ipv4":      "isIPv4",
	"ipv6":      "isIPv6",
	"email":     "isEmail",
	"intStr":    "isIntString",
	"intString": "isIntString",
}

// ValidatorName get real validator name.
func ValidatorName(name string) string {
	if rName, ok := validatorAliases[name]; ok {
		return rName
	}

	return name
}

/*************************************************************
 * global validators
 *************************************************************/

// global validators. contains built-in and user custom
var (
	validators map[string]interface{}
	// validator func reflect.Value
	validatorValues = map[string]reflect.Value{
		// int value
		"lt":  reflect.ValueOf(Lt),
		"gt":  reflect.ValueOf(Gt),
		"min": reflect.ValueOf(Min),
		"max": reflect.ValueOf(Max),
		// value check
		"enum":     reflect.ValueOf(Enum),
		"notIn":    reflect.ValueOf(NotIn),
		"between":  reflect.ValueOf(Between),
		"regexp":   reflect.ValueOf(Regexp),
		"isEqual":  reflect.ValueOf(IsEqual),
		"intEqual": reflect.ValueOf(IntEqual),
		"notEqual": reflect.ValueOf(NotEqual),
		// data type check
		"isInt":     reflect.ValueOf(IsInt),
		"isMap":     reflect.ValueOf(IsMap),
		"isUint":    reflect.ValueOf(IsUint),
		"isBool":    reflect.ValueOf(IsBool),
		"isFloat":   reflect.ValueOf(IsFloat),
		"isInts":    reflect.ValueOf(IsInts),
		"isArray":   reflect.ValueOf(IsArray),
		"isSlice":   reflect.ValueOf(IsSlice),
		"isString":  reflect.ValueOf(IsString),
		"isStrings": reflect.ValueOf(IsStrings),
		// string
		"isIntString": reflect.ValueOf(IsIntString),
		// length
		"minLength":   reflect.ValueOf(MinLength),
		"maxLength":   reflect.ValueOf(MaxLength),
		"lengthEqual": reflect.ValueOf(Length),
		// common
		"isIP":    reflect.ValueOf(IsIP),
		"isIPv4":  reflect.ValueOf(IsIPv4),
		"isIPv6":  reflect.ValueOf(IsIPv6),
		"isEmail": reflect.ValueOf(IsEmail),
	}
)

// AddValidators to the global validators map
func AddValidators(m map[string]interface{}) {
	for name, checkFunc := range m {
		AddValidator(name, checkFunc)
	}
}

// AddValidator to the pkg. checkFunc must return a bool
func AddValidator(name string, checkFunc interface{}) {
	if validators == nil {
		validators = make(map[string]interface{})
	}

	validators[name] = checkFunc
	validatorValues[name] = checkValidatorFunc(name, checkFunc)
}

// get validator func's reflect.Value
func validatorValue(name string) (reflect.Value, bool) {
	if v, ok := validatorValues[name]; ok {
		return v, true
	}

	return reflect.Value{}, false
}

func checkValidatorFunc(name string, fn interface{}) reflect.Value {
	fv := reflect.ValueOf(fn)

	// is nil or not is func
	if fn == nil || fv.Kind() != reflect.Func {
		panicf("validator '%s'. 'checkFunc' parameter is invalid, it must be an func", name)
	}

	ft := fv.Type()
	if ft.NumOut() != 1 || ft.Out(0).Kind() != reflect.Bool {
		panicf("validator '%s' func must be return a bool value.", name)
	}

	return fv
}

/*************************************************************
 * validators for current validation
 *************************************************************/

// AddValidators to the Validation
func (v *Validation) AddValidators(m map[string]interface{}) {
	for name, checkFunc := range m {
		v.AddValidator(name, checkFunc)
	}
}

// AddValidator to the Validation. checkFunc must return a bool
func (v *Validation) AddValidator(name string, checkFunc interface{}) {
	if v.validatorFuncs == nil {
		v.validatorFuncs = make(map[string]interface{})
	}

	v.validatorFuncs[name] = checkFunc
	v.validatorValues[name] = checkValidatorFunc(name, checkFunc)
}

// ValidatorValue get by name
func (v *Validation) ValidatorValue(name string) (fv reflect.Value, ok bool) {
	name = ValidatorName(name)

	// if v.data is StructData instance.
	if sd, ok := v.data.(*StructData); ok {
		fv, ok = sd.FuncValue(name)
		if ok {
			return fv, true
		}
	}

	// current validation
	if fv, ok = v.validatorValues[name]; ok {
		return
	}

	// global validators
	if fv, ok = validatorValues[name]; ok {
		return
	}

	return
}

// ValidatorFunc get by name
func (v *Validation) ValidatorFunc(name string) interface{} {
	name = ValidatorName(name)
	if fn, ok := v.validatorFuncs[name]; ok {
		return fn
	}

	if fn, ok := validators[name]; ok {
		return fn
	}

	panicf("the validator %s not exists!", name)
	return nil
}

// HasValidator check
func (v *Validation) HasValidator(name string) bool {
	if _, ok := v.validatorFuncs[name]; ok {
		return true
	}

	_, ok := validators[name]
	return ok
}

/*************************************************************
 * context validators
 *************************************************************/

// Required field val check
func (v *Validation) Required(val interface{}) bool {
	return !ValueIsEmpty(reflect.ValueOf(val))
}

// EqField value should EQ the dst field
func (v *Validation) EqField(val interface{}, dstField string) bool {
	// get dst field value.
	dstVal, has := v.Get(dstField)
	if !has {
		return false
	}

	return val == dstVal
}

// NeField value should not equal the dst field
func (v *Validation) NeField(val interface{}, dstField string) bool {
	// get dst field value.
	dstVal, has := v.Get(dstField)
	if !has {
		return false
	}

	return val != dstVal
}

// GtField value should GT the dst field
func (v *Validation) GtField(val interface{}, dstField string) bool {
	// get dst field value.
	dstVal, has := v.Get(dstField)
	if !has {
		return false
	}

	return ValueLen(reflect.ValueOf(val)) > ValueLen(reflect.ValueOf(dstVal))
}

// GteField value should GTE the dst field
func (v *Validation) GteField(val interface{}, dstField string) bool {
	// get dst field value.
	dstVal, has := v.Get(dstField)
	if !has {
		return false
	}

	return ValueLen(reflect.ValueOf(val)) >= ValueLen(reflect.ValueOf(dstVal))
}

// LtField value should LT the dst field
func (v *Validation) LtField(val interface{}, dstField string) bool {
	// get dst field value.
	dstVal, has := v.Get(dstField)
	if !has {
		return false
	}

	return ValueLen(reflect.ValueOf(val)) < ValueLen(reflect.ValueOf(dstVal))
}

// LteField value should LTE the dst field
func (v *Validation) LteField(val interface{}, dstField string) bool {
	// get dst field value.
	dstVal, has := v.Get(dstField)
	if !has {
		return false
	}

	return ValueLen(reflect.ValueOf(val)) <= ValueLen(reflect.ValueOf(dstVal))
}

/*************************************************************
 * global: basic validators
 *************************************************************/

// IsEmpty of the value
func IsEmpty(val interface{}) bool {
	if val == nil {
		return true
	}

	if rv, ok := val.(reflect.Value); ok {
		return ValueIsEmpty(rv)
	}

	return ValueIsEmpty(reflect.ValueOf(val))
}

/*************************************************************
 * global: type validators
 *************************************************************/

// IsUint string
func IsUint(str string) bool {
	_, err := strconv.ParseUint(str, 10, 32)
	return err == nil
}

// IsBool string.
func IsBool(str string) bool {
	_, err := strconv.ParseBool(str)
	return err == nil
}

// IsFloat string
func IsFloat(str string) bool {
	return rxFloat.MatchString(str)
}

// IsArray check
func IsArray(val interface{}) (ok bool) {
	if val == nil {
		return false
	}

	var rv reflect.Value
	if rv, ok = val.(reflect.Value); !ok {
		rv = reflect.ValueOf(val)
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return rv.Kind() == reflect.Array
}

// IsSlice check
func IsSlice(val interface{}) (ok bool) {
	if val == nil {
		return false
	}

	var rv reflect.Value
	if rv, ok = val.(reflect.Value); !ok {
		rv = reflect.ValueOf(val)
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return rv.Kind() == reflect.Slice
}

// IsInts is int slice check
func IsInts(val interface{}) (ok bool) {
	if val == nil {
		return false
	}

	switch val.(type) {
	case []int:
		return true
	case reflect.Value:
		if val.(reflect.Value).Kind() == reflect.Slice {
			_, ok = val.([]int)
		}
	}

	return
}

// IsStrings is string slice check
func IsStrings(val interface{}) (ok bool) {
	if val == nil {
		return false
	}

	switch val.(type) {
	case []string:
		return true
	case reflect.Value:
		if val.(reflect.Value).Kind() == reflect.Slice {
			_, ok = val.([]string)
		}
	}

	return
}

// IsMap check
func IsMap(val interface{}) (ok bool) {
	if val == nil {
		return false
	}

	var rv reflect.Value
	if rv, ok = val.(reflect.Value); !ok {
		rv = reflect.ValueOf(val)
	}

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return rv.Kind() == reflect.Map
}

// IsInt check, and support length check
func IsInt(val interface{}, minAndMax ...int64) (ok bool) {
	if val == nil {
		return false
	}

	var rv reflect.Value
	if rv, ok = val.(reflect.Value); !ok {
		rv = reflect.ValueOf(val)
	}

	intVal, isInt := ValueInt64(rv)

	// @todo convert string to int?
	// if !isInt && rv.Kind() == reflect.String {
	// }

	argLn := len(minAndMax)
	if argLn == 0 { // only check type
		return isInt
	}

	if !isInt {
		return false
	}

	// value check
	minVal := minAndMax[0]
	if argLn == 1 { // only min length check.
		return intVal >= minVal
	}

	maxVal := minAndMax[1]

	// min and max length check
	return intVal >= minVal && intVal <= maxVal
}

// IsString check, and support length check.
// Usage:
// 	ok := IsString(val)
// 	ok := IsString(val, 5) // with min len check
// 	ok := IsString(val, 5, 12) // with min and max len check
func IsString(val interface{}, minAndMaxLen ...int) (ok bool) {
	if val == nil {
		return false
	}

	var rv reflect.Value
	if rv, ok = val.(reflect.Value); !ok {
		rv = reflect.ValueOf(val)
	}

	argLn := len(minAndMaxLen)
	isStr := rv.Type().Kind() == reflect.String

	// only check type
	if argLn == 0 {
		return isStr
	}

	if !isStr {
		return false
	}

	// length check
	strLen := rv.Len()
	minLen := minAndMaxLen[0]

	// only min length check.
	if argLn == 1 {
		return strLen >= minLen
	}

	// min and max length check
	maxLen := minAndMaxLen[1]
	return strLen >= minLen && strLen <= maxLen
}

/*************************************************************
 * global: string validators
 *************************************************************/

// IsIntString check. eg "10"
func IsIntString(str string) bool {
	return filter.String(str).CanInt()
}

// IsASCII string.
func IsASCII(str string) bool {
	return rxASCII.MatchString(str)
}

// IsPrintableASCII string.
func IsPrintableASCII(str string) bool {
	return rxPrintableASCII.MatchString(str)
}

// IsBase64 string.
func IsBase64(str string) bool {
	return rxBase64.MatchString(str)
}

// IsDataURI string.
func IsDataURI(str string) bool {
	return rxDataURI.MatchString(str)
}

// IsMultiByte string.
func IsMultiByte(str string) bool {
	return rxMultiByte.MatchString(str)
}

// IsISBN10 string.
func IsISBN10(str string) bool {
	return rxISBN10.MatchString(str)
}

// IsISBN13 string.
func IsISBN13(str string) bool {
	return rxISBN13.MatchString(str)
}

// IsHexadecimal string.
func IsHexadecimal(str string) bool {
	return rxHexadecimal.MatchString(str)
}

// IsHexColor string.
func IsHexColor(str string) bool {
	return rxHexColor.MatchString(str)
}

// IsRGBColor string.
func IsRGBColor(str string) bool {
	return rxRGBColor.MatchString(str)
}

// IsAlpha string.
func IsAlpha(str string) bool {
	return rxAlpha.MatchString(str)
}

// IsAlphaNum string.
func IsAlphaNum(str string) bool {
	return rxAlphaNum.MatchString(str)
}

// IsNumber string.
func IsNumber(str string) bool {
	return rxNumber.MatchString(str)
}

// IsFilePath string
func IsFilePath(str string) bool {
	return false
}

// IsEmail check
func IsEmail(str string) bool {
	return rxEmail.MatchString(str)
}

func isIPv6(str string) bool {
	return rxIP.MatchString(str)
}

func IsSSN(str string) bool {
	return rxSSN.MatchString(str)
}

// IsIP is the validation function for validating if the field's value is a valid v4 or v6 IP address.
func IsIP(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil
}

// IsIPv4 is the validation function for validating if a value is a valid v4 IP address.
func IsIPv4(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() != nil
}

// IsIPv6 is the validation function for validating if the field's value is a valid v6 IP address.
func IsIPv6(str string) bool {
	ip := net.ParseIP(str)
	return ip != nil && ip.To4() == nil
}

// IsMAC is the validation function for validating if the field's value is a valid MAC address.
func IsMAC(str string) bool {
	_, err := net.ParseMAC(str)
	return err == nil
}

// IsCIDRv4 is the validation function for validating if the field's value is a valid v4 CIDR address.
func IsCIDRv4(str string) bool {
	ip, _, err := net.ParseCIDR(str)
	return err == nil && ip.To4() != nil
}

// IsCIDRv6 is the validation function for validating if the field's value is a valid v6 CIDR address.
func IsCIDRv6(str string) bool {
	ip, _, err := net.ParseCIDR(str)
	return err == nil && ip.To4() == nil
}

// IsCIDR is the validation function for validating if the field's value is a valid v4 or v6 CIDR address.
func IsCIDR(str string) bool {
	_, _, err := net.ParseCIDR(str)
	return err == nil
}

// IsJSON check if the string is valid JSON (note: uses json.Unmarshal).
func IsJSON(str string) bool {
	var js json.RawMessage
	return Unmarshal([]byte(str), &js) == nil
}

/*************************************************************
 * global: compare validators
 *************************************************************/

// IsEqual check
func IsEqual(val, wantVal interface{}) bool {
	equal, err := eq(reflect.ValueOf(val), reflect.ValueOf(wantVal))
	if err != nil {
		return false
	}

	return equal
}

// NotEqual check
func NotEqual(val, wantVal interface{}) bool {
	return !IsEqual(val, wantVal)
}

// IntEqual check
func IntEqual(val interface{}, wantVal int64) bool {
	intVal, isInt := IntVal(val)
	if !isInt {
		return false
	}

	return intVal == wantVal
}

// Gt check value greater dst value
func Gt(val interface{}, dstVal int64) bool {
	intVal, isInt := IntVal(val)
	if !isInt {
		return false
	}

	return intVal > dstVal
}

// Min check value greater or equal dst value. for int(8-64), uint(8-64). alias `Gte`
func Min(val interface{}, min int64) bool {
	intVal, isInt := IntVal(val)
	if !isInt {
		return false
	}

	return intVal >= min
}

// Lt less than dst value
func Lt(val interface{}, dstVal int64) bool {
	intVal, isInt := IntVal(val)
	if !isInt {
		return false
	}

	return intVal < dstVal
}

// Max less than or equal dst value. for int(8-64), uint(8-64). alias `Lte`
func Max(val interface{}, max int64) bool {
	intVal, isInt := IntVal(val)
	if !isInt {
		return false
	}

	return intVal <= max
}

// Between int value in the given range.
func Between(val interface{}, min, max int64) bool {
	rv := reflect.ValueOf(val)

	intVal, isInt := ValueInt64(rv)
	if !isInt {
		return false
	}

	return intVal >= min && intVal <= max
}

// Regexp match value string
func Regexp(str string, pattern string) bool {
	ok, _ := regexp.MatchString(pattern, str)
	return ok
}

/*************************************************************
 * global: array, slice, map validators
 *************************************************************/

// Enum value should be in the given enum(strings, ints, uints).
func Enum(val interface{}, enum interface{}) bool {
	if val == nil {
		return false
	}

	rv := reflect.ValueOf(val)

	// if is string value
	if rv.Kind() == reflect.String {
		strVal := val.(string)

		switch ss := enum.(type) {
		case []string:
			for _, strItem := range ss {
				if strVal == strItem { // exists
					return true
				}
			}
		}

		return false
	}

	// if is int value
	intVal, isInt := ValueInt64(rv)
	if !isInt {
		return false
	}

	if int64s, ok := toInt64Slice(enum); ok {
		for _, i64 := range int64s {
			if intVal == i64 { // exists
				return true
			}
		}
	}

	return false
}

// NotIn value should be not in the given enum(strings, ints, uints).
func NotIn(val interface{}, enum interface{}) bool {
	return false == Enum(val, enum)
}

/*************************************************************
 * global: length validators
 *************************************************************/

// Length equal check for string, array, slice, map
func Length(val interface{}, wantLen int) bool {
	ln := CalcLength(val)
	if ln == -1 {
		return false
	}

	return ln == wantLen
}

// MinLength check for string, array, slice, map
func MinLength(val interface{}, minLen int) bool {
	ln := CalcLength(val)
	if ln == -1 {
		return false
	}

	return ln >= minLen
}

// MaxLength check for string, array, slice, map
func MaxLength(val interface{}, maxLen int) bool {
	ln := CalcLength(val)
	if ln == -1 {
		return false
	}

	return ln <= maxLen
}

// ByteLength check string's length
func ByteLength(str string, params ...string) bool {
	if len(params) == 2 {
		min := filter.MustInt(params[0])
		max := filter.MustInt(params[1])
		strLen := len(str)

		return strLen >= min && strLen <= max
	}

	return false
}

// RuneLength check string's length (including multi byte strings)
func RuneLength(str string, minLen int, maxLen ...int) bool {
	// strLen := len([]rune(str))
	strLen := utf8.RuneCountInString(str)

	// only min length check.
	if len(maxLen) == 0 {
		return strLen >= minLen
	}

	// min and max length check
	return strLen >= minLen && strLen <= maxLen[1]
}

// StringLength check string's length (including multi byte strings)
func StringLength(str string, minLen int, maxLen ...int) bool {
	return RuneLength(str, minLen, maxLen...)
}

/*************************************************************
 * global: date/time validators
 *************************************************************/
