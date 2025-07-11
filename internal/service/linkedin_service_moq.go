// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/you/linkedinify/internal/model"
	"sync"
)

// Ensure, that LinkedInServiceInteractorMock does implement LinkedInServiceInteractor.
// If this is not the case, regenerate this file with moq.
var _ LinkedInServiceInteractor = &LinkedInServiceInteractorMock{}

// LinkedInServiceInteractorMock is a mock implementation of LinkedInServiceInteractor.
//
//	func TestSomethingThatUsesLinkedInServiceInteractor(t *testing.T) {
//
//		// make and configure a mocked LinkedInServiceInteractor
//		mockedLinkedInServiceInteractor := &LinkedInServiceInteractorMock{
//			HistoryFunc: func(ctx context.Context, userID uuid.UUID, page int, pageSize int) ([]model.LinkedInPost, error) {
//				panic("mock out the History method")
//			},
//			TransformFunc: func(ctx context.Context, userID uuid.UUID, text string) (string, error) {
//				panic("mock out the Transform method")
//			},
//		}
//
//		// use mockedLinkedInServiceInteractor in code that requires LinkedInServiceInteractor
//		// and then make assertions.
//
//	}
type LinkedInServiceInteractorMock struct {
	// HistoryFunc mocks the History method.
	HistoryFunc func(ctx context.Context, userID uuid.UUID, page int, pageSize int) ([]model.LinkedInPost, error)

	// TransformFunc mocks the Transform method.
	TransformFunc func(ctx context.Context, userID uuid.UUID, text string) (string, error)

	// calls tracks calls to the methods.
	calls struct {
		// History holds details about calls to the History method.
		History []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserID is the userID argument value.
			UserID uuid.UUID
			// Page is the page argument value.
			Page int
			// PageSize is the pageSize argument value.
			PageSize int
		}
		// Transform holds details about calls to the Transform method.
		Transform []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserID is the userID argument value.
			UserID uuid.UUID
			// Text is the text argument value.
			Text string
		}
	}
	lockHistory   sync.RWMutex
	lockTransform sync.RWMutex
}

// History calls HistoryFunc.
func (mock *LinkedInServiceInteractorMock) History(ctx context.Context, userID uuid.UUID, page int, pageSize int) ([]model.LinkedInPost, error) {
	if mock.HistoryFunc == nil {
		panic("LinkedInServiceInteractorMock.HistoryFunc: method is nil but LinkedInServiceInteractor.History was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		UserID   uuid.UUID
		Page     int
		PageSize int
	}{
		Ctx:      ctx,
		UserID:   userID,
		Page:     page,
		PageSize: pageSize,
	}
	mock.lockHistory.Lock()
	mock.calls.History = append(mock.calls.History, callInfo)
	mock.lockHistory.Unlock()
	return mock.HistoryFunc(ctx, userID, page, pageSize)
}

// HistoryCalls gets all the calls that were made to History.
// Check the length with:
//
//	len(mockedLinkedInServiceInteractor.HistoryCalls())
func (mock *LinkedInServiceInteractorMock) HistoryCalls() []struct {
	Ctx      context.Context
	UserID   uuid.UUID
	Page     int
	PageSize int
} {
	var calls []struct {
		Ctx      context.Context
		UserID   uuid.UUID
		Page     int
		PageSize int
	}
	mock.lockHistory.RLock()
	calls = mock.calls.History
	mock.lockHistory.RUnlock()
	return calls
}

// Transform calls TransformFunc.
func (mock *LinkedInServiceInteractorMock) Transform(ctx context.Context, userID uuid.UUID, text string) (string, error) {
	if mock.TransformFunc == nil {
		panic("LinkedInServiceInteractorMock.TransformFunc: method is nil but LinkedInServiceInteractor.Transform was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		UserID uuid.UUID
		Text   string
	}{
		Ctx:    ctx,
		UserID: userID,
		Text:   text,
	}
	mock.lockTransform.Lock()
	mock.calls.Transform = append(mock.calls.Transform, callInfo)
	mock.lockTransform.Unlock()
	return mock.TransformFunc(ctx, userID, text)
}

// TransformCalls gets all the calls that were made to Transform.
// Check the length with:
//
//	len(mockedLinkedInServiceInteractor.TransformCalls())
func (mock *LinkedInServiceInteractorMock) TransformCalls() []struct {
	Ctx    context.Context
	UserID uuid.UUID
	Text   string
} {
	var calls []struct {
		Ctx    context.Context
		UserID uuid.UUID
		Text   string
	}
	mock.lockTransform.RLock()
	calls = mock.calls.Transform
	mock.lockTransform.RUnlock()
	return calls
}
