// Code generated by MockGen. DO NOT EDIT.
// Source: internal/ports/email_sender.go
//
// Generated by this command:
//
//	mockgen -source=internal/ports/email_sender.go -destination=internal/adapters/mocks/mock_email_sender.go -package=mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	ports "transaction-processor/internal/ports"

	gomock "go.uber.org/mock/gomock"
	mail "gopkg.in/mail.v2"
)

// MockEmailSender is a mock of EmailSender interface.
type MockEmailSender struct {
	ctrl     *gomock.Controller
	recorder *MockEmailSenderMockRecorder
	isgomock struct{}
}

// MockEmailSenderMockRecorder is the mock recorder for MockEmailSender.
type MockEmailSenderMockRecorder struct {
	mock *MockEmailSender
}

// NewMockEmailSender creates a new mock instance.
func NewMockEmailSender(ctrl *gomock.Controller) *MockEmailSender {
	mock := &MockEmailSender{ctrl: ctrl}
	mock.recorder = &MockEmailSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEmailSender) EXPECT() *MockEmailSenderMockRecorder {
	return m.recorder
}

// SendSummaryEmail mocks base method.
func (m *MockEmailSender) SendSummaryEmail(recipient string, summary ports.EmailSummary) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendSummaryEmail", recipient, summary)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendSummaryEmail indicates an expected call of SendSummaryEmail.
func (mr *MockEmailSenderMockRecorder) SendSummaryEmail(recipient, summary any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendSummaryEmail", reflect.TypeOf((*MockEmailSender)(nil).SendSummaryEmail), recipient, summary)
}

// MockMailDialer is a mock of MailDialer interface.
type MockMailDialer struct {
	ctrl     *gomock.Controller
	recorder *MockMailDialerMockRecorder
	isgomock struct{}
}

// MockMailDialerMockRecorder is the mock recorder for MockMailDialer.
type MockMailDialerMockRecorder struct {
	mock *MockMailDialer
}

// NewMockMailDialer creates a new mock instance.
func NewMockMailDialer(ctrl *gomock.Controller) *MockMailDialer {
	mock := &MockMailDialer{ctrl: ctrl}
	mock.recorder = &MockMailDialerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMailDialer) EXPECT() *MockMailDialerMockRecorder {
	return m.recorder
}

// DialAndSend mocks base method.
func (m_2 *MockMailDialer) DialAndSend(m ...*mail.Message) error {
	m_2.ctrl.T.Helper()
	varargs := []any{}
	for _, a := range m {
		varargs = append(varargs, a)
	}
	ret := m_2.ctrl.Call(m_2, "DialAndSend", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// DialAndSend indicates an expected call of DialAndSend.
func (mr *MockMailDialerMockRecorder) DialAndSend(m ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DialAndSend", reflect.TypeOf((*MockMailDialer)(nil).DialAndSend), m...)
}

// MockMailMessage is a mock of MailMessage interface.
type MockMailMessage struct {
	ctrl     *gomock.Controller
	recorder *MockMailMessageMockRecorder
	isgomock struct{}
}

// MockMailMessageMockRecorder is the mock recorder for MockMailMessage.
type MockMailMessageMockRecorder struct {
	mock *MockMailMessage
}

// NewMockMailMessage creates a new mock instance.
func NewMockMailMessage(ctrl *gomock.Controller) *MockMailMessage {
	mock := &MockMailMessage{ctrl: ctrl}
	mock.recorder = &MockMailMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMailMessage) EXPECT() *MockMailMessageMockRecorder {
	return m.recorder
}

// SetBody mocks base method.
func (m *MockMailMessage) SetBody(contentType, body string, settings ...mail.PartSetting) {
	m.ctrl.T.Helper()
	varargs := []any{contentType, body}
	for _, a := range settings {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "SetBody", varargs...)
}

// SetBody indicates an expected call of SetBody.
func (mr *MockMailMessageMockRecorder) SetBody(contentType, body any, settings ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{contentType, body}, settings...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetBody", reflect.TypeOf((*MockMailMessage)(nil).SetBody), varargs...)
}

// SetHeader mocks base method.
func (m *MockMailMessage) SetHeader(field string, value ...string) {
	m.ctrl.T.Helper()
	varargs := []any{field}
	for _, a := range value {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "SetHeader", varargs...)
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockMailMessageMockRecorder) SetHeader(field any, value ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{field}, value...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockMailMessage)(nil).SetHeader), varargs...)
}

// MockMailMessageFactory is a mock of MailMessageFactory interface.
type MockMailMessageFactory struct {
	ctrl     *gomock.Controller
	recorder *MockMailMessageFactoryMockRecorder
	isgomock struct{}
}

// MockMailMessageFactoryMockRecorder is the mock recorder for MockMailMessageFactory.
type MockMailMessageFactoryMockRecorder struct {
	mock *MockMailMessageFactory
}

// NewMockMailMessageFactory creates a new mock instance.
func NewMockMailMessageFactory(ctrl *gomock.Controller) *MockMailMessageFactory {
	mock := &MockMailMessageFactory{ctrl: ctrl}
	mock.recorder = &MockMailMessageFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMailMessageFactory) EXPECT() *MockMailMessageFactoryMockRecorder {
	return m.recorder
}

// NewMessage mocks base method.
func (m *MockMailMessageFactory) NewMessage() ports.MailMessage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewMessage")
	ret0, _ := ret[0].(ports.MailMessage)
	return ret0
}

// NewMessage indicates an expected call of NewMessage.
func (mr *MockMailMessageFactoryMockRecorder) NewMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewMessage", reflect.TypeOf((*MockMailMessageFactory)(nil).NewMessage))
}
