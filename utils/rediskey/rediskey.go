package rediskey

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"strings"

	"github.com/mitchellh/hashstructure/v2"
	"rederinghub.io/utils"
)

const (
	concatStr = '_'
)

type key struct {
	sb strings.Builder
}

// Beauty beauty redis key to format <prefix_uuid_md5ValueOfParams>
// Example: Beauty("aa").WithUUID(123).WithParams("code")
func Beauty(prefix string) *key {
	k := &key{}
	k.sb.WriteString(prefix)
	return k
}

func (k *key) String() string {
	return k.sb.String()
}

// WithUUID concat uuid to key
func (k *key) WithUUID(uuid any) *key {
	k.sb.WriteRune(concatStr)

	uuidStr := fmt.Sprint(uuid)
	k.sb.WriteString(uuidStr)
	return k
}

// WithParams concat params to key
func (k *key) WithParams(params ...interface{}) *key {
	for _, item := range params {
		k.sb.WriteRune(concatStr)
		k.sb.WriteString(fmt.Sprintf("%v", item))
	}

	return k
}

func (k *key) WithStructHash(str interface{}, opts *hashstructure.HashOptions) *key {
	k.sb.WriteRune(concatStr)
	k.sb.WriteString(utils.HashStruct(str, opts))
	return k
}

// WithParamsHash md5 params and concat to key
func (k *key) WithParamsHash(params ...interface{}) *key {
	k.sb.WriteRune(concatStr)
	k.sb.WriteString(md5Params(params...))
	return k
}

func md5Params(params ...interface{}) string {
	if len(params) == 0 {
		return ""
	}

	var buffer bytes.Buffer
	for _, v := range params {
		buffer.WriteString(fmt.Sprintf("%v", v))
	}
	return fmt.Sprintf("%x", md5.Sum(buffer.Bytes()))
}
