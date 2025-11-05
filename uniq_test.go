package main

import "testing"

func TestCheckLines(t *testing.T) {
	// Простые тесты в одну строку
	if got := checkLines("Hello", Params{}); got != "Hello" {
		t.Errorf("want Hello, got %s", got)
	}
	if got := checkLines("HELLO", Params{I: true}); got != "hello" {
		t.Errorf("want hello, got %s", got)
	}
	if got := checkLines("a b c", Params{Num_fields: 1}); got != "b c" {
		t.Errorf("want 'b c', got %s", got)
	}
	if got := checkLines("abc", Params{Num_chars: 1}); got != "bc" {
		t.Errorf("want 'bc', got %s", got)
	}
}

func TestCheckCDU(t *testing.T) {
	// Базовые тесты
	lines := []string{"a", "a", "b"}

	// Без флагов
	result := checkCDU(lines, Params{})
	if len(result) != 2 || result[0] != "a" || result[1] != "b" {
		t.Errorf("basic unique failed: %v", result)
	}

	// С счетчиком
	result = checkCDU(lines, Params{C: true})
	if len(result) != 2 || result[0] != "2 a" || result[1] != "1 b" {
		t.Errorf("count failed: %v", result)
	}

	// Только дубликаты
	result = checkCDU([]string{"a", "a", "b"}, Params{D: true})
	if len(result) != 1 || result[0] != "a" {
		t.Errorf("duplicates failed: %v", result)
	}

	// Только уникальные
	result = checkCDU([]string{"a", "a", "b"}, Params{U: true})
	if len(result) != 1 || result[0] != "b" {
		t.Errorf("unique failed: %v", result)
	}
}

func TestUniqLinesComplex(t *testing.T) {
	// Реальный пример из задания
	lines := []string{
		"I love music.",
		"I love music.",
		"I love music.",
		"",
		"I love music of Kartik.",
		"I love music of Kartik.",
		"Thanks.",
	}

	// Без параметров
	result := checkCDU(lines, Params{})
	expected := []string{"I love music.", "", "I love music of Kartik.", "Thanks."}

	if len(result) != len(expected) {
		t.Fatalf("length mismatch: got %d, want %d", len(result), len(expected))
	}
	for i := range expected {
		if result[i] != expected[i] {
			t.Errorf("at %d: got %q, want %q", i, result[i], expected[i])
		}
	}
}
