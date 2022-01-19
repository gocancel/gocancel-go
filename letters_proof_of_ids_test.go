package gocancel

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestLettersService_DownloadProofOfID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("/api/v1/letters/a/proof_of_ids/b", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testHeader(t, r, "Accept", mediaTypeLetterProofOfID)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", "attachment; filename=proof.png")
		fmt.Fprint(w, "Hello World")
	})

	ctx := context.Background()
	reader, _, err := client.Letters.DownloadProofOfID(ctx, "a", "b")
	if err != nil {
		t.Fatalf("Letters.DownloadProofOfID returned error: %v", err)
	}

	content, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Errorf("Letters.DownloadProofOfID returned bad reader: %v", err)
	}

	want := []byte("Hello World")
	if !bytes.Equal(content, want) {
		t.Errorf("Letters.DownloadProofOfID returned %+v, want %+v", content, want)
	}

	const methodName = "DownloadProofOfID"
	testBadOptions(t, methodName, func() (err error) {
		_, _, err = client.Letters.DownloadProofOfID(ctx, "\n", "\n")
		return err
	})
}
