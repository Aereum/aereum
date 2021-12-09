package util

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

type JSONBuilder struct {
	Encode strings.Builder
}

func (j *JSONBuilder) putGeneral(fieldName, value string) {
	if j.Encode.Len() > 0 {
		fmt.Fprintf(&j.Encode, `,"%v":%v`, fieldName, value)
	} else {
		fmt.Fprintf(&j.Encode, `"%v":%v`, fieldName, value)
	}
}

func (j *JSONBuilder) PutTime(fieldName string, t time.Time) {
	j.putGeneral(fieldName, t.Format(time.RFC3339))
}

func (j *JSONBuilder) PutUint64(fieldName string, value uint64) {
	j.putGeneral(fieldName, fmt.Sprintf("%v", value))
}

func (j *JSONBuilder) PutHex(fieldName string, value []byte) {
	if len(value) == 0 {
		return
	}
	j.putGeneral(fieldName, fmt.Sprintf(`"0x%v"`, hex.EncodeToString(value)))
}

func (j *JSONBuilder) PutBase64(fieldName string, value []byte) {
	if len(value) == 0 {
		return
	}
	j.putGeneral(fieldName, fmt.Sprintf(`"%v"`, base64.StdEncoding.EncodeToString(value)))
}

func (j *JSONBuilder) PutString(fieldName, value string) {
	j.putGeneral(fieldName, fmt.Sprintf(`"%v"`, value))
}

func (j *JSONBuilder) PutJSON(fieldName, value string) {
	j.putGeneral(fieldName, value)
}

func (j *JSONBuilder) ToString() string {
	return fmt.Sprintf(`{%v}`, j.Encode.String())
}
