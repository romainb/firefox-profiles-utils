// Code generated by mockery v2.9.4. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"

	regexp "regexp"
)

// FirefoxProfiles is an autogenerated mock type for the FirefoxProfiles type
type FirefoxProfiles struct {
	mock.Mock
}

// GetProfilesList provides a mock function with given fields:
func (_m *FirefoxProfiles) GetProfilesList() ([]string, error) {
	ret := _m.Called()

	var r0 []string
	if rf, ok := ret.Get(0).(func() []string); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProfilesMatching provides a mock function with given fields: regex
func (_m *FirefoxProfiles) GetProfilesMatching(regex *regexp.Regexp) ([]string, error) {
	ret := _m.Called(regex)

	var r0 []string
	if rf, ok := ret.Get(0).(func(*regexp.Regexp) []string); ok {
		r0 = rf(regex)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*regexp.Regexp) error); ok {
		r1 = rf(regex)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetProfilesPath provides a mock function with given fields:
func (_m *FirefoxProfiles) GetProfilesPath() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// IsProfileUsed provides a mock function with given fields: profileName
func (_m *FirefoxProfiles) IsProfileUsed(profileName string) (bool, error) {
	ret := _m.Called(profileName)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(profileName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(profileName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}