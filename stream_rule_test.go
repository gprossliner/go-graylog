package graylog_test

import (
	"testing"

	"github.com/suzuki-shunsuke/go-graylog/test"
)

func TestGetStreamRules(t *testing.T) {
	test.TestGetStreamRules(t)
}

func TestCreateStreamRule(t *testing.T) {
	test.TestCreateStreamRule(t)
}

// func TestUpdateStreamRule(t *testing.T) {
// 	server, client, err := testutil.GetServerAndClient()
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	defer server.Close()
// 	indexSet := dummyIndexSet()
// 	is, _, err := client.CreateIndexSet(indexSet)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	stream := dummyStream()
// 	stream.IndexSetID = is.ID
// 	is, _, err = client.CreateStream(stream)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
//
// 	server.streams[stream.ID] = *stream
// 	stream.Description = "changed!"
// 	updatedStream, _, err := client.UpdateStream(stream.ID, stream)
// 	if err != nil {
// 		t.Fatal("Failed to UpdateStream", err)
// 	}
// 	if updatedStream == nil {
// 		t.Fatal("UpdateStream() == nil, nil")
// 	}
// 	if updatedStream.Title != stream.Title {
// 		t.Fatalf(
// 			"updatedStream.Title == %s, wanted %s",
// 			updatedStream.Title, stream.Title)
// 	}
// 	if _, _, err := client.UpdateStream("", stream); err == nil {
// 		t.Fatal("id is required")
// 	}
// 	if _, _, err := client.UpdateStream("h", stream); err == nil {
// 		t.Fatal(`no stream whose id is "h"`)
// 	}
// }
