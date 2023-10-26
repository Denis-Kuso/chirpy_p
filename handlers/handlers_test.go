package handlers

import (
    "testing"
    "net/http"
    "net/http/httptest"
    "io/ioutil"
    "fmt"
    "strings"
)

type UserFeedback struct {
    valid bool `json:"valid"`
}
type Test struct {
    name string
    server *httptest.Server
    isValid bool
    expStatusCode int
    expBody string
}

func TestValidateChirp(t *testing.T) {
    tests := []Test{
        {"basic-request",
        httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
        r *http.Request) {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{ "body": "This is an opinion I need to share with the world"}`)) 
        })),
        true,
        200,
        `"valid": true`,
    },
    {"invalid body length",
    httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,
        r *http.Request) {
            w.WriteHeader(http.StatusOK)
            w.Write([]byte(`{ "body": "This is an opinion I need to share with the world, so badly that I am willing to write soooooo man chaaarsss to share it with everyone, like yeah, so that an actual error will pop out. At least it should should it not?"}`))
        })),
    false,
    400,
    `"valid":false`,
    }}
    for _, test := range tests {
		t.Run(test.name, func(t *testing.T){
			defer test.server.Close()
            //implement
            req, _ := http.Get(test.server.URL)
            body, _ := ioutil.ReadAll(req.Body)
            if test.expStatusCode != req.StatusCode {
                t.Errorf("FAILED with %s. Wrong status code! Expected: %d, got: %d",test.name, test.expStatusCode, req.StatusCode)
            }
            fmt.Println(body)
        })
    }
// test when request data is ok
// test when request data has text of invalid text length
// test error on our side
}

func TestHandlerCheckStatus(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(CheckStatus))
    resp, err := http.Get(server.URL)
    if err != nil {
        t.Fatalf("err: %v",err)
    }
    // status
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Err: expected %d, recevied: %d", resp.StatusCode, http.StatusOK)
    }
    defer resp.Body.Close()

    // body
    expectedText := "OK"
    b, err := ioutil.ReadAll(resp.Body)
    if expectedText != string(b) {
        t.Fatalf("err: expected: %s, got: %s", expectedText, string(b))
    }
}

func TestResetViews(t *testing.T) {
    s := ApiState {
        ViewCount:10, // to test reseting
    }
    server := httptest.NewServer(http.HandlerFunc(s.ResetViews))
    resp, err := http.Get(server.URL)
    if err != nil {
        t.Fatalf("err: %v",err)
    }
    // status
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Err: expected %d, recevied: %d", resp.StatusCode, http.StatusOK)
    }
    defer resp.Body.Close()

    expectedCount := 0
    if s.ViewCount != expectedCount {
        t.Errorf("Err: expected %d, got: %d", expectedCount, s.ViewCount)
    }

    expectedReply := "Hits reset to 0"
    b, err := ioutil.ReadAll(resp.Body)
    if expectedReply != string(b) {
        t.Fatalf("err: expected: %s, got: %s", expectedReply, string(b))
    }
    // test reseting
}

func TestShowPageViews(t *testing.T) {
    s := ApiState {
        ViewCount:0,
    }
    server := httptest.NewServer(http.HandlerFunc(s.ShowPageViews))
    resp, err := http.Get(server.URL)
    if err != nil {
        t.Fatalf("err: %v",err)
    }
    // status
    if resp.StatusCode != http.StatusOK {
        t.Errorf("Err: expected %d, recevied: %d", resp.StatusCode, http.StatusOK)
    }
    defer resp.Body.Close()

    expectedBody := "Welcome, Chirpy Admin"

    b, err := ioutil.ReadAll(resp.Body)
    if !strings.Contains(string(b), expectedBody) {
        t.Fatalf("err: expected body to contain: %s, got: %s", expectedBody, string(b))
    }
}

