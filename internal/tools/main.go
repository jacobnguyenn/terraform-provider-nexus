package tools

import (
	"os"
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func InterfaceSliceToStringSlice(data []interface{}) []string {
	result := make([]string, len(data))
	for i, v := range data {
		result[i] = v.(string)
	}
	return result
}

func StringSliceToInterfaceSlice(strings []string) []interface{} {
	s := make([]interface{}, len(strings))
	for i, v := range strings {
		s[i] = string(v)
	}
	return s
}

func GetEnv(key, fallback string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return fallback
	}
	return value
}

// Copied from https://siongui.github.io/2018/03/09/go-match-common-element-in-two-array/
func Intersection(a, b []int) (c []int) {
	m := make(map[int]bool)

	for _, item := range a {
		m[item] = true
	}

	for _, item := range b {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return
}

func GetIntPointer(number int) *int {
	return &number
}

func GetStringPointer(s string) *string {
	return &s
}

func GetBoolPointer(b bool) *bool {
	return &b
}

func ConvertStringSet(set *schema.Set) []string {
	s := make([]string, 0, set.Len())
	for _, v := range set.List() {
		s = append(s, v.(string))
	}
	sort.Strings(s)

	return s
}

const (
	SortKey                     = "order"
	NegativeCacheDefaultEnabled = false
	NegativeCacheDefaultTTL     = 1440
)

func SortSliceByKey(s []interface{}, key string) {
	sort.Slice(s, func(i, j int) bool {
		return s[i].(map[string]interface{})[key].(int) < s[j].(map[string]interface{})[key].(int)
	})
}
