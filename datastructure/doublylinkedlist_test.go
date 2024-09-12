package datastructure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDoublyLinked_AddToHead(t *testing.T) {
	dll := NewDoublyLinked[int]()
	if !dll.IsEmpty() {
		t.Error("Expected list to be empty initially")
	}

	dll.AddToHead(1)
	dll.AddToHead(2)

	if dll.Count() != 2 {
		t.Errorf("Expected count to be 2, got %d", dll.Count())
	}

	if dll.Head().Value != 2 {
		t.Errorf("Expected head Value to be 2, got %d", dll.Head().Value)
	}

	if dll.Tail().Value != 1 {
		t.Errorf("Expected tail Value to be 1, got %d", dll.Tail().Value)
	}
}

func TestDoublyLinked_AddToTail(t *testing.T) {
	dll := NewDoublyLinked[int]()
	if !dll.IsEmpty() {
		t.Error("Expected list to be empty initially")
	}
	assert.Nil(t, dll.Tail())

	dll.AddToTail(1)
	dll.AddToTail(2)

	if dll.Count() != 2 {
		t.Errorf("Expected count to be 2, got %d", dll.Count())
	}

	if dll.Head().Value != 1 {
		t.Errorf("Expected head Value to be 1, got %d", dll.Head().Value)
	}

	if dll.Tail().Value != 2 {
		t.Errorf("Expected tail Value to be 2, got %d", dll.Tail().Value)
	}
}

func TestDoublyLinked_RemoveFromHead(t *testing.T) {
	dll := NewDoublyLinked[int]()
	dll.AddToTail(1)
	dll.AddToTail(2)

	dll.RemoveFromHead()

	if dll.Count() != 1 {
		t.Errorf("Expected count to be 1 after removal, got %d", dll.Count())
	}

	if dll.Head().Value != 2 {
		t.Errorf("Expected head Value to be 2 after removal, got %d", dll.Head().Value)
	}

	dll.RemoveFromHead()
	if !dll.IsEmpty() {
		t.Error("Expected list to be empty after removing all elements")
	}
}

func TestDoublyLinked_RemoveFromTail(t *testing.T) {
	dll := NewDoublyLinked[int]()
	dll.AddToTail(1)
	dll.AddToTail(2)

	dll.RemoveFromTail()

	if dll.Count() != 1 {
		t.Errorf("Expected count to be 1 after removal, got %d", dll.Count())
	}

	if dll.Tail().Value != 1 {
		t.Errorf("Expected tail Value to be 1 after removal, got %d", dll.Tail().Value)
	}

	dll.RemoveFromTail()
	if !dll.IsEmpty() {
		t.Error("Expected list to be empty after removing all elements")
	}
}

func TestDoublyLinked_RemoveSpecificNode(t *testing.T) {
	dll := NewDoublyLinked[int]()
	node1 := dll.AddToTail(1)
	node2 := dll.AddToTail(2)

	dll.Remove(node1)

	if dll.Count() != 1 {
		t.Errorf("Expected count to be 1 after removal, got %d", dll.Count())
	}

	if dll.Head().Value != 2 {
		t.Errorf("Expected head Value to be 2 after removal, got %d", dll.Head().Value)
	}

	dll.Remove(node2)
	if !dll.IsEmpty() {
		t.Error("Expected list to be empty after removing all elements")
	}
}

func TestDoublyLinked_IsEmpty(t *testing.T) {
	dll := NewDoublyLinked[int]()
	if !dll.IsEmpty() {
		t.Error("Expected list to be empty initially")
	}
	dll.AddToTail(1)
	if dll.IsEmpty() {
		t.Error("Expected list to not be empty after adding an element")
	}
	dll.RemoveFromTail()
	if !dll.IsEmpty() {
		t.Error("Expected list to be empty after removing all elements")
	}
}

func TestDoublyLinked_Range(t *testing.T) {
	dll := NewDoublyLinked[int]()
	dll.AddToTail(1)
	dll.AddToTail(2)
	dll.AddToTail(3)

	result := []int{}
	dll.Range(func(val int) bool {
		result = append(result, val)
		return true
	})

	expected := []int{1, 2, 3}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected %d, got %d", v, result[i])
		}
	}
}

// Test for Range with early exit
func TestDoublyLinked_RangeWithEarlyExit(t *testing.T) {
	dll := NewDoublyLinked[int]()
	dll.AddToTail(1)
	dll.AddToTail(2)
	dll.AddToTail(3)

	result := []int{}
	dll.Range(func(val int) bool {
		result = append(result, val)
		return val != 2 // stop iteration early
	})

	expected := []int{1, 2}
	for i, v := range expected {
		if result[i] != v {
			t.Errorf("Expected %d, got %d", v, result[i])
		}
	}
}

// Test for various edge cases when removing nodes
func TestDoublyLinked_RemoveEdgeCases(t *testing.T) {
	// Edge case: Remove from an empty list
	t.Run("RemoveFromEmptyList", func(t *testing.T) {
		dll := NewDoublyLinked[int]()
		dll.Remove(dll.Head()) // Should not panic or change anything
		if dll.Count() != 0 {
			t.Errorf("Expected count to be 0, got %d", dll.Count())
		}
	})

	// Edge case: Remove a nil node
	t.Run("RemoveNilNode", func(t *testing.T) {
		dll := NewDoublyLinked[int]()
		dll.Remove(nil) // Should not panic or change anything
		if dll.Count() != 0 {
			t.Errorf("Expected count to be 0, got %d", dll.Count())
		}
	})

	// Edge case: Remove head node in a list with one element
	t.Run("RemoveHeadWithOneElement", func(t *testing.T) {
		dll := NewDoublyLinked[int]()
		dll.AddToHead(1)
		dll.Remove(dll.Head()) // Should remove the only element
		if !dll.IsEmpty() {
			t.Error("Expected list to be empty after removing the only element")
		}
	})

	// Edge case: Remove tail node in a list with one element
	t.Run("RemoveTailWithOneElement", func(t *testing.T) {
		dll := NewDoublyLinked[int]()
		dll.AddToTail(1)
		dll.Remove(dll.Tail()) // Should remove the only element
		if !dll.IsEmpty() {
			t.Error("Expected list to be empty after removing the only element")
		}
	})

	// Edge case: Remove head node directly in a larger list
	t.Run("RemoveHeadNodeInLargerList", func(t *testing.T) {
		dll := NewDoublyLinked[int]()
		dll.AddToHead(1)
		dll.AddToHead(2)
		dll.Remove(dll.Head()) // Should remove the head (Value 2)
		if dll.Head().Value != 1 {
			t.Errorf("Expected head Value to be 1, got %d", dll.Head().Value)
		}
	})

	// Edge case: Remove tail node directly in a larger list
	t.Run("RemoveTailNodeInLargerList", func(t *testing.T) {
		dll := NewDoublyLinked[int]()
		dll.AddToTail(3)
		dll.AddToTail(4)
		dll.Remove(dll.Tail()) // Should remove the tail (Value 4)
		if dll.Tail().Value != 3 {
			t.Errorf("Expected tail Value to be 3, got %d", dll.Tail().Value)
		}
	})

	// Edge case: Remove node that has already been removed
	t.Run("RemoveNodeThatWasAlreadyRemoved", func(t *testing.T) {
		dll := NewDoublyLinked[int]()
		node := dll.AddToTail(5)
		dll.Remove(node) // First removal
		dll.Remove(node) // Second removal, should do nothing
		if dll.Count() != 0 {
			t.Errorf("Expected count to be 0 after trying to remove the same node twice, got %d", dll.Count())
		}
	})
}

func TestDoublyLinked_MoveToHead(t *testing.T) {
	dll := NewDoublyLinked[int]()
	dll.AddToTail(1)
	node := dll.AddToTail(2)
	dll.AddToTail(3)
	dll.MoveToHead(node)
	if dll.Head().Value != 2 {
		t.Errorf("Expected head Value to be 2, got %d", dll.Head().Value)
	}
}

func TestDoublyLinked_MoveToTail(t *testing.T) {
	dll := NewDoublyLinked[int]()
	dll.AddToTail(1)
	node := dll.AddToTail(2)
	dll.AddToTail(3)
	dll.MoveToTail(node)
	if dll.Tail().Value != 2 {
		t.Errorf("Expected tail Value to be 2, got %d", dll.Tail().Value)
	}
}
