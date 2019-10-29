// Copyright 2018 The TCell Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use file except in compliance with the License.
// You may obtain a copy of the license at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tcell

import (
	"testing"
)

func mkTestScreen(t *testing.T, charset string) SimulationScreen {
	s := NewSimulationScreen(charset)
	if s == nil {
		t.Fatalf("Failed to get simulation screen")
	}
	if e := s.Init(); e != nil {
		t.Fatalf("Failed to initialize screen: %v", e)
	}
	return s
}

func TestInitScreen(t *testing.T) {

	s := mkTestScreen(t, "")
	defer s.Fini()

	if x, y := s.Size(); x != 80 || y != 25 {
		t.Fatalf("Size should be 80, 25, was %v, %v", x, y)
	}
	if s.CharacterSet() != "UTF-8" {
		t.Fatalf("Character Set (%v) not UTF-8", s.CharacterSet())
	}
	if b, x, y := s.GetContents(); len(b) != x*y || x != 80 || y != 25 {
		t.Fatalf("Contents (%v, %v, %v) wrong", len(b), x, y)
	}
}

func TestClearScreen(t *testing.T) {
	s := mkTestScreen(t, "")
	defer s.Fini()
	s.Clear()
	b, x, y := s.GetContents()
	if len(b) != x*y || x != 80 || y != 25 {
		t.Fatalf("Contents (%v, %v, %v) wrong", len(b), x, y)
	}
	for i := 0; i < x*y; i++ {
		if len(b[i].Runes) == 1 && b[i].Runes[0] != ' ' {
			t.Errorf("Incorrect contents at %v: %v", i, b[i].Runes)
		}
		if b[i].Style != StyleDefault {
			t.Errorf("Incorrect style at %v: %v", i, b[i].Style)
		}
	}
}

func TestSetCell(t *testing.T) {
	st := StyleDefault.Background(ColorRed).Blink(true)
	s := mkTestScreen(t, "")
	defer s.Fini()
	s.SetCell(2, 5, st, '@')
	b, _, _ := s.GetContents()
	s.Show()
	if len(b) != 80*25 {
		t.Fatalf("Wrong content size")
	}
	cell := &b[5*80+2]
	if len(cell.Runes) != 1 || len(cell.Bytes) != 1 ||
		cell.Runes[0] != '@' || cell.Bytes[0] != '@' ||
		cell.Style != st {
		t.Errorf("Incorrect cell content: %v", cell)
	}
}

func TestResize(t *testing.T) {
	st := StyleDefault.Background(ColorYellow).Underline(true)
	s := mkTestScreen(t, "")
	defer s.Fini()
	s.SetCell(2, 5, st, '&')
	b, x, y := s.GetContents()
	s.Show()

	cell := &b[5*80+2]
	if len(cell.Runes) != 1 || len(cell.Bytes) != 1 ||
		cell.Runes[0] != '&' || cell.Bytes[0] != '&' ||
		cell.Style != st {
		t.Errorf("Incorrect cell content: %v", cell)
	}
	s.SetSize(30, 10)
	s.Show()
	b2, x2, y2 := s.GetContents()
	if len(b2) == len(b) || x2 == x || y2 == y {
		t.Errorf("Screen parameters should not match")
	}

	cell2 := &b[5*80+2]
	if len(cell2.Runes) != 1 || len(cell2.Bytes) != 1 ||
		cell2.Runes[0] != '&' || cell2.Bytes[0] != '&' ||
		cell2.Style != st {
		t.Errorf("Incorrect cell content after resize: %v", cell2)
	}
}