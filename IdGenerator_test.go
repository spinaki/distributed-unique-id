package main

import (
        "testing"
        "fmt"
        "encoding/json"
        "time"
        "github.com/stretchr/testify/assert"
        "strings"
)

func TestIDRange(t *testing.T) {
        idRange, err := GenerateIDRange(nil)
        s, err := json.Marshal(idRange)
        fmt.Println(string(s))
        if (err != nil) {
                t.Fatal("id range cannot be marshalled")
        }
        assert.Condition(t, func() bool { return idRange.LowerBound > uint64(150116117052998656)},
                "LowerBound should be greater than previous lowerbound")
        assert.Condition(t, func() bool { return idRange.LowerBound > uint64(150116117052998911)},
                "LowerBound should be greater than previous upperbound")
        assert.Condition(t, func() bool { return idRange.UpperBound > uint64(150116117052998911)},
                "UpperBound should be greater than previous upperbound")
        assert.Equal(t, uint64(255), (idRange.UpperBound - idRange.LowerBound), "Upper and Lower Bound Difference Mismatch")

        idRange, err = GenerateIDRange(&Settings{StartTime:time.Now()})
        s, err = json.Marshal(idRange)
        if (err != nil) {
                t.Fatal("id range cannot be marshalled")
        }
        //assert.Equal(t, uint64(82944), idRange.LowerBound, "LowerBound mismatch")
        //assert.Equal(t, uint64(83199), idRange.UpperBound, "UpperBound mismatch")
        assert.Equal(t, uint64(255), (idRange.UpperBound - idRange.LowerBound), "Upper and Lower Bound Difference Mismatch")
}

func TestIDList(t *testing.T) {
        idList, err := GenerateIDList(nil)
        if (err != nil) {
                t.Fatal("idList not generated")
        }
        s, err := json.Marshal(idList)
        if (err != nil) {
                t.Fatal("id list cannot be marshalled")
        }
        assert.Equal(t, 256, len(idList.List), "Length of ID List should be 256")
        assert.Equal(t, 256, cap(idList.List), "Capacity of ID List should be 256")
        assert.Condition(t, func() bool {
                val := idList.List[0]
                return val > uint64(150549601139639296)},
                "LowerBound should be greater than previous lowerbound")
        assert.Condition(t, func() bool {
                val := idList.List[0]
                return val > uint64(150549601139639551)},
                "LowerBound should be greater than previous upperbound")
        assert.Condition(t, func() bool {
                val :=  idList.List[255]
                return val > uint64(150549601139639551)},
                "Upperbound should be greater than previous upperbound")
        assert.Condition(t, func() bool {
                val := idList.List[0]
                return val > uint64(150549601139639296)},
                "Upperbound should be greater than previous lowerbound")
        lastVal := idList.List[255]
        firstVal := idList.List[0]
        assert.Equal(t, uint64(255), (lastVal - firstVal), "Upper and Lower Bound Difference Mismatch")
        fmt.Println(string(s))
}

func TestStringId(t *testing.T) {
        ids := GenerateRandomStringId(32, 2)
        assert.Equal(t, 44, len(ids[0]), "ID length should be 44")
        assert.Condition(t, func() bool {return strings.HasSuffix(ids[0], "=")}, "ID  should end with =")
        strIdList := &StringIDList{List:ids}
        s, err := json.Marshal(strIdList)
        if (err != nil) {
                t.Fatal("id list cannot be marshalled")
        }
        fmt.Println(string(s))
}

func TestSingleID(t *testing.T) {
        idRange, err := GenerateIDRange(nil)
        s, err := json.Marshal(idRange)
        if err != nil {
                fmt.Println(err)
        }
        fmt.Println(string(s))
}