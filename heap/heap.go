// Copy-Pasted from the standard library
// This is a modification to standard library before it provides a standard generic based library.
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package heap provides heap operations for any type that implements
// heap.Interface. A heap is a tree with the property that each node is the
// minimum-valued node in its subtree.
//
// The minimum element in the tree is the root, at index 0.
//
// A heap is a common way to implement a priority queue. To build a priority
// queue, implement the Heap interface with the (negative) priority as the
// ordering for the Less method, so Push adds items while Pop removes the
// highest-priority item from the queue. The Examples include such an
// implementation; the file example_pq_test.go has the complete source.
//
package heap

import "sort"

// The Interface type describes the requirements
// for a type using the routines in this package.
// Any type that implements it may be used as a
// min-heap with the following invariants (established after
// Init has been called or if the data is empty or sorted):
//
//	!h.Less(j, i) for 0 <= i < h.Len() and 2*i+1 <= j <= 2*i+2 and j < h.Len()
//
// Note that Push and Pop in this interface are for package heap's
// implementation to call. To add and remove things from the heap,
// use heap.Push and heap.Pop.
type Interface[T any] interface {
	sort.Interface
	Push(x T) // add x as element Len()
	Pop() T   // remove and return element Len() - 1.
}

// Init establishes the heap invariants required by the other routines in this package.
// Init is idempotent with respect to the heap invariants
// and may be called whenever the heap invariants may have been invalidated.
// The complexity is O(n) where n = h.Len().
func Init[H Interface[T], T any](h H) {
	// heapify
	n := h.Len()
	for i := n/2 - 1; i >= 0; i-- {
		down[H, T](h, i, n)
	}
}

// Push pushes the element x onto the heap.
// The complexity is O(log n) where n = h.Len().
func Push[H Interface[T], T any](h H, x T) {
	h.Push(x)
	up[H, T](h, h.Len()-1)
}

// Pop removes and returns the minimum element (according to Less) from the heap.
// The complexity is O(log n) where n = h.Len().
// Pop is equivalent to Remove(h, 0).
func Pop[H Interface[T], T any](h H) T {
	n := h.Len() - 1
	h.Swap(0, n)
	down[H, T](h, 0, n)
	return h.Pop()
}

// Remove removes and returns the element at index i from the heap.
// The complexity is O(log n) where n = h.Len().
func Remove[H Interface[T], T any](h H, i int) T {
	n := h.Len() - 1
	if n != i {
		h.Swap(i, n)
		if !down[H, T](h, i, n) {
			up[H, T](h, i)
		}
	}
	return h.Pop()
}

// Fix re-establishes the heap ordering after the element at index i has changed its value.
// Changing the value of the element at index i and then calling Fix is equivalent to,
// but less expensive than, calling Remove(h, i) followed by a Push of the new value.
// The complexity is O(log n) where n = h.Len().
func Fix[H Interface[T], T any](h H, i int) {
	if !down[H, T](h, i, h.Len()) {
		up[H, T](h, i)
	}
}

func up[H Interface[T], T any](h H, j int) {
	for {
		i := (j - 1) / 2 // parent
		if i == j || !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		j = i
	}
}

func down[H Interface[T], T any](h H, i0, n int) bool {
	i := i0
	for {
		j1 := 2*i + 1
		if j1 >= n || j1 < 0 { // j1 < 0 after int overflow
			break
		}
		j := j1 // left child
		if j2 := j1 + 1; j2 < n && h.Less(j2, j1) {
			j = j2 // = 2*i + 2  // right child
		}
		if !h.Less(j, i) {
			break
		}
		h.Swap(i, j)
		i = j
	}
	return i > i0
}

// IsHeap checks if the input is a properly heap-ified.
func IsHeap[H Interface[T], T any](h H) bool {
loop:
	for i := 0; i < h.Len(); i++ {
		lc := i*2 + 1
		if lc >= h.Len() || lc < 0 {
			break loop
		}
		if !h.Less(i, lc) {
			return false
		}
		rc := i*2 + 2
		if rc >= h.Len() || rc < 0 {
			break loop
		}
		if !h.Less(i, rc) {
			return false
		}
	}

	return true
}
