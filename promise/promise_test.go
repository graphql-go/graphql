package promise

import (
	"errors"
	"testing"
)

func TestPromiseValueChaining(t *testing.T) {
	n := 0
	Resolve(1).Then(func(v interface{}) interface{} {
		n = v.(int)
		if n != 1 {
			t.Fatalf("expected 1, got %v", n)
		}
		return n + 1
	}).Then(func(v interface{}) interface{} {
		n = v.(int)
		if n != 2 {
			t.Fatalf("expected 2, got %v", n)
		}
		return n + 1
	}).Then(func(v interface{}) interface{} {
		n = v.(int)
		if n != 3 {
			t.Fatalf("expected 3, got %v", n)
		}
		return nil
	}).Schedule()
	if n != 3 {
		t.Fatalf("expected 3, got %v", n)
	}
}

func TestCatch(t *testing.T) {
	var val interface{}
	var err error
	Resolve(1).Then(func(interface{}) interface{} {
		return Reject(errors.New("reject"))
	}).Then(func(interface{}) interface{} {
		return nil
	}).Catch(func(caught error) interface{} {
		err = caught
		return "foo"
	}).Then(func(value interface{}) interface{} {
		val = value
		return nil
	}).Schedule()
	if err == nil {
		t.Fatalf("expected non-nil error")
	}
	if val != "foo" {
		t.Fatalf("expected \"foo\", got %v", val)
	}
}

func TestAll(t *testing.T) {
	p1 := Resolve(1)
	var p2 *Promise
	p3 := Resolve(3)
	var result []interface{}
	All([]interface{}{p1, p2, p3, 4}).Then(func(value interface{}) interface{} {
		result = value.([]interface{})
		return nil
	}).Schedule()
	if len(result) != 4 {
		t.Fatalf("expected 4 results, got %v", len(result))
	}
	if n, _ := result[0].(int); n != 1 {
		t.Fatalf("expected 1, got %v", n)
	}
	if result[1] != nil {
		t.Fatalf("expected nil, got %v", result[1])
	}
	if n, _ := result[2].(int); n != 3 {
		t.Fatalf("expected 3, got %v", n)
	}
	if n, _ := result[3].(int); n != 4 {
		t.Fatalf("expected 4, got %v", n)
	}
}

func TestAll_Reject(t *testing.T) {
	p1 := Resolve(1)
	p2 := Reject(errors.New("foo"))
	p3 := Resolve(3)
	var result []interface{}
	var rejectReason error
	All([]interface{}{p1, p2, p3, 4}).Then(func(value interface{}) interface{} {
		result = value.([]interface{})
		return nil
	}).Catch(func(err error) interface{} {
		rejectReason = err
		return nil
	}).Schedule()
	if result != nil {
		t.Fatalf("expected nil result, got %v", result)
	}
	if rejectReason == nil {
		t.Fatalf("expected non-nil reject reason")
	}
}

func TestAll_Schedule(t *testing.T) {
	step := 0
	p1 := New(func(resolve func(interface{}), reject func(error)) {
		if step > 1 {
			resolve(1)
		}
	})
	p2 := New(func(resolve func(interface{}), reject func(error)) {
		if step > 0 {
			resolve(2)
		}
	})
	var result []interface{}
	all := All([]interface{}{p1, p2}).Then(func(value interface{}) interface{} {
		result = value.([]interface{})
		return nil
	})

	if all.Schedule() {
		t.Fatalf("expected false")
	}

	step++
	if !all.Schedule() {
		t.Fatalf("expected true")
	}

	step++
	if !all.Schedule() {
		t.Fatalf("expected true")
	}

	if len(result) != 2 {
		t.Fatalf("expected 2 results, got %v", len(result))
	}
	if n, _ := result[0].(int); n != 1 {
		t.Fatalf("expected 1, got %v", n)
	}
	if n, _ := result[1].(int); n != 2 {
		t.Fatalf("expected 2, got %v", n)
	}
}
