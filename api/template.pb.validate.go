// Code generated by protoc-gen-validate. DO NOT EDIT.
// source: template.proto

package api

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"net/url"
	"regexp"
	"sort"
	"strings"
	"time"
	"unicode/utf8"

	"google.golang.org/protobuf/types/known/anypb"
)

// ensure the imports are used
var (
	_ = bytes.MinRead
	_ = errors.New("")
	_ = fmt.Print
	_ = utf8.UTFMax
	_ = (*regexp.Regexp)(nil)
	_ = (*strings.Reader)(nil)
	_ = net.IPv4len
	_ = time.Duration(0)
	_ = (*url.URL)(nil)
	_ = (*mail.Address)(nil)
	_ = anypb.Any{}
	_ = sort.Sort
)

// Validate checks the field values on GetTemplateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetTemplateRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetTemplateRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetTemplateRequestMultiError, or nil if none found.
func (m *GetTemplateRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetTemplateRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if len(errors) > 0 {
		return GetTemplateRequestMultiError(errors)
	}

	return nil
}

// GetTemplateRequestMultiError is an error wrapping multiple validation errors
// returned by GetTemplateRequest.ValidateAll() if the designated constraints
// aren't met.
type GetTemplateRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetTemplateRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetTemplateRequestMultiError) AllErrors() []error { return m }

// GetTemplateRequestValidationError is the validation error returned by
// GetTemplateRequest.Validate if the designated constraints aren't met.
type GetTemplateRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetTemplateRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetTemplateRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetTemplateRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetTemplateRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetTemplateRequestValidationError) ErrorName() string {
	return "GetTemplateRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetTemplateRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetTemplateRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetTemplateRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetTemplateRequestValidationError{}

// Validate checks the field values on Template with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Template) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Template with the rules defined in
// the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in TemplateMultiError, or nil
// if none found.
func (m *Template) ValidateAll() error {
	return m.validate(true)
}

func (m *Template) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Name

	// no validation rules for TokenId

	// no validation rules for Symbol

	// no validation rules for MetadataImage

	if len(errors) > 0 {
		return TemplateMultiError(errors)
	}

	return nil
}

// TemplateMultiError is an error wrapping multiple validation errors returned
// by Template.ValidateAll() if the designated constraints aren't met.
type TemplateMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TemplateMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TemplateMultiError) AllErrors() []error { return m }

// TemplateValidationError is the validation error returned by
// Template.Validate if the designated constraints aren't met.
type TemplateValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TemplateValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TemplateValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TemplateValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TemplateValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TemplateValidationError) ErrorName() string { return "TemplateValidationError" }

// Error satisfies the builtin error interface
func (e TemplateValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTemplate.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TemplateValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TemplateValidationError{}

// Validate checks the field values on TemplateParam with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *TemplateParam) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TemplateParam with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in TemplateParamMultiError, or
// nil if none found.
func (m *TemplateParam) ValidateAll() error {
	return m.validate(true)
}

func (m *TemplateParam) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for TypeI

	// no validation rules for Max

	// no validation rules for Min

	// no validation rules for Decimal

	if m.Value != nil {
		// no validation rules for Value
	}

	if len(errors) > 0 {
		return TemplateParamMultiError(errors)
	}

	return nil
}

// TemplateParamMultiError is an error wrapping multiple validation errors
// returned by TemplateParam.ValidateAll() if the designated constraints
// aren't met.
type TemplateParamMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TemplateParamMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TemplateParamMultiError) AllErrors() []error { return m }

// TemplateParamValidationError is the validation error returned by
// TemplateParam.Validate if the designated constraints aren't met.
type TemplateParamValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TemplateParamValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TemplateParamValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TemplateParamValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TemplateParamValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TemplateParamValidationError) ErrorName() string { return "TemplateParamValidationError" }

// Error satisfies the builtin error interface
func (e TemplateParamValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTemplateParam.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TemplateParamValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TemplateParamValidationError{}

// Validate checks the field values on NftInfo with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *NftInfo) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on NftInfo with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in NftInfoMultiError, or nil if none found.
func (m *NftInfo) ValidateAll() error {
	return m.validate(true)
}

func (m *NftInfo) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for NetworkType

	// no validation rules for ChainId

	// no validation rules for TokenId

	// no validation rules for ContractAddress

	if len(errors) > 0 {
		return NftInfoMultiError(errors)
	}

	return nil
}

// NftInfoMultiError is an error wrapping multiple validation errors returned
// by NftInfo.ValidateAll() if the designated constraints aren't met.
type NftInfoMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m NftInfoMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m NftInfoMultiError) AllErrors() []error { return m }

// NftInfoValidationError is the validation error returned by NftInfo.Validate
// if the designated constraints aren't met.
type NftInfoValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e NftInfoValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e NftInfoValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e NftInfoValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e NftInfoValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e NftInfoValidationError) ErrorName() string { return "NftInfoValidationError" }

// Error satisfies the builtin error interface
func (e NftInfoValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sNftInfo.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = NftInfoValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = NftInfoValidationError{}

// Validate checks the field values on GetTemplateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetTemplateResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetTemplateResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetTemplateResponseMultiError, or nil if none found.
func (m *GetTemplateResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetTemplateResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	for idx, item := range m.GetTemplate() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, GetTemplateResponseValidationError{
						field:  fmt.Sprintf("Template[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, GetTemplateResponseValidationError{
						field:  fmt.Sprintf("Template[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return GetTemplateResponseValidationError{
					field:  fmt.Sprintf("Template[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	// no validation rules for Total

	if len(errors) > 0 {
		return GetTemplateResponseMultiError(errors)
	}

	return nil
}

// GetTemplateResponseMultiError is an error wrapping multiple validation
// errors returned by GetTemplateResponse.ValidateAll() if the designated
// constraints aren't met.
type GetTemplateResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetTemplateResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetTemplateResponseMultiError) AllErrors() []error { return m }

// GetTemplateResponseValidationError is the validation error returned by
// GetTemplateResponse.Validate if the designated constraints aren't met.
type GetTemplateResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetTemplateResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetTemplateResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetTemplateResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetTemplateResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetTemplateResponseValidationError) ErrorName() string {
	return "GetTemplateResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetTemplateResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetTemplateResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetTemplateResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetTemplateResponseValidationError{}

// Validate checks the field values on GetTemplateDetailRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetTemplateDetailRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetTemplateDetailRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetTemplateDetailRequestMultiError, or nil if none found.
func (m *GetTemplateDetailRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *GetTemplateDetailRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for TokenId

	// no validation rules for ChainId

	// no validation rules for ContractAddress

	if len(errors) > 0 {
		return GetTemplateDetailRequestMultiError(errors)
	}

	return nil
}

// GetTemplateDetailRequestMultiError is an error wrapping multiple validation
// errors returned by GetTemplateDetailRequest.ValidateAll() if the designated
// constraints aren't met.
type GetTemplateDetailRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetTemplateDetailRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetTemplateDetailRequestMultiError) AllErrors() []error { return m }

// GetTemplateDetailRequestValidationError is the validation error returned by
// GetTemplateDetailRequest.Validate if the designated constraints aren't met.
type GetTemplateDetailRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetTemplateDetailRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetTemplateDetailRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetTemplateDetailRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetTemplateDetailRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetTemplateDetailRequestValidationError) ErrorName() string {
	return "GetTemplateDetailRequestValidationError"
}

// Error satisfies the builtin error interface
func (e GetTemplateDetailRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetTemplateDetailRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetTemplateDetailRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetTemplateDetailRequestValidationError{}

// Validate checks the field values on GetTemplateDetailResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *GetTemplateDetailResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on GetTemplateDetailResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// GetTemplateDetailResponseMultiError, or nil if none found.
func (m *GetTemplateDetailResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *GetTemplateDetailResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	if all {
		switch v := interface{}(m.GetNftInfo()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetTemplateDetailResponseValidationError{
					field:  "NftInfo",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetTemplateDetailResponseValidationError{
					field:  "NftInfo",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetNftInfo()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetTemplateDetailResponseValidationError{
				field:  "NftInfo",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for Fee

	// no validation rules for FeeToken

	// no validation rules for MintMaxSupply

	// no validation rules for MintTotalSupply

	// no validation rules for Script

	// no validation rules for ScriptType

	// no validation rules for Creator

	// no validation rules for CustomUri

	// no validation rules for ProjectName

	// no validation rules for ClientSeed

	if all {
		switch v := interface{}(m.GetParamsTemplate()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, GetTemplateDetailResponseValidationError{
					field:  "ParamsTemplate",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, GetTemplateDetailResponseValidationError{
					field:  "ParamsTemplate",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetParamsTemplate()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return GetTemplateDetailResponseValidationError{
				field:  "ParamsTemplate",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for MinterNftInfo

	if len(errors) > 0 {
		return GetTemplateDetailResponseMultiError(errors)
	}

	return nil
}

// GetTemplateDetailResponseMultiError is an error wrapping multiple validation
// errors returned by GetTemplateDetailResponse.ValidateAll() if the
// designated constraints aren't met.
type GetTemplateDetailResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m GetTemplateDetailResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m GetTemplateDetailResponseMultiError) AllErrors() []error { return m }

// GetTemplateDetailResponseValidationError is the validation error returned by
// GetTemplateDetailResponse.Validate if the designated constraints aren't met.
type GetTemplateDetailResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e GetTemplateDetailResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e GetTemplateDetailResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e GetTemplateDetailResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e GetTemplateDetailResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e GetTemplateDetailResponseValidationError) ErrorName() string {
	return "GetTemplateDetailResponseValidationError"
}

// Error satisfies the builtin error interface
func (e GetTemplateDetailResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sGetTemplateDetailResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = GetTemplateDetailResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = GetTemplateDetailResponseValidationError{}

// Validate checks the field values on Param with the rules defined in the
// proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *Param) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on Param with the rules defined in the
// proto definition for this message. If any rules are violated, the result is
// a list of violation errors wrapped in ParamMultiError, or nil if none found.
func (m *Param) ValidateAll() error {
	return m.validate(true)
}

func (m *Param) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for TypeValue

	// no validation rules for Max

	// no validation rules for Min

	// no validation rules for Decimal

	// no validation rules for Value

	// no validation rules for Editable

	if len(errors) > 0 {
		return ParamMultiError(errors)
	}

	return nil
}

// ParamMultiError is an error wrapping multiple validation errors returned by
// Param.ValidateAll() if the designated constraints aren't met.
type ParamMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ParamMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ParamMultiError) AllErrors() []error { return m }

// ParamValidationError is the validation error returned by Param.Validate if
// the designated constraints aren't met.
type ParamValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ParamValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ParamValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ParamValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ParamValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ParamValidationError) ErrorName() string { return "ParamValidationError" }

// Error satisfies the builtin error interface
func (e ParamValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sParam.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ParamValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ParamValidationError{}

// Validate checks the field values on ParamsTemplate with the rules defined in
// the proto definition for this message. If any rules are violated, the first
// error encountered is returned, or nil if there are no violations.
func (m *ParamsTemplate) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on ParamsTemplate with the rules defined
// in the proto definition for this message. If any rules are violated, the
// result is a list of violation errors wrapped in ParamsTemplateMultiError,
// or nil if none found.
func (m *ParamsTemplate) ValidateAll() error {
	return m.validate(true)
}

func (m *ParamsTemplate) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Seed

	for idx, item := range m.GetParams() {
		_, _ = idx, item

		if all {
			switch v := interface{}(item).(type) {
			case interface{ ValidateAll() error }:
				if err := v.ValidateAll(); err != nil {
					errors = append(errors, ParamsTemplateValidationError{
						field:  fmt.Sprintf("Params[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			case interface{ Validate() error }:
				if err := v.Validate(); err != nil {
					errors = append(errors, ParamsTemplateValidationError{
						field:  fmt.Sprintf("Params[%v]", idx),
						reason: "embedded message failed validation",
						cause:  err,
					})
				}
			}
		} else if v, ok := interface{}(item).(interface{ Validate() error }); ok {
			if err := v.Validate(); err != nil {
				return ParamsTemplateValidationError{
					field:  fmt.Sprintf("Params[%v]", idx),
					reason: "embedded message failed validation",
					cause:  err,
				}
			}
		}

	}

	if len(errors) > 0 {
		return ParamsTemplateMultiError(errors)
	}

	return nil
}

// ParamsTemplateMultiError is an error wrapping multiple validation errors
// returned by ParamsTemplate.ValidateAll() if the designated constraints
// aren't met.
type ParamsTemplateMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m ParamsTemplateMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m ParamsTemplateMultiError) AllErrors() []error { return m }

// ParamsTemplateValidationError is the validation error returned by
// ParamsTemplate.Validate if the designated constraints aren't met.
type ParamsTemplateValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e ParamsTemplateValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e ParamsTemplateValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e ParamsTemplateValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e ParamsTemplateValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e ParamsTemplateValidationError) ErrorName() string { return "ParamsTemplateValidationError" }

// Error satisfies the builtin error interface
func (e ParamsTemplateValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sParamsTemplate.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = ParamsTemplateValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = ParamsTemplateValidationError{}

// Validate checks the field values on TemplateRenderingRequest with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TemplateRenderingRequest) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TemplateRenderingRequest with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TemplateRenderingRequestMultiError, or nil if none found.
func (m *TemplateRenderingRequest) ValidateAll() error {
	return m.validate(true)
}

func (m *TemplateRenderingRequest) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for TokenId

	if all {
		switch v := interface{}(m.GetParams()).(type) {
		case interface{ ValidateAll() error }:
			if err := v.ValidateAll(); err != nil {
				errors = append(errors, TemplateRenderingRequestValidationError{
					field:  "Params",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		case interface{ Validate() error }:
			if err := v.Validate(); err != nil {
				errors = append(errors, TemplateRenderingRequestValidationError{
					field:  "Params",
					reason: "embedded message failed validation",
					cause:  err,
				})
			}
		}
	} else if v, ok := interface{}(m.GetParams()).(interface{ Validate() error }); ok {
		if err := v.Validate(); err != nil {
			return TemplateRenderingRequestValidationError{
				field:  "Params",
				reason: "embedded message failed validation",
				cause:  err,
			}
		}
	}

	// no validation rules for ChainId

	// no validation rules for ContractAddress

	if len(errors) > 0 {
		return TemplateRenderingRequestMultiError(errors)
	}

	return nil
}

// TemplateRenderingRequestMultiError is an error wrapping multiple validation
// errors returned by TemplateRenderingRequest.ValidateAll() if the designated
// constraints aren't met.
type TemplateRenderingRequestMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TemplateRenderingRequestMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TemplateRenderingRequestMultiError) AllErrors() []error { return m }

// TemplateRenderingRequestValidationError is the validation error returned by
// TemplateRenderingRequest.Validate if the designated constraints aren't met.
type TemplateRenderingRequestValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TemplateRenderingRequestValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TemplateRenderingRequestValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TemplateRenderingRequestValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TemplateRenderingRequestValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TemplateRenderingRequestValidationError) ErrorName() string {
	return "TemplateRenderingRequestValidationError"
}

// Error satisfies the builtin error interface
func (e TemplateRenderingRequestValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTemplateRenderingRequest.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TemplateRenderingRequestValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TemplateRenderingRequestValidationError{}

// Validate checks the field values on TemplateRenderingResponse with the rules
// defined in the proto definition for this message. If any rules are
// violated, the first error encountered is returned, or nil if there are no violations.
func (m *TemplateRenderingResponse) Validate() error {
	return m.validate(false)
}

// ValidateAll checks the field values on TemplateRenderingResponse with the
// rules defined in the proto definition for this message. If any rules are
// violated, the result is a list of violation errors wrapped in
// TemplateRenderingResponseMultiError, or nil if none found.
func (m *TemplateRenderingResponse) ValidateAll() error {
	return m.validate(true)
}

func (m *TemplateRenderingResponse) validate(all bool) error {
	if m == nil {
		return nil
	}

	var errors []error

	// no validation rules for Glb

	// no validation rules for Image

	if len(errors) > 0 {
		return TemplateRenderingResponseMultiError(errors)
	}

	return nil
}

// TemplateRenderingResponseMultiError is an error wrapping multiple validation
// errors returned by TemplateRenderingResponse.ValidateAll() if the
// designated constraints aren't met.
type TemplateRenderingResponseMultiError []error

// Error returns a concatenation of all the error messages it wraps.
func (m TemplateRenderingResponseMultiError) Error() string {
	var msgs []string
	for _, err := range m {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

// AllErrors returns a list of validation violation errors.
func (m TemplateRenderingResponseMultiError) AllErrors() []error { return m }

// TemplateRenderingResponseValidationError is the validation error returned by
// TemplateRenderingResponse.Validate if the designated constraints aren't met.
type TemplateRenderingResponseValidationError struct {
	field  string
	reason string
	cause  error
	key    bool
}

// Field function returns field value.
func (e TemplateRenderingResponseValidationError) Field() string { return e.field }

// Reason function returns reason value.
func (e TemplateRenderingResponseValidationError) Reason() string { return e.reason }

// Cause function returns cause value.
func (e TemplateRenderingResponseValidationError) Cause() error { return e.cause }

// Key function returns key value.
func (e TemplateRenderingResponseValidationError) Key() bool { return e.key }

// ErrorName returns error name.
func (e TemplateRenderingResponseValidationError) ErrorName() string {
	return "TemplateRenderingResponseValidationError"
}

// Error satisfies the builtin error interface
func (e TemplateRenderingResponseValidationError) Error() string {
	cause := ""
	if e.cause != nil {
		cause = fmt.Sprintf(" | caused by: %v", e.cause)
	}

	key := ""
	if e.key {
		key = "key for "
	}

	return fmt.Sprintf(
		"invalid %sTemplateRenderingResponse.%s: %s%s",
		key,
		e.field,
		e.reason,
		cause)
}

var _ error = TemplateRenderingResponseValidationError{}

var _ interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
} = TemplateRenderingResponseValidationError{}
