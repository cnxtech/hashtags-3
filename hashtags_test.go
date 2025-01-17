/*
   Copyright © 2019 M.Watermann, 10247 Berlin, Germany
                  All rights reserved
              EMail : <support@mwat.de>
*/

package hashtags

import (
	"log"
	"os"
	"reflect"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

func delDB(aFilename string) string {
	os.Remove(aFilename)

	return aFilename
} // delDB()

func Test_tSourceList_indexOf(t *testing.T) {
	sl1 := &tSourceList{
		"one",
		"two",
		"three",
		"four",
		"five",
	}
	type args struct {
		aID string
	}
	tests := []struct {
		name string
		sl   *tSourceList
		args args
		want int
	}{
		// TODO: Add test cases.
		{" 1", sl1, args{"one"}, 0},
		{" 2", sl1, args{"two"}, 1},
		{" 3", sl1, args{"three"}, 2},
		{" 4", sl1, args{"four"}, 3},
		{" 5", sl1, args{"five"}, 4},
		{" 6", sl1, args{"six"}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.indexOf(tt.args.aID); got != tt.want {
				t.Errorf("tSourceList.indexOf() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_tSourceList_indexOf()

func Test_tSourceList_removeID(t *testing.T) {
	sl1 := &tSourceList{
		"one",
		"two",
		"three",
		"four",
		"five",
	}
	wl1 := &tSourceList{
		"two",
		"three",
		"four",
		"five",
	}
	wl2 := &tSourceList{
		"two",
		"three",
		"four",
	}
	wl3 := &tSourceList{
		"two",
		"four",
	}
	type args struct {
		aID string
	}
	tests := []struct {
		name string
		sl   *tSourceList
		args args
		want *tSourceList
	}{
		// TODO: Add test cases.
		{" 1", sl1, args{"one"}, wl1},
		{" 2", sl1, args{"five"}, wl2},
		{" 3", sl1, args{"three"}, wl3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.removeID(tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tSourceList.IDremove() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_tSourceList_removeID()

func Test_tSourceList_renameID(t *testing.T) {
	sl1 := &tSourceList{
		"one",
		"two",
		"three",
	}
	wl1 := &tSourceList{
		"four",
		"one",
		"two",
	}
	wl2 := &tSourceList{
		"one",
		"six",
		"two",
	}
	type args struct {
		aOldID string
		aNewID string
	}
	tests := []struct {
		name string
		sl   *tSourceList
		args args
		want *tSourceList
	}{
		// TODO: Add test cases.
		{" 1", sl1, args{"three", "four"}, wl1},
		{" 2", sl1, args{"four", "six"}, wl2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.renameID(tt.args.aOldID, tt.args.aNewID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tSourceList.renameID() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_tSourceList_renameID()

func Test_tSourceList_String(t *testing.T) {
	sl1 := &tSourceList{
		"one",
		"two",
		"three",
	}
	wl1 := "one\nthree\ntwo"
	sl2 := &tSourceList{}
	wl2 := ""
	tests := []struct {
		name string
		sl   *tSourceList
		want string
	}{
		// TODO: Add test cases.
		{" 1", sl1, wl1},
		{" 2", sl2, wl2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sl.String(); got != tt.want {
				t.Errorf("TSourceList.String() = %v, want %v", got, tt.want)
			}
		})
	}
} // Test_tSourceList_String()

func TestNew(t *testing.T) {
	fn := delDB("hashlist.db")
	fn2 := delDB("does.not.exist")
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1, _ := New(fn)
	hl1.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	hl1.Store()
	hl2, _ := New(fn2)
	type args struct {
		aFilename string
	}
	tests := []struct {
		name    string
		args    args
		want    *THashList
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", args{fn}, hl1, false},
		{" 2", args{fn2}, hl2, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.aFilename)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestNew()

func TestTHashList_Checksum(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id3, id1},
		},
		mtx: new(sync.RWMutex),
	}
	h1a, _ := New(fn)
	h1a.HashAdd(hash1, id2).
		HashAdd(hash2, id1).
		HashAdd(hash2, id3)
	w1 := h1a.Checksum()
	hl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id2},
			hash2: &tSourceList{id2, id3},
		},
		mtx: new(sync.RWMutex),
	}
	h2a, _ := New(fn)
	h2a.HashAdd(hash1, id1).
		HashAdd(hash1, id2).
		HashAdd(hash2, id2).
		HashAdd(hash2, id3).
		HashAdd(hash1, id2).
		HashAdd(hash2, id3)
	w2 := h2a.Checksum()
	tests := []struct {
		name     string
		hl       *THashList
		wantRSum uint32
	}{
		// TODO: Add test cases.
		{" 1", hl1, w1},
		{" 2", hl2, w2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			atomic.StoreUint32(&tt.hl.µChange, 0)
			if gotRSum := tt.hl.Checksum(); gotRSum != tt.wantRSum {
				t.Errorf("THashList.Checksum() = %v, want %v", gotRSum, tt.wantRSum)
			}
		})
	}
} // TestTHashList_Checksum()

func TestTHashList_Clear(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1, _ := New(fn)
	hl1.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	tests := []struct {
		name string
		hl   *THashList
		want int
	}{
		// TODO: Add test cases.
		{" 1", hl1, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.Clear().Len(); got != tt.want {
				t.Errorf("THashList.Clear() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_Clear()

func TestTHashList_CountedList(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "@mention1", "#another3"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := []TCountItem{
		{1, hash1},
		{2, hash2},
	}
	hl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id3},
			hash3: &tSourceList{id1, id2, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl2 := []TCountItem{
		{3, hash3},
		{1, hash1},
		{2, hash2},
	}
	tests := []struct {
		name      string
		hl        *THashList
		wantRList []TCountItem
	}{
		// TODO: Add test cases.
		{" 1", hl1, wl1},
		{" 2", hl2, wl2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRList := tt.hl.CountedList(); !reflect.DeepEqual(gotRList, tt.wantRList) {
				t.Errorf("THashList.CountedList() = %v, want %v", gotRList, tt.wantRList)
			}
		})
	}
} // TestTHashList_CountedList()

func TestTHashList_HashAdd(t *testing.T) {
	hash1, hash2 := "#hash2", "#hash1"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id1},
			hash1: &tSourceList{id2, id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id2, id1},
			hash1: &tSourceList{id2, id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl3 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id2, id1},
			hash1: &tSourceList{id2, id3, id1},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aHash string
		aID   string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1, id1}, wl1},
		{" 2", wl1, args{hash2, id2}, wl2},
		{" 3", wl2, args{hash1, id3}, wl3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.HashAdd(tt.args.aHash, tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.HashAdd() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_HashAdd()

func TestTHashList_HashLen(t *testing.T) {
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id3, id1},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aHash string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want int
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1}, 1},
		{" 2", hl1, args{hash2}, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.HashLen(tt.args.aHash); got != tt.want {
				t.Errorf("THashList.HashLen() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_HashLen()

func TestTHashList_HashList(t *testing.T) {
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2, id1},
			hash2: &tSourceList{id2, id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := []string{
		id2,
		id1,
	}
	type args struct {
		aHash string
	}
	tests := []struct {
		name      string
		hl        *THashList
		args      args
		wantRList []string
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1}, wl1},
		{" 2", hl1, args{hash2}, wl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRList := tt.hl.HashList(tt.args.aHash); !reflect.DeepEqual(gotRList, tt.wantRList) {
				t.Errorf("THashList.HashList() = %v, want %v", gotRList, tt.wantRList)
			}
		})
	}
} // TestTHashList_HashList()

func TestTHashList_HashRemove(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id2},
			hash2: &tSourceList{id1, id2},
		},
		fn:  fn,
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id2},
			hash2: &tSourceList{id1, id2},
		},
		fn:  fn,
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id1, id2},
		},
		fn:  fn,
		mtx: new(sync.RWMutex),
	}
	wl3 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id2},
		},
		fn:  fn,
		mtx: new(sync.RWMutex),
	}
	wl4, _ := New(fn)
	type args struct {
		aHash string
		aID   string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{hash1, id1}, wl1},
		{" 2", hl1, args{hash1, id2}, wl2},
		{" 3", hl1, args{hash2, id1}, wl3},
		{" 4", hl1, args{hash2, id2}, wl4},
		{" 5", hl1, args{hash1, id1}, wl4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.HashRemove(tt.args.aHash, tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.HashRemove() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_HashRemove()

func TestTHashList_IDlist(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id2},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	var wl0 []string
	wl1 := []string{hash1, hash3}
	wl2 := []string{hash1, hash2}
	wl3 := []string{hash2, hash3}
	type args struct {
		aID string
	}
	tests := []struct {
		name      string
		hl        *THashList
		args      args
		wantRList []string
	}{
		// TODO: Add test cases.
		{" 0", hl1, args{"@does.not.exist"}, wl0},
		{" 1", hl1, args{id1}, wl1},
		{" 2", hl1, args{id2}, wl2},
		{" 3", hl1, args{id3}, wl3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRList := tt.hl.IDlist(tt.args.aID); !reflect.DeepEqual(gotRList, tt.wantRList) {
				t.Errorf("THashList.IDlist() = %v, want %v", gotRList, tt.wantRList)
			}
		})
	}
} // TestTHashList_IDlist()

func TestTHashList_IDremove(t *testing.T) {
	// fn := delDB("hashlist.db")
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl3, _ := New("")
	type args struct {
		aID string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{id1}, wl1},
		{" 2", wl1, args{id2}, wl2},
		{" 3", wl2, args{id3}, wl3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.IDremove(tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.IDremove() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_IDremove()

func TestTHashList_IDrename(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3, id4, id5, id6 := "id_e", "id_a", "id_c", "id_g", "id_j", "id_k"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3, id1},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id3, id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3, id4},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id3, id4},
		},
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3, id4},
			hash2: &tSourceList{id3, id5},
			hash3: &tSourceList{id3, id4},
		},
		mtx: new(sync.RWMutex),
	}
	wl3 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id4, id6},
			hash2: &tSourceList{id5, id6},
			hash3: &tSourceList{id4, id6},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aOldID string
		aNewID string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{id1, id4}, wl1},
		{" 2", wl1, args{id2, id5}, wl2},
		{" 3", wl2, args{id3, id6}, wl3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.IDrename(tt.args.aOldID, tt.args.aNewID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.IDrename() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_IDrename()

func TestTHashList_IDupdate(t *testing.T) {
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id2, id3},
			hash2: &tSourceList{id1, id2},
		},
		mtx: new(sync.RWMutex),
	}
	tx1 := []byte("blabla " + hash1 + " blabla " + hash3 + " blabla")
	wl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id2, id3},
			hash2: &tSourceList{id2},
			hash3: &tSourceList{id1},
		},
		mtx: new(sync.RWMutex),
	}
	tx2 := []byte("blabla blabla blabla")
	wl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id3},
			hash3: &tSourceList{id1},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aID   string
		aText []byte
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{id1, tx1}, wl1},
		{" 2", hl1, args{id2, tx2}, wl2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.IDupdate(tt.args.aID, tt.args.aText); got.String() != tt.want.String() {
				t.Errorf("THashList.Update() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_IDupdate()

func TestTHashList_Len(t *testing.T) {
	fn := delDB("hashlist.db")
	hl1, _ := New(fn)
	hl2, _ := New(fn)
	hl2.HashAdd("#hash", "source")
	hl3, _ := New(fn)
	hl3.HashAdd("#hash2", "source1")
	hl4, _ := New(fn)
	hl4.HashAdd("#hash2", "source1").
		HashAdd("#hash3", "source2")
	tests := []struct {
		name string
		hl   *THashList
		want int
	}{
		// TODO: Add test cases.
		{" 1", hl1, 0},
		{" 2", hl2, 1},
		{" 3", hl3, 1},
		{" 4", hl4, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.Len(); got != tt.want {
				t.Errorf("THashList.Len() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_Len()

func TestTHashList_LenTotal(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := "id_c", "id_a", "id_b"
	hl1, _ := New(fn)
	hl1.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	hl2, _ := New(fn)
	hl2.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2).
		HashAdd(hash3, id1).
		HashAdd(hash3, id2).
		HashAdd(hash3, id3)
	tests := []struct {
		name       string
		hl         *THashList
		wantRCount int
	}{
		// TODO: Add test cases.
		{" 1", hl1, 6},
		{" 2", hl2, 10},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRCount := tt.hl.LenTotal(); gotRCount != tt.wantRCount {
				t.Errorf("THashList.LenTotal() = %v, want %v", gotRCount, tt.wantRCount)
			}
		})
	}
} // TestTHashList_LenTotal()

func TestTHashList_Load(t *testing.T) {
	// hl1, _ := New("")
	// hl1.SetFilename("load.txt")
	// hl1.xLoad()
	// hl1.SetFilename("load.db")
	// hl1.store()
	fn := delDB("hashlist.db")
	fn2 := delDB(".does.not.exist")
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1, _ := New(fn)
	hl1.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	hl1.store()
	hl1.Clear()
	hl2, _ := New(fn2)
	hl3, _ := New("")
	tests := []struct {
		name    string
		hl      *THashList
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", hl1, 2, false},
		{" 2", hl2, 0, false},
		{" 3", hl3, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.hl.Load()
			if (err != nil) != tt.wantErr {
				t.Errorf("THashList.Load() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Len() != tt.want {
				t.Errorf("THashList.Load() = %v, want %v", got.Len(), tt.want)
			}
		})
	}
} // TestTHashList_Load()

func TestTHashList_parseID(t *testing.T) {
	// fn := delDB("hashlist.db")
	hash1, hash2, hash3, hash4 := "#HÄSCH1", "#hash2", "#hash3", "#hash4"
	lh1 := strings.ToLower(hash1)
	id1, id2, id3, id4, id5 := "id_c", "id_a", "id_b", "id_d", "id_e"
	hl1, _ := New("")
	tx1 := []byte("blabla " + hash1 + " blabla " + hash3 + ". Blabla")
	wl1 := &THashList{
		hl: tHashMap{
			lh1:   &tSourceList{id1},
			hash3: &tSourceList{id1},
		},
		mtx: new(sync.RWMutex),
	}
	hl2 := &THashList{
		hl: tHashMap{
			lh1:   &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	tx2 := []byte("blabla " + hash2 + ". Blabla " + hash3 + " blabla")
	wl2 := &THashList{
		hl: tHashMap{
			lh1:   &tSourceList{id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id2, id3},
		},
		mtx: new(sync.RWMutex),
	}
	hl3, _ := New("")
	tx3 := []byte("\n> #KurzErklärt #Zensurheberrecht verhindern – \n> [Glyphosat-Gutachten selbst anfragen!](https://fragdenstaat.de/aktionen/zensurheberrecht-2019/)\n")
	wl3 := &THashList{
		hl: tHashMap{
			"#kurzerklärt":      &tSourceList{id3},
			"#zensurheberrecht": &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	hl4, _ := New("")
	tx4 := []byte("blabla **" + hash1 + "** blabla\n\n_" + hash3 + "_")
	wl4 := &THashList{
		hl: tHashMap{
			lh1:   &tSourceList{id4},
			hash3: &tSourceList{id4},
		},
		mtx: new(sync.RWMutex),
	}
	hl5, _ := New("")
	tx5 := []byte(`blabla&#39; **` + hash2 + `** blabla\n\<a href="page#` + hash3 + `">txt</a> ` + hash4)
	wl5 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id5},
			hash4: &tSourceList{id5},
		},
		mtx: new(sync.RWMutex),
	}
	type args struct {
		aID   string
		aText []byte
	}
	tests := []struct {
		name   string
		fields *THashList //fields
		args   args
		want   *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{id1, tx1}, wl1},
		{" 2", hl2, args{id2, tx2}, wl2},
		{" 3", hl3, args{id3, tx3}, wl3},
		{" 4", hl4, args{id4, tx4}, wl4},
		{" 5", hl5, args{id5, tx5}, wl5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hl := &THashList{
				fn:      tt.fields.fn,
				hl:      tt.fields.hl,
				mtx:     tt.fields.mtx,
				µChange: tt.fields.µChange,
				µCC:     tt.fields.µCC,
			}
			if got := hl.parseID(tt.args.aID, tt.args.aText); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.parseID() = '%v',\nwant '%v'", got, tt.want)
			}
		})
	}
} // TestTHashList_parseID()

func TestTHashList_remove(t *testing.T) {
	// fn := delDB("hashlist.db")
	hash1, hash2, hash3 := "#hash1", "#hash2", "#hash3"
	id1, id2, id3 := "id_3", "id_1", "id_2"
	hl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id1, id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl1 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id2, id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl2 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id1, id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl3 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
			hash3: &tSourceList{id1},
		},
		mtx: new(sync.RWMutex),
	}
	wl4 := &THashList{
		hl: tHashMap{
			hash1: &tSourceList{id3},
			hash2: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl5 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl6 := &THashList{
		hl: tHashMap{
			hash2: &tSourceList{id3},
		},
		mtx: new(sync.RWMutex),
	}
	wl7, _ := New("")
	type args struct {
		aDelim  byte
		aMapIdx string
		aID     string
	}
	tests := []struct {
		name string
		hl   *THashList
		args args
		want *THashList
	}{
		// TODO: Add test cases.
		{" 1", hl1, args{'#', hash1, id1}, wl1},
		{" 2", wl1, args{'#', hash2, id2}, wl2},
		{" 3", wl2, args{'#', hash3, id3}, wl3},
		{" 4", wl3, args{'#', hash3, id1}, wl4},
		{" 5", wl4, args{'#', hash1, id3}, wl5},
		{" 6", wl5, args{'#', hash1, id3}, wl6},
		{" 7", wl6, args{'#', hash2, id3}, wl7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.remove(tt.args.aDelim, tt.args.aMapIdx, tt.args.aID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("THashList.remove() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_remove()

func TestTHashList_store(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2 := "#hash1", "#Zensurheberrecht"
	id1, id2 := "id_c", "id_a"
	hl1, _ := New(fn)
	hl1.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	hl2, _ := New("")
	hl2.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	tests := []struct {
		name    string
		hl      *THashList
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{" 1", hl1, 91, false},
		{" 2", hl2, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.hl.store()
			if (err != nil) != tt.wantErr {
				t.Errorf("THashList.Store() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("THashList.Store() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_Store()

func TestTHashList_String(t *testing.T) {
	fn := delDB("hashlist.db")
	hash1, hash2 := "#hash1", "#hash2"
	id1, id2 := "id_c", "id_a"
	hl1, _ := New(fn)
	hl1.HashAdd(hash1, id1).
		HashAdd(hash2, id2).
		HashAdd(hash2, id1).
		HashAdd(hash1, id2)
	wl1 := "[" + hash1 + "]\n" + id2 + "\n" + id1 +
		"\n[" + hash2 + "]\n" + id2 + "\n" + id1 + "\n"
	tests := []struct {
		name string
		hl   *THashList
		want string
	}{
		// TODO: Add test cases.
		{" 1", hl1, wl1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.hl.String(); got != tt.want {
				t.Errorf("THashList.String() = %v, want %v", got, tt.want)
			}
		})
	}
} // TestTHashList_String()

func Benchmark_LoadTxT(b *testing.B) {
	hl, _ := New("")
	hl.SetFilename("load.txt")
	UseBinaryStorage = false
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := hl.Load(); nil != err {
			log.Printf("LoadTxt(): %v", err)
		}
	}
} // Benchmark_LoadTxt()

func Benchmark_LoadBin(b *testing.B) {
	hl, _ := New("")
	hl.SetFilename("load.db")
	UseBinaryStorage = true
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		if _, err := hl.Load(); nil != err {
			log.Printf("LoadBin(): %v", err)
		}
	}
} // Benchmark_LoadBin()

func Benchmark_StoreTxt(b *testing.B) {
	UseBinaryStorage = false
	hl, _ := New("load.txt")
	hl.Load()

	for n := 0; n < b.N; n++ {
		if _, err := hl.Store(); nil != err {
			log.Printf("StoreTxt(): %v", err)
		}
	}
} // Benchmark_StoreTxt()

func Benchmark_StoreBin(b *testing.B) {
	UseBinaryStorage = true
	hl, _ := New("load.db")
	hl.Load()

	for n := 0; n < b.N; n++ {
		if _, err := hl.Store(); nil != err {
			log.Printf("StoreBin(): %v", err)
		}
	}
} // Benchmark_StoreBin()
