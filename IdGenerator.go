package main

import (
        "time"
        "crypto/rand"
        "encoding/base64"
)

type IDRange struct {
        LowerBound uint64 `json:"lower_bound"`
        UpperBound uint64 `json:"upper_bound"`
        MachineId uint16 `json:"machine_id"`
}

type IDList struct {
        List []uint64 `json:"id_list"`
        MachineId uint16 `json:"machine_id"`
}

type StringIDList struct {
        List []string `json:"id_list"`
}

const KeyLengthInBytes = 32

func initSnowFlake(st *Settings ) *SnowFlake {
        if (st == nil) {
                st = &Settings{}
                // Default start-time is the Jan 01, 2014 and the ID generator should work for 174 yeard from then
                st.StartTime = time.Date(2014, 1, 1, 0, 0, 0, 0, time.UTC)
        }
        sf := NewSnowFlake(*st)
        if sf == nil {
                panic("snowFlake not created")
        }
        return sf

}

func GenerateIDRange(settings *Settings ) (*IDRange, error) {
        snowFlake := initSnowFlake(settings)
        lower, upper, err := snowFlake.NextIDRange()
        if err != nil {
                return nil, err
        }
        idRange := &IDRange{LowerBound:lower, UpperBound:upper, MachineId:snowFlake.machineID}
        return idRange, nil
}


func GenerateIDList(settings *Settings) (*IDList, error) {
        snowFlake := initSnowFlake(settings)
        ids, err := snowFlake.NextIDs()
        if err != nil {
                return nil, err
        }
        idList := &IDList{List:ids, MachineId:snowFlake.machineID}
        return idList, nil
}

// returns a set of random string ids
func GenerateRandomStringId(keyLength int, numIds int) []string {
        ids := make([]string, 0, numIds)
        i := 0
        for; i < numIds; {
                ids = append(ids, generateRandomString(keyLength))
                i++
        }
        return ids
}

// source https://elithrar.github.io/article/generating-secure-random-numbers-crypto-rand/

// GenerateRandomBytes returns securely generated random bytes.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomBytes(n int) ([]byte, error) {
        b := make([]byte, n)
        _, err := rand.Read(b)
        // Note that err == nil only if we read len(b) bytes.
        if err != nil {
                return nil, err
        }

        return b, nil
}

// GenerateRandomString returns a URL-safe, base64 encoded
// securely generated random string.
// It will return an error if the system's secure random
// number generator fails to function correctly, in which
// case the caller should not continue.
func generateRandomString(s int) string {
        b, err := generateRandomBytes(s)
        if (err != nil) {
                panic ("error while generating random bytes")
        }
        return base64.URLEncoding.EncodeToString(b)
}

