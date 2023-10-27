package database

import (
    "testing"
    "os"
    "log"
)

// Strategy
// no data
// one item
// several items
// match length and content
func TestGetChirps(t *testing.T) {
}


// file exists, return nil (should be database.json)
// file does not exist, should then exist
func TestEnsureDB(t *testing.T) {
    type Test struct {
        Name string
        Expected error
    }
//    db, err = NewDB("./database.json")
//    if err != nil {
//        log.Fatalf("SETUP FAILED due to %v\n",err)
//    }
    testFile := "./test_database.json"
    db := DB {
        path: testFile,
    }
    tests := []Test{
        {"file does not exist",
        nil,},
        {"file does exist",
        nil,},
    }
    for _, test := range tests {
        err := db.ensureDB()
        if err != test.Expected {
            t.Errorf("FAILED on case: %s\n, expected: %v, got:%v\n", test.Name, test.Expected, err)
        }
    }
    err := os.Remove(testFile)
    if err != nil {
        log.Fatalf("Error: %v during removal of %s\n", err, testFile)
    }
}


func TestLoadDB(t *testing.T) {
}
