package sync15

import (
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/juruen/rmapi/archive"
	"github.com/juruen/rmapi/model"
)

func metadataJSON(t *testing.T, d *BlobDoc) map[string]interface{} {
	t.Helper()
	b, err := json.Marshal(d.Metadata)
	if err != nil {
		t.Fatal(err)
	}
	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		t.Fatal(err)
	}
	return m
}

func TestNewBlobDocSetsCreatedTime(t *testing.T) {
	before := time.Now().UnixNano() / 1000000
	m := metadataJSON(t, NewBlobDoc("doc", "id", model.DocumentType, ""))
	after := time.Now().UnixNano() / 1000000

	ct, ok := m["createdTime"].(string)
	if !ok || ct == "" {
		t.Fatalf("createdTime missing from uploaded metadata: %v", m)
	}
	v, err := strconv.ParseInt(ct, 10, 64)
	if err != nil || v < before || v > after {
		t.Fatalf("createdTime %q not a current millisecond timestamp (want %d..%d)", ct, before, after)
	}
	if m["createdTime"] != m["lastModified"] {
		t.Errorf("createdTime %v should equal lastModified %v on a new doc", m["createdTime"], m["lastModified"])
	}
}

func TestMetadataRoundTripPreservesCreatedTime(t *testing.T) {
	in := []byte(`{"visibleName":"x","type":"DocumentType","parent":"","lastModified":"2","createdTime":"1"}`)
	var meta archive.MetadataFile
	if err := json.Unmarshal(in, &meta); err != nil {
		t.Fatal(err)
	}
	out, _ := json.Marshal(meta)
	m := map[string]interface{}{}
	_ = json.Unmarshal(out, &m)
	if m["createdTime"] != "1" {
		t.Fatalf("createdTime lost on round-trip (mv would strip it): %s", out)
	}
}

func TestMetadataRoundTripOmitsAbsentCreatedTime(t *testing.T) {
	var meta archive.MetadataFile
	_ = json.Unmarshal([]byte(`{"visibleName":"x","lastModified":"2"}`), &meta)
	out, _ := json.Marshal(meta)
	m := map[string]interface{}{}
	_ = json.Unmarshal(out, &m)
	if _, present := m["createdTime"]; present {
		t.Fatalf("createdTime should stay absent when it was absent: %s", out)
	}
}
