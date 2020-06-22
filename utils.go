package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/hbagdi/go-kong/kong"
	"reflect"
	"strings"
)

// 字符串包含
func contains(obj interface{}, target interface{}) bool {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true
		}
	}
	return false
}

// 提取域名列表B中后缀为域名列表A的值
func domainsContains(a []string, b []string) []string {
	set := make([]string, 0)
	for i := 0; i < len(a); i++ {
		al := strings.Split(a[i], ".")
		av := strings.Join(al[1:], ".")
		if len(al) < 3 {
			av = a[i]
		}
		
		for j := 0; j < len(b); j++ {
			bl := strings.Split(b[j], ".")
			bv := strings.Join(bl[1:], ".")
			if len(bl) < 3 {
				bv = b[j]
			}
			if bv == av {
				set = append(set, b[j])
			}
		}
	}
	return sliceDedup(set)
}

// 列表排重
func sliceDedup(source []string) []string {
	result := make([]string, 0, len(source))
	temp := map[string]struct{}{}
	for _, item := range source {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

// StringSlice 将一列字符串值转换为一列字符串指针
func stringSlice(src []string) []*string {
	dst := make([]*string, len(src))
	for i := 0; i < len(src); i++ {
		dst[i] = &(src[i])
	}
	return dst
}

// StringValueSlice 将一列字符串指针转换为一列字符串值
func stringValueSlice(src []*string) []string {
	dst := make([]string, len(src))
	for i := 0; i < len(src); i++ {
		if src[i] != nil {
			dst[i] = *(src[i])
		}
	}
	return dst
}

// StringMap converts a string map of string values into a string
// map of string pointers
func stringMap(src map[string]string) map[string]*string {
	dst := make(map[string]*string)
	for k, val := range src {
		v := val
		dst[k] = &v
	}
	return dst
}

// StringValueMap converts a string map of string pointers into a string
// map of string values
func stringValueMap(src map[string]*string) map[string]string {
	dst := make(map[string]string)
	for k, val := range src {
		if val != nil {
			dst[k] = *val
		}
	}
	return dst
}

func isEmptyString(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}

type SNI struct {
	Name             *string `json:"name,omitempty" yaml:"name,omitempty"`
	CreatedAt        *int64  `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	SslCertificateId *string `json:"ssl_certificate_id,omitempty" yaml:"ssl_certificate_id,omitempty"`
}

// Create creates a SNI in Kong0.13.
// If an ID is specified, it will be used to
// create a sni in Kong, otherwise an ID
// is auto-generated.
func SinCreate(client *kong.Client, ctx context.Context, sni *SNI) (*SNI, error) {

	queryPath := "/snis"
	method := "POST"
	req, err := client.NewRequest(method, queryPath, nil, sni)

	if err != nil {
		return nil, err
	}

	var createdSNI SNI
	_, err = client.Do(ctx, req, &createdSNI)
	if err != nil {
		return nil, err
	}
	return &createdSNI, nil
}

// Get fetches a SNI in Kong.
func SinGet(client *kong.Client, ctx context.Context,
	sinName *string) (bool, error) {

	if isEmptyString(sinName) {
		return false, errors.New(
			"sinName cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/snis/%v", *sinName)
	req, err := client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return false, err
	}

	var sni SNI
	_, err = client.Do(ctx, req, &sni)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Update updates a SNI in Kong
func SinUpdate(client *kong.Client, ctx context.Context, sni *SNI) (*SNI, error) {

	if isEmptyString(sni.Name) {
		return nil, errors.New("Name cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/snis/%v", *sni.Name)
	req, err := client.NewRequest("PATCH", endpoint, nil, sni)
	if err != nil {
		return nil, err
	}

	var updatedAPI SNI
	_, err = client.Do(ctx, req, &updatedAPI)
	if err != nil {
		return nil, err
	}
	return &updatedAPI, nil
}

// Delete deletes a SNI in Kong
func SinDelete(client *kong.Client, ctx context.Context, sinName *string) error {

	if isEmptyString(sinName) {
		return errors.New("usernameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/snis/%v", *sinName)
	req, err := client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = client.Do(ctx, req, nil)
	return err
}
