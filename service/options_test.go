package service

import (
	"testing"
	"unsafe"
)

/*
	There's a good reason not to compare functions.
	If the functions don't have unique identities themselves we don't need to allocate memory for them.
	Also, the compiler may choose to inline the function making the same function have more than one address
	or no callable address at all. The spec prohibits function comparison to be able to optimize your code.
	It's a good thing and you should embrace it. 
	But still comparing the pointers to the functions is a good way to check if the functions are the same.
*/

func TestWithPreHooks(t *testing.T) {
	hook1 := func() {}
	hook2 := func() {}
	opts := &serviceOptions{}
	option := WithPreHooks(hook1, hook2)
	option(opts)

	if len(opts.preHooks) != 2 {
		t.Errorf("expected 2 pre-hooks, got %d", len(opts.preHooks))
	}

	// check pre-hooks are set correctly
	if *(*unsafe.Pointer)(unsafe.Pointer(&opts.preHooks[0])) != *(*unsafe.Pointer)(unsafe.Pointer(&hook1)) {
		t.Errorf("pre-hooks not set correctly")
	}

	// check pre-hooks are set correctly
	if *(*unsafe.Pointer)(unsafe.Pointer(&opts.preHooks[1])) != *(*unsafe.Pointer)(unsafe.Pointer(&hook2)) {
		t.Errorf("pre-hooks not set correctly")
	}
}

func TestWithIgnorePreRunHooksError(t *testing.T) {
	opts := &serviceOptions{}
	option := WithIgnorePreRunHooksError(true)
	option(opts)

	if !opts.ignorePreRunHooksError {
		t.Errorf("expected ignorePreRunHooksError to be true, got false")
	}
}

func TestWithPostHooks(t *testing.T) {
	hook1 := func() {}
	hook2 := func() {}
	opts := &serviceOptions{}
	option := WithPostHooks(hook1, hook2)
	option(opts)

	if len(opts.postHooks) != 2 {
		t.Errorf("expected 2 post-hooks, got %d", len(opts.postHooks))
	}

	// check post-hooks are set correctly
	if *(*unsafe.Pointer)(unsafe.Pointer(&opts.postHooks[0])) != *(*unsafe.Pointer)(unsafe.Pointer(&hook1)) {
		t.Errorf("post-hooks not set correctly")
	}

	// check post-hooks are set correctly
	if *(*unsafe.Pointer)(unsafe.Pointer(&opts.postHooks[1])) != *(*unsafe.Pointer)(unsafe.Pointer(&hook2)) {
		t.Errorf("post-hooks not set correctly")
	}
}

func TestWithIgnorePostRunHooksError(t *testing.T) {
	opts := &serviceOptions{}
	option := WithIgnorePostRunHooksError(true)
	option(opts)

	if !opts.ignorePostRunHooksError {
		t.Errorf("expected ignorePostRunHooksError to be true, got false")
	}
}
