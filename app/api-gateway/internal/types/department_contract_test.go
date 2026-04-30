package types

import (
	"reflect"
	"testing"
)

func TestDepartmentTypesExposeStableContract(t *testing.T) {
	typ := reflect.TypeOf(DepartmentInfo{})
	for _, field := range []string{"Id", "Title", "Description", "ParentId", "CreateTime"} {
		if _, ok := typ.FieldByName(field); !ok {
			t.Fatalf("DepartmentInfo missing %s", field)
		}
	}
}
