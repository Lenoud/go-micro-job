package types

import (
	"reflect"
	"testing"
)

func TestUserTypesDoNotExposeLegacyPushEmailFields(t *testing.T) {
	for _, tc := range []struct {
		name  string
		value interface{}
	}{
		{name: "UserInfo", value: UserInfo{}},
		{name: "UpdateUserReq", value: UpdateUserReq{}},
		{name: "UpdateUserInfoReq", value: UpdateUserInfoReq{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			typ := reflect.TypeOf(tc.value)
			if _, ok := typ.FieldByName("PushEmail"); ok {
				t.Fatalf("%s still exposes PushEmail", tc.name)
			}
			if _, ok := typ.FieldByName("PushSwitch"); ok {
				t.Fatalf("%s still exposes PushSwitch", tc.name)
			}
		})
	}
}
