package promise

import "reflect"

type Promise struct {
	isFullfilled bool
	value        interface{}

	isRejected bool
	err        error

	next         []*Promise
	dependencies []*Promise

	source func(resolve func(interface{}), reject func(error))

	parent      *Promise
	onFulfilled func(value interface{}) interface{}
	onRejected  func(value error) interface{}
}

// New returns a new Promise, which is very much like the JavaScript equivalent, but with one
// exception: the function given to New should not actually perform any work asynchronously (or if
// it does, it should be done transparently). The function will be invoked when the promise is
// scheduled. If the promise cannot be fulfilled yet, simply don't invoke resolve until the next
// time the function is called.
func New(f func(resolve func(interface{}), reject func(error))) *Promise {
	return &Promise{
		source: f,
	}
}

type Thenable interface {
	Then(f func(value interface{}) interface{}) *Promise
}

// Then appends a handler to be promise, invoking it when the promise is fulfilled. If the handler
// returns a value, it'll be passed as input to the next handler in the chain. If the handler
// returns another promise, the next handler in the chain will receive that promise's value when it
// is fulfilled.
func (p *Promise) Then(onFulfilled func(value interface{}) interface{}) *Promise {
	newPromise := &Promise{
		parent:      p,
		onFulfilled: onFulfilled,
	}
	p.next = append(p.next, newPromise)
	return newPromise
}

// Returns a Promise and deals with rejected cases only. The Promise returned by Catch is rejected
// if onRejected throws an error or returns a Promise which is itself rejected; otherwise, it is
// resolved.
func (p *Promise) Catch(onRejected func(err error) interface{}) *Promise {
	newPromise := &Promise{
		parent:     p,
		onRejected: onRejected,
	}
	p.next = append(p.next, newPromise)
	return newPromise
}

// Schedule invokes pending functions for unfulfilled promises and returns true if any progress was
// made.
func (p *Promise) Schedule() (didProgress bool) {
	for i := 0; ; i++ {
		didProgress := false
		for _, dependency := range p.dependencies {
			if !dependency.isFullfilled {
				if dependency.Schedule() {
					didProgress = true
				}
			}
		}
		if !p.isFullfilled && !p.isRejected {
			if p.source != nil {
				p.source(func(value interface{}) {
					p.isFullfilled = true
					p.value = value
					didProgress = true
				}, func(err error) {
					p.isRejected = true
					p.err = err
					didProgress = true
				})
			} else if p.parent != nil {
				if p.parent.isFullfilled || p.parent.isRejected {
					if p.parent.isFullfilled {
						p.isFullfilled = true
						p.value = p.parent.value
						if p.onFulfilled != nil {
							p.value = p.onFulfilled(p.value)
						}
					} else {
						if p.onRejected != nil {
							p.isFullfilled = true
							p.value = p.onRejected(p.parent.err)
						} else {
							p.isRejected = true
							p.err = p.parent.err
						}
					}
					didProgress = true
					if promise, ok := p.value.(*Promise); ok {
						for _, next := range p.next {
							next.parent = promise
						}
						promise.next = append(promise.next, p.next...)
						p.next = []*Promise{promise}
					}
				} else {
					didProgress = p.parent.Schedule()
				}
			}
		}
		if p.isFullfilled {
			for _, next := range p.next {
				if next.Schedule() {
					didProgress = true
				}
			}
		}
		if !didProgress {
			return i > 0
		}
	}
}

// Returns a Promise object that is resolved with the given value. If the value is a Thenable (i.e.
// has a Then method), the returned promise will "follow" that thenable, adopting its eventual
// state; otherwise the returned promise will be fulfilled with the value.
func Resolve(value interface{}) *Promise {
	if thenable, ok := value.(Thenable); ok {
		return thenable.Then(func(value interface{}) interface{} {
			return value
		})
	}
	return New(func(resolve func(interface{}), reject func(error)) {
		resolve(value)
	})
}

// Returns a Promise that is rejected with the given reason.
func Reject(reason error) *Promise {
	return New(func(resolve func(interface{}), reject func(error)) {
		reject(reason)
	})
}

// All returns a single Promise that resolves when all of the promises in the argument have resolved
// or when the iterable argument contains no promises. It rejects with the reason of the first
// promise that rejects.
func All(iterable interface{}) *Promise {
	v := reflect.ValueOf(iterable)
	result := make([]interface{}, v.Len())
	var rejectReason error
	remaining := 0
	all := New(func(resolve func(interface{}), reject func(error)) {
		if rejectReason != nil {
			reject(rejectReason)
		} else if remaining == 0 {
			resolve(result)
		}
	})
	for i := 0; i < v.Len(); i++ {
		value := v.Index(i).Interface()
		promise, ok := value.(*Promise)
		if !ok {
			result[i] = value
			continue
		} else if promise == nil {
			continue
		}
		i := i
		remaining++
		all.dependencies = append(all.dependencies, promise.Then(func(value interface{}) interface{} {
			result[i] = value
			remaining--
			return nil
		}).Catch(func(err error) interface{} {
			rejectReason = err
			return nil
		}))
	}
	return all
}
