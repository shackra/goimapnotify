// Code generated by mockigo. DO NOT EDIT.

//go:generate mockigo 

package main

import client "github.com/emersion/go-imap/client"
import goimap "github.com/emersion/go-imap"
import match "github.com/subtle-byte/mockigo/match"
import mock "github.com/subtle-byte/mockigo/mock"
import time "time"

var _ = match.Any[int]

type IdleClientMock struct {
	mock *mock.Mock
}

func NewIdleClientMock(t mock.Testing) *IdleClientMock {
	t.Helper()
	return &IdleClientMock{mock: mock.NewMock(t)}
}

type _IdleClientMock_Expecter struct {
	mock *mock.Mock
}

func (_mock *IdleClientMock) EXPECT() _IdleClientMock_Expecter {
	 return _IdleClientMock_Expecter{mock: _mock.mock}
}

type _IdleClientMock_IdleWithFallback_Call struct {
	*mock.Call
}

func (_mock *IdleClientMock) IdleWithFallback(_a0 <-chan struct{}, _a1 time.Duration) (error) {
	_mock.mock.T.Helper()
	_results := _mock.mock.Called("IdleWithFallback", _a0, _a1)
	_r0 := _results.Error(0)
	return _r0
}

func (_expecter _IdleClientMock_Expecter) IdleWithFallback(_a0 match.Arg[<-chan struct{}], _a1 match.Arg[time.Duration]) _IdleClientMock_IdleWithFallback_Call {
	return _IdleClientMock_IdleWithFallback_Call{Call: _expecter.mock.ExpectCall("IdleWithFallback", _a0.Matcher, _a1.Matcher)}
}

func (_call _IdleClientMock_IdleWithFallback_Call) Return(_r0 error) _IdleClientMock_IdleWithFallback_Call {
	_call.Call.Return(_r0)
	return _call
}

func (_call _IdleClientMock_IdleWithFallback_Call) RunReturn(f func(<-chan struct{}, time.Duration) (error)) _IdleClientMock_IdleWithFallback_Call {
	_call.Call.RunReturn(f)
	return _call
}

type _IdleClientMock_Logout_Call struct {
	*mock.Call
}

func (_mock *IdleClientMock) Logout() (error) {
	_mock.mock.T.Helper()
	_results := _mock.mock.Called("Logout", )
	_r0 := _results.Error(0)
	return _r0
}

func (_expecter _IdleClientMock_Expecter) Logout() _IdleClientMock_Logout_Call {
	return _IdleClientMock_Logout_Call{Call: _expecter.mock.ExpectCall("Logout", )}
}

func (_call _IdleClientMock_Logout_Call) Return(_r0 error) _IdleClientMock_Logout_Call {
	_call.Call.Return(_r0)
	return _call
}

func (_call _IdleClientMock_Logout_Call) RunReturn(f func() (error)) _IdleClientMock_Logout_Call {
	_call.Call.RunReturn(f)
	return _call
}

type _IdleClientMock_Select_Call struct {
	*mock.Call
}

func (_mock *IdleClientMock) Select(_a0 string, _a1 bool) (*goimap.MailboxStatus, error) {
	_mock.mock.T.Helper()
	_results := _mock.mock.Called("Select", _a0, _a1)
	var _r0 *goimap.MailboxStatus
	if _got := _results.Get(0); _got != nil {
		_r0 = _got.(*goimap.MailboxStatus)
	}
	_r1 := _results.Error(1)
	return _r0, _r1
}

func (_expecter _IdleClientMock_Expecter) Select(_a0 match.Arg[string], _a1 match.Arg[bool]) _IdleClientMock_Select_Call {
	return _IdleClientMock_Select_Call{Call: _expecter.mock.ExpectCall("Select", _a0.Matcher, _a1.Matcher)}
}

func (_call _IdleClientMock_Select_Call) Return(_r0 *goimap.MailboxStatus, _r1 error) _IdleClientMock_Select_Call {
	_call.Call.Return(_r0, _r1)
	return _call
}

func (_call _IdleClientMock_Select_Call) RunReturn(f func(string, bool) (*goimap.MailboxStatus, error)) _IdleClientMock_Select_Call {
	_call.Call.RunReturn(f)
	return _call
}

type _IdleClientMock_SetUpdates_Call struct {
	*mock.Call
}

func (_mock *IdleClientMock) SetUpdates(_a0 chan<- client.Update) () {
	_mock.mock.T.Helper()
	_mock.mock.Called("SetUpdates", _a0)
}

func (_expecter _IdleClientMock_Expecter) SetUpdates(_a0 match.Arg[chan<- client.Update]) _IdleClientMock_SetUpdates_Call {
	return _IdleClientMock_SetUpdates_Call{Call: _expecter.mock.ExpectCall("SetUpdates", _a0.Matcher)}
}

func (_call _IdleClientMock_SetUpdates_Call) Return() _IdleClientMock_SetUpdates_Call {
	_call.Call.Return()
	return _call
}

func (_call _IdleClientMock_SetUpdates_Call) RunReturn(f func(chan<- client.Update) ()) _IdleClientMock_SetUpdates_Call {
	_call.Call.RunReturn(f)
	return _call
}