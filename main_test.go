package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyEdit_Insert(t *testing.T) {
	// --- Given ---
	text := "Hello World"
	edit := Edit{
		Position: 5,
		Text:     ",",
		Type:     "insert",
	}

	// --- When ---
	err := applyEdit(&text, edit)

	// --- Then ---
	assert.NoError(t, err)
	exp := "Hello, World"
	assert.Equal(t, exp, text)
}

func TestApplyEdit_Delete(t *testing.T) {
	// --- Given ---
	text := "Hello, World"
	edit := Edit{
		Position: 5,
		Length:   1,
		Type:     "delete",
	}

	// --- When ---
	err := applyEdit(&text, edit)

	// --- Then ---
	assert.NoError(t, err)
	exp := "Hello World"
	assert.Equal(t, exp, text)
}

func TestApplyEdit_InvalidPosition(t *testing.T) {
	// --- Given ---
	text := "Hello"
	edit := Edit{
		Position: -1,
		Text:     "Test",
		Type:     "insert",
	}

	// --- When ---
	err := applyEdit(&text, edit)

	// --- Then ---
	assert.Error(t, err)
}

func TestApplyEdit_InvalidLength(t *testing.T) {
	// --- Given ---
	text := "Hello"
	edit := Edit{
		Position: 0,
		Length:   10,
		Type:     "delete",
	}

	// --- When ---
	err := applyEdit(&text, edit)

	// --- Then ---
	assert.Error(t, err)
}

func TestApplyEdit_UnknownType(t *testing.T) {
	// --- Given ---
	text := "Hello"
	edit := Edit{
		Position: 0,
		Text:     "Test",
		Type:     "replace",
	}

	// --- When ---
	err := applyEdit(&text, edit)

	// --- Then ---
	assert.Error(t, err)
}
