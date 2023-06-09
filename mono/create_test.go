package mono

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/jjeffcaii/reactor-go"
)

func TestMonoCreate_SubscribeWith(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewMockSubscriber(ctrl)
	s.EXPECT().OnSubscribe(gomock.Any(), gomock.Any()).Do(MockRequestInfinite).Times(1)
	s.EXPECT().OnNext(gomock.Eq(1)).Times(1)
	s.EXPECT().OnComplete().Times(1)
	s.EXPECT().OnError(gomock.Any()).Times(0)

	newMonoCreate(func(ctx context.Context, s Sink) {
		s.Success(1)
	}).SubscribeWith(context.Background(), s)

	fakeErr := errors.New("fake error")

	s = NewMockSubscriber(ctrl)
	s.EXPECT().OnSubscribe(gomock.Any(), gomock.Any()).Do(MockRequestInfinite).Times(1)
	s.EXPECT().OnNext(gomock.Any()).Times(0)
	s.EXPECT().OnComplete().Times(0)
	s.EXPECT().OnError(gomock.Eq(fakeErr)).Times(1)

	newMonoCreate(func(ctx context.Context, s Sink) {
		s.Error(fakeErr)
	}).SubscribeWith(context.Background(), s)
}

func TestMonoCreate_Panic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewMockSubscriber(ctrl)

	s.EXPECT().OnSubscribe(gomock.Any(), gomock.Any()).Do(MockRequestInfinite).Times(2)
	s.EXPECT().OnNext(gomock.Any()).Times(0)
	s.EXPECT().OnError(gomock.Any()).Times(2)
	s.EXPECT().OnComplete().Times(0)

	newMonoCreate(func(ctx context.Context, s Sink) {
		panic("fake panic")
	}).SubscribeWith(context.Background(), s)

	newMonoCreate(func(ctx context.Context, s Sink) {
		panic(errors.New("fake error"))
	}).SubscribeWith(context.Background(), s)
}

func TestMonoCreate_OnNextPanic(t *testing.T) {
	var cnt int
	onError := reactor.OnError(func(e error) {
		cnt++
	})
	m := newMonoCreate(func(ctx context.Context, s Sink) {
		s.Success(1)
	})
	m.SubscribeWith(context.Background(),
		reactor.NewSubscriber(
			reactor.OnNext(func(v reactor.Any) error {
				panic("fake panic")
			}),
			onError,
		))
	m.SubscribeWith(context.Background(),
		reactor.NewSubscriber(
			reactor.OnNext(func(v reactor.Any) error {
				panic(errors.New("fake error"))
			}),
			onError,
		))
}

func TestMonoCreate_Cancel(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	s := NewMockSubscriber(ctrl)
	s.EXPECT().OnSubscribe(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, su reactor.Subscription) {
			su.Cancel()
		}).
		Times(1)
	s.EXPECT().OnNext(gomock.Any()).Times(0)
	s.EXPECT().OnComplete().Times(0)
	s.EXPECT().OnError(gomock.Eq(reactor.ErrSubscribeCancelled)).Times(1)

	newMonoCreate(func(ctx context.Context, s Sink) {
		s.Success(1)
	}).SubscribeWith(context.Background(), s)
}
