package httpmod

import (
	"github.com/ua-parser/uap-go/uaparser"
	"gitlab.com/snap-clickstaff/go-common/comtypes"
)

var userAgentParser = comtypes.NewSingleton(func() interface{} {
	return uaparser.NewFromSaved()
})

func UserAgentGetSystemParser() *uaparser.Parser {
	return userAgentParser.Get().(*uaparser.Parser)
}

func UserAgentParse(uaString string) *uaparser.Client {
	return UserAgentGetSystemParser().Parse(uaString)
}
