package cuckoo

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
	"testing"
)

func TestInsertion(t *testing.T) {
	cf := NewFilter(1000000)
	fd, err := os.Open("/usr/share/dict/words")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	scanner := bufio.NewScanner(fd)

	var values [][]byte
	var lineCount uint
	for scanner.Scan() {
		s := []byte(scanner.Text())
		if cf.InsertUnique(s) {
			lineCount++
		}
		values = append(values, s)
	}

	count := cf.Count()
	if count != lineCount {
		t.Errorf("Expected count = %d, instead count = %d", lineCount, count)
	}

	for _, v := range values {
		cf.Delete(v)
	}

	count = cf.Count()
	if count != 0 {
		t.Errorf("Expected count = 0, instead count == %d", count)
	}
}

func TestEncodeDecode(t *testing.T) {
	cf := NewFilter(8)
	cf.buckets = []bucket{
		[4]fingerprint{1, 2, 3, 4},
		[4]fingerprint{5, 6, 7, 8},
	}
	cf.count = 8
	encoded := cf.Encode()
	got, err := Decode(encoded)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(cf, got) {
		t.Errorf("Decode = %v, want %v, encoded = %v", got, cf, encoded)
	}
}
