package validator

import (
	"testing"

	"snippet.darieldejesus.com/internal/assert"
)

func TestValidator(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		validator := &Validator{}
		assert.Equal(t, validator.Valid(), true)
	})

	t.Run("AddFieldError", func(t *testing.T) {
		key := "field"
		val := "This is an error"
		validator := &Validator{}
		validator.AddFieldError(key, val)
		assert.Equal(t, validator.Valid(), false)
		assert.Equal(t, validator.FieldErrors[key], val)
	})

	t.Run("AddNonFieldError", func(t *testing.T) {
		val := "This is an error"
		validator := &Validator{}
		validator.AddNonFieldError(val)
		assert.Equal(t, validator.Valid(), false)
		assert.Equal(t, validator.NonFieldErrors[0], val)
	})

	t.Run("CheckField", func(t *testing.T) {
		key := "field"
		val := "This is an error"
		validator := &Validator{}
		validator.CheckField(true, key, val)
		assert.Equal(t, validator.FieldErrors[key], "")
		validator.CheckField(false, key, val)
		assert.Equal(t, validator.FieldErrors[key], val)
	})
}

func TestNotBlank(t *testing.T) {
	assert.Equal(t, NotBlank(""), false)
	assert.Equal(t, NotBlank("value"), true)
}

func TestMaxChars(t *testing.T) {
	assert.Equal(t, MaxChars("Platano power", 15), true)
	assert.Equal(t, MaxChars("Platano power was here", 15), false)
	assert.Equal(t, MaxChars("Dariel wæs hërë", 15), true)
}

func TestPermittedValue(t *testing.T) {
	assert.Equal(t, PermittedValue(5, 0, 5, 10), true)
	assert.Equal(t, PermittedValue(5, 5), true)
	assert.Equal(t, PermittedValue(0, 0), true)
	assert.Equal(t, PermittedValue(-5, -5), true)
	assert.Equal(t, PermittedValue(5, 0, 6, 12), false)
	assert.Equal(t, PermittedValue(5, -5), false)
	assert.Equal(t, PermittedValue(-5, 5), false)
	assert.Equal(t, PermittedValue(5), false)
}

func TestMinChars(t *testing.T) {
	assert.Equal(t, MinChars("Platano power", 10), true)
	assert.Equal(t, MinChars("Platano", 10), false)
	assert.Equal(t, MinChars("Dariel wæ", 10), false)
	assert.Equal(t, MinChars("Dariel wæs", 10), true)
}

func TestMatches(t *testing.T) {
	assert.Equal(t, Matches("contact@darieldejesus.com", EmailRX), true)
	assert.Equal(t, Matches("contact@darieldejesus.", EmailRX), false)
}
