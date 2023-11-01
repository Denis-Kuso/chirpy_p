package database

import (
    "testing"
    "fmt"
    "os"
    "log"
    "errors"
    "reflect"
)

/* strategy:
empty string
non-empty string
*/ 
func TestCreateChirp(t *testing.T) {
    subtests := []struct{
        name string
        InputBody string
        expectedData Chirp
        expectedErr error
    }{
        {"non-empty string",
        "hello, world",
        Chirp{
            Id:1,
            Body: "hello, world",
        },
        nil},
        {"empty string",
        "",
        Chirp{
            Id: 2,
            Body: ""},
        nil,},
    }//TODO add testcases

    // use actual file
    db,err := NewDB("./test_database.json")
    if err != nil{
        t.Fatalf("Error loading database: %v\n",err)
    }
    err = db.ensureDB()
    if err != nil {
        t.Fatalf("err: %v\n",err)
    }

    for _,subtest := range subtests {
        t.Run(subtest.name, func(t *testing.T) {
            result, err := db.CreateChirp(subtest.InputBody)
            fmt.Println(result)
            fmt.Println(subtest.expectedData)
            if !errors.Is(err, subtest.expectedErr){
                t.Errorf("expected error (%v), got error (%v)", subtest.expectedErr, err)
            }
            if !reflect.DeepEqual(result, subtest.expectedData) {
                t.Errorf("expected %v, got %v\n",subtest.expectedData, result)
            }
        })
    }
    // mock writting to file
}


/* stragegy:
nil dbStructure
non-nil dbStructure
*/
func TestWrite(t *testing.T) {
}


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

// successful load
 // no items
 // nonzero items
// unsuccessful load
 // unmarshalling
 // file does not exist
func TestLoadDB(t *testing.T) {
}
