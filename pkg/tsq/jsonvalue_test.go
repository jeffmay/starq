package tsq_test

import (
	"testing"

	"github.com/jeffmay/starq/pkg/tsq"
	"github.com/stretchr/testify/require"
)

func TestInvalidJSON(t *testing.T) {
	obj := tsq.JSONObject(map[string]*tsq.JSONValue{})
	obj.MustGetString()
}

func TestJSONNull(t *testing.T) {
	null1 := tsq.JSONNull()
	require.True(t, null1.IsNull())
	null2 := tsq.JSONNull()
	require.Equal(t, null1, null2)
	require.True(t, null1.IsEqual(null2))
}

func TestJSONString(t *testing.T) {
	str := "hello"
	str1 := tsq.JSONString(str)
	require.Equal(t, str, str1.MustGetString())
	str2 := tsq.JSONString(str)
	require.True(t, str1.IsEqual(str2))
	// this is required to coerce the type from rawString to String to match str1
	require.Equal(t, str, str2.MustGetString())
	require.Equal(t, str1, str2)
}

func TestJSONNumber(t *testing.T) {
	num := 42
	num1 := tsq.JSONNumber(num)
	require.Equal(t, int64(42), num1.MustGetInt64())
	require.Equal(t, float64(42.0), num1.MustGetFloat64())
	num2 := tsq.JSONNumber(num)
	require.Equal(t, num1, num2)
	require.True(t, num1.IsEqual(num2))
}

func TestJSONBool(t *testing.T) {
	bool1 := tsq.JSONBool(true)
	require.True(t, bool1.MustGetBool())
	bool2 := tsq.JSONBool(true)
	require.Equal(t, bool1, bool2)
	require.True(t, bool1.IsEqual(bool2))
}

func TestJSONArray(t *testing.T) {
	a := []*tsq.JSONValue{
		tsq.JSONInt64(1),
		tsq.JSONFloat64(2),
		tsq.JSONString("3"),
	}
	arr1 := tsq.JSONArray(a...)
	require.Len(t, arr1.MustGetArray(), 3)
	arr2 := tsq.JSONArray(a...)
	require.Equal(t, arr1, arr2)
	require.True(t, arr1.IsEqual(arr2))
}

func TestJSONObject(t *testing.T) {
	m := map[string]*tsq.JSONValue{
		"a": tsq.JSONNumber(1),
		"b": tsq.JSONNumber(2),
		"c": tsq.JSONNumber(3),
	}
	obj1 := tsq.JSONObject(m)
	require.Len(t, obj1.MustGetObject(), 3)
	obj2 := tsq.JSONObject(m)
	require.Equal(t, obj1, obj2)
	require.True(t, obj1.IsEqual(obj2))
}

func TestParseObject(t *testing.T) {
	obj := tsq.ParseJSON(`{
		"str": "abc",
		"int": 42,
		"float": 1.23,
		"bool": false,
		"arr": ["1", "2", "3"],
		"null": null,
		"1": {
			"str": "def",
			"int": 84,
			"float": 2.46,
			"bool": true,
			"arr": ["4", 5, "6"],
			"2": {
				"null": null,
				"nestedarr": [
					[{ "value": "needle" }]
				]
			}
		}
	}`)
	require.Equal(t, "abc", obj.MustGetString("str"))
	require.EqualValues(t, 42, obj.MustGetInt64("int"))
	require.EqualValues(t, 1.23, obj.MustGetFloat64("float"))
	require.False(t, obj.MustGetBool("bool"))
	expectArr := tsq.JSONArray(
		tsq.JSONString("1"),
		tsq.JSONString("2"),
		tsq.JSONString("3"),
	).MustGetArray()
	require.Equal(t, expectArr, obj.MustGetArray("arr"))
	require.Equal(t, tsq.JSONNull(), obj.MustGet("null"))
	nested1 := obj.MustGet("1")
	require.Equal(t, "def", nested1.MustGetBool("str"))
	require.EqualValues(t, 84, nested1.MustGetInt64("int"))
}

func TestJSONInt64Pretty(t *testing.T) {
	num := tsq.JSONInt64(42)
	require.Equal(t, "42", num.Pretty())
}

func TestJSONFloat64Pretty(t *testing.T) {
	num := tsq.JSONFloat64(1.23)
	require.Equal(t, "1.23", num.Pretty())
}

func TestJSONStringPretty(t *testing.T) {
	num := tsq.JSONString("🤪")
	require.Equal(t, `"🤪"`, num.Pretty())
}

func TestJSONArrayPretty(t *testing.T) {
	arr := tsq.JSONArray(
		tsq.JSONArray(
			tsq.JSONInt64(1),
			tsq.JSONInt64(2),
		),
	)
	expected := "[\n  [\n    1,\n    2\n  ]\n]"
	require.Equal(t, expected, arr.Pretty())
}

func TestJSONObjectPretty(t *testing.T) {
	obj := tsq.JSONObject(map[string]*tsq.JSONValue{
		"nested": tsq.JSONObject(map[string]*tsq.JSONValue{
			"a": tsq.JSONInt64(1),
			"b": tsq.JSONInt64(2),
		}),
	})
	expected := "{\n  \"nested\": {\n    \"a\": 1,\n    \"b\": 2\n  }\n}"
	require.Equal(t, expected, obj.Pretty())
}
