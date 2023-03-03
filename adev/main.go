package main

import (
	"context"
	"fmt"
	"github.com/adshao/go-binance/v2"
	"github.com/sirupsen/logrus"
	"github.com/uvite/gvm/engine"
	"github.com/uvite/gvm/tart/fixedpoint"
	"github.com/uvite/gvm/tart/floats"
	"github.com/uvite/gvm/tart/types"
	"syscall"

	"os"
	"os/signal"
	"time"
)

var (
	apiKey    = "your api key"
	secretKey = "your secret key"
	interval  = "1m"
	symbol    = "ETHUSDT"
)

func main() {

	gvm, _ := engine.NewGvm()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	gvm.Ctx = ctx

	err := gvm.Load("jma.js")
	fmt.Println(err)
	//gvm.Run()
	//
	////closes := &floats.Slice{}
	////high := &floats.Slice{}
	//////b.Set("high", high)
	////low := &floats.Slice{}
	//////b.Set("low", low)
	//rtOpts := lib.RuntimeOptions{}
	////b.Set("close", close)
	//r, err := genv.GetSimpleRunner(b, "/script.js", fmt.Sprintf(`
	//				%s
	//		`, sourceData.Data), fs, rtOpts)
	//
	//fmt.Println(err)
	//ch := make(chan metrics.SampleContainer, 100)
	//defer close(ch)
	//go func() { // read the channel so it doesn't block
	//	for range ch {
	//	}
	//}()
	//
	//

	//initVU, err := r.NewVU(ctx,1, 1, ch)
	//vu := initVU.Activate(&lib.VUActivationParams{RunContext: ctx})

	//go func() {
	//	kline(gvm)
	//}()
	kline(gvm)
	WaitForSignal(ctx, syscall.SIGINT, syscall.SIGTERM)
	//

	//for i := 0; i < 5; i++ {
	//err = vu.RunOnce()
	//fmt.Println(err)
	//}

	//var client = binance.NewClient(apiKey, secretKey)
	//
	//klines, err := client.NewKlinesService().Symbol(symbol).
	//	Interval(interval).Do(context.Background())
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//for _, event := range klines {
	//	//fmt.Println(event.CloseTime)
	//	time1 := types.NewTimeFromUnix(0, event.CloseTime*int64(time.Millisecond))
	//	isClosed := time.Now().After(time1.Time())
	//	if isClosed {
	//		closes.Push(fixedpoint.MustNewFromString(event.Close).Float64())
	//		high.Push(fixedpoint.MustNewFromString(event.High).Float64())
	//		low.Push(fixedpoint.MustNewFromString(event.Low).Float64())
	//	}
	//
	//	//rsx.Push(fixedpoint.MustNewFromString(event.Low).Float64())
	//
	//}
	////if call, ok := goja.AssertFunction(exports.Get("default")); ok {
	////	if _, err = call(goja.Undefined()); err != nil {
	////
	////	}
	////}
	//
	//err = vu.RunOnce()
	//fmt.Println(err)
	//wsKlineHandler := func(event *binance.WsKlineEvent) {
	//
	//	if event.Kline.IsFinal {
	//
	//		closes.Push(fixedpoint.MustNewFromString(event.Kline.Close).Float64())
	//		high.Push(fixedpoint.MustNewFromString(event.Kline.High).Float64())
	//		low.Push(fixedpoint.MustNewFromString(event.Kline.Low).Float64())
	//
	//		//if call, ok := goja.AssertFunction(exports.Get("default")); ok {
	//		//	if _, err = call(goja.Undefined()); err != nil {
	//		//
	//		//	}
	//		//}
	//		err = vu.RunOnce()
	//		fmt.Println(err)
	//
	//	}
	//}
	//errHandler := func(err error) {
	//	fmt.Println(err)
	//}
	//doneC, _, err := binance.WsKlineServe(symbol, interval, wsKlineHandler, errHandler)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//<-doneC
}

func WaitForSignal(ctx context.Context, signals ...os.Signal) os.Signal {
	var sigC = make(chan os.Signal, 1)
	signal.Notify(sigC, signals...)
	defer signal.Stop(sigC)

	select {
	case sig := <-sigC:
		logrus.Warnf("%v", sig)
		return sig

	case <-ctx.Done():
		return nil

	}

	return nil
}

func kline(gvm *engine.Gvm) {

	closes := &floats.Slice{}
	high := &floats.Slice{}

	low := &floats.Slice{}
	gvm.Set("close", closes)
	gvm.Set("high", high)
	gvm.Set("low", low)

	var client = binance.NewClient(apiKey, secretKey)

	klines, err := client.NewKlinesService().Symbol(symbol).
		Interval(interval).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, event := range klines {
		//fmt.Println(event.CloseTime)
		time1 := types.NewTimeFromUnix(0, event.CloseTime*int64(time.Millisecond))
		isClosed := time.Now().After(time1.Time())
		if isClosed {
			closes.Push(fixedpoint.MustNewFromString(event.Close).Float64())
			high.Push(fixedpoint.MustNewFromString(event.High).Float64())
			low.Push(fixedpoint.MustNewFromString(event.Low).Float64())
		}

	}

	wsKlineHandler := func(event *binance.WsKlineEvent) {

		if event.Kline.IsFinal {

			closes.Push(fixedpoint.MustNewFromString(event.Kline.Close).Float64())
			high.Push(fixedpoint.MustNewFromString(event.Kline.High).Float64())
			low.Push(fixedpoint.MustNewFromString(event.Kline.Low).Float64())
			fmt.Println(event.Kline.EndTime)
			_, err = gvm.Run()
			fmt.Println(err)

		}
	}
	errHandler := func(err error) {
		fmt.Println(err)
	}
	doneC, _, err := binance.WsKlineServe(symbol, interval, wsKlineHandler, errHandler)
	if err != nil {
		fmt.Println(err)
		return
	}
	<-doneC

}
