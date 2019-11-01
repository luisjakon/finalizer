package finalizer

import (
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

////
////  Finalizer -
//// 	This is a tiny utility package for ensuring user-defined functions must run when the program exits
//// 	even when the program is terminated via an os interrupt such as SIGTERM
////
////  Example
////  -------
////  	func main() {
////		defer finalizer.MustOnExit(func() { fmt.Println("I must run no matter what!") })()
////		...
////  	}
////
type finalizr func()

func init() {
	var doOnce sync.Once
	doOnce.Do(func() {
		signalchan := make(chan os.Signal)
		signal.Notify(signalchan, os.Interrupt, syscall.SIGTERM)
		go func() {
			<-signalchan
			runtime.GC() // <- required by runtime.SetFinalizer(...)
			os.Exit(int(syscall.SIGTERM))
		}()
	})
}

func MustOnExit(fn func()) func() {
	if fn == nil {
		return func() {}
	}
	f := finalizr(fn)
	runtime.SetFinalizer(&f, (*finalizr).OnExit)
	return fn
}

func (f *finalizr) OnExit() {
	(*f)()
}
