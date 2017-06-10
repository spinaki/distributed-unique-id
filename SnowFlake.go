package uniqueidgenerator
import (
        "errors"
        "net"
        "sync"
        "time"
        "io/ioutil"
        "net/http"
        "os"
        "fmt"
)
// These constants are the bit lengths of SnowFlake ID parts.
const (
        BitLenTime      = 39                               // bit length of time
        BitLenMachineID = 63 - BitLenTime - BitLenSequence // bit length of machine id (16 bits)
        BitLenSequence  = 8                                // bit length of sequence number
)

// Settings configures SnowFlake:
//
// StartTime is the time since which the SnowFlake time is defined as the elapsed time.
// If StartTime is 0, the start time of the SnowFlake is set to "2014-09-01 00:00:00 +0000 UTC".
// If StartTime is ahead of the current time, SnowFlake is not created.
//
// MachineID returns the unique ID of the SnowFlake instance.
// If MachineID returns an error, SnowFlake is not created.
// If MachineID is nil, default MachineID is used.
// Default MachineID returns the lower 16 bits of the private IP address.
//
// CheckMachineID validates the uniqueness of the machine ID.
// If CheckMachineID returns false, SnowFlake is not created.
// If CheckMachineID is nil, no validation is done.
type Settings struct {
        StartTime      time.Time
        MachineID      func() (uint16, error)
        CheckMachineID func(uint16) bool
}

// SnowFlake is a distributed unique ID generator.
type SnowFlake struct {
        mutex       *sync.Mutex
        startTime   int64
        elapsedTime int64
        sequence    uint16
        machineID   uint16
}

// NewSnowFlake returns a new SnowFlake configured with the given Settings.
// NewSnowFlake returns nil in the following cases:
// - Settings.StartTime is ahead of the current time.
// - Settings.MachineID returns an error.
// - Settings.CheckMachineID returns false.
func NewSnowFlake(st Settings) *SnowFlake {
        sf := new(SnowFlake)
        sf.mutex = new(sync.Mutex)
        // why is it set to max value ?
        sf.sequence = uint16( 1 << BitLenSequence - 1)
        if st.StartTime.After(time.Now()) {
                return nil
        }
        if st.StartTime.IsZero() {
                sf.startTime = toSnowFlakeTime(time.Date(2014, 9, 1, 0, 0, 0, 0, time.UTC))
        } else {
                sf.startTime = toSnowFlakeTime(st.StartTime)
        }

        var err error
        if st.MachineID == nil {
                sf.machineID, err = lower16BitPrivateIP()
        } else {
                sf.machineID, err = st.MachineID()
        }
        if err != nil || (st.CheckMachineID != nil && !st.CheckMachineID(sf.machineID)) {
                return nil
        }

        return sf
}

// elapsedTime, machine-id and sequence
func (sf *SnowFlake) NextIDs() ([]uint64, error) {
        sf.mutex.Lock()
        defer sf.mutex.Unlock()
        sf.elapsedTime = currentElapsedTime(sf.startTime)
        sf.sequence = 0
        const maxSequence = uint16(1 << BitLenSequence - 1)
        idList := make([]uint64, 0, maxSequence + 1)
        for; sf.sequence <= maxSequence; {
                id, err := sf.toID()
                if err != nil {
                        return nil, err
                }
                idList = append(idList, id)
                sf.sequence = (sf.sequence + 1)
        }
        return idList, nil
}

// returns the lowerMost and upperMost values in the unique ID range
func (sf *SnowFlake) NextIDRange () (uint64, uint64, error) {
        sf.mutex.Lock()
        defer sf.mutex.Unlock()
        sf.elapsedTime = currentElapsedTime(sf.startTime)
        sf.sequence = 0
        lower, err := sf.toID()
        if (err != nil) {
                return 0, 0, err
        }
        sf.sequence = uint16(1 << BitLenSequence - 1)
        upper, err := sf.toID()
        if (err != nil) {
                return 0, 0, err
        }
        return lower, upper, nil
}

// NextID generates a next unique ID.
// After the SnowFlake time overflows, NextID returns an error.
// ONLY USED in Testing ??
func (sf *SnowFlake) NextID() (uint64, error) {
        const maskSequence = uint16(1 << BitLenSequence - 1)

        sf.mutex.Lock()
        defer sf.mutex.Unlock()
        current := currentElapsedTime(sf.startTime)
        //fmt.Println(sf.elapsedTime, current)
        if sf.elapsedTime < current {
                // this is only executed the first time
                // this will be executed if the elapsedTime is not set correctly to current time
                fmt.Println("elapsedTime less than current")
                sf.elapsedTime = current
                sf.sequence = 0
        } else { // sf.elapsedTime >= current
                sf.sequence = (sf.sequence + 1) & maskSequence
                if sf.sequence == 0 {
                        sf.elapsedTime++
                        overtime := sf.elapsedTime - current
                        fmt.Println("SLEEP FOR DURATION: ", sleepTime(overtime))
                        time.Sleep(sleepTime((overtime)))
                }
        }
        return sf.toID()
}

func (sf *SnowFlake) NextIDRange1 () (uint64, uint64, error) {
        sf.mutex.Lock()
        defer sf.mutex.Unlock()
        const maskSequence = uint16(1 << BitLenSequence - 1)
        current := currentElapsedTime(sf.startTime)
        fmt.Println(sf.elapsedTime, current)
        if sf.elapsedTime < current {
                // this is only executed the first time
                // this will be executed if the elapsedTime is not set correctly to current time
                fmt.Println("elapsedTime less than current")
                sf.elapsedTime = current
                sf.sequence = 0
        } else {
                sf.sequence = (sf.sequence + 1) & maskSequence
                if sf.sequence == 0  {
                        sf.elapsedTime++
                        overtime := sf.elapsedTime - current
                        fmt.Println("SLEEP FOR DURATION: ", sleepTime(overtime))
                        time.Sleep(sleepTime((overtime)))
                }
        }
        lower, err := sf.toID()
        if (err != nil) {
                return 0, 0, err
        }
        sf.sequence = uint16(1 << BitLenSequence - 1)
        upper, err := sf.toID()
        if (err != nil) {
                return 0, 0, err
        }
        return lower, upper, nil

}

const snowFlakeTimeUnitScaleFactor = 1e7 // nsec, i.e. 10 msec convert unit of nano-sec to 10 msec.

func toSnowFlakeTime(t time.Time) int64 {
        return t.UTC().UnixNano() / snowFlakeTimeUnitScaleFactor
}

func currentElapsedTime(startTime int64) int64 {
        return toSnowFlakeTime(time.Now()) - startTime
}

func sleepTime(overtime int64) time.Duration {
        return time.Duration(overtime)*10*time.Millisecond -
                time.Duration(time.Now().UTC().UnixNano() % snowFlakeTimeUnitScaleFactor) * time.Nanosecond
}

func (sf *SnowFlake) toID() (uint64, error) {
        if sf.elapsedTime >= 1<<BitLenTime {
                return 0, errors.New("over the time limit")
        }
        // Time-Sequence-MachineID
        //return uint64(sf.elapsedTime)<<(BitLenSequence+BitLenMachineID) |
        //        uint64(sf.sequence)<<BitLenMachineID |
        //        uint64(sf.machineID), nil

        // Time-MachineID-Sequence
        return uint64(sf.elapsedTime) << (BitLenSequence + BitLenMachineID) |
                uint64(sf.machineID) << BitLenSequence |
                uint64(sf.sequence), nil
}

func privateIPv4() (net.IP, error) {
        as, err := net.InterfaceAddrs()
        if err != nil {
                return nil, err
        }

        for _, a := range as {
                ipnet, ok := a.(*net.IPNet)
                if !ok || ipnet.IP.IsLoopback() {
                        continue
                }

                ip := ipnet.IP.To4()
                if isPrivateIPv4(ip) {
                        return ip, nil
                }
        }
        return nil, errors.New("no private ip address")
}

func amazonEC2PrivateIPv4() (net.IP, error) {
        // URL to retrieve instance metadata in an AWS EC2 instance:
        // http://docs.aws.amazon.com/en_us/AWSEC2/latest/UserGuide/ec2-instance-metadata.html
        timeout := time.Duration( 10 * time.Millisecond)
        client := http.Client{
                Timeout: timeout,
        }
        res, err := client.Get("http://169.254.169.254/latest/meta-data/local-ipv4")
        if err != nil {
                return nil, err
        }
        defer res.Body.Close()

        body, err := ioutil.ReadAll(res.Body)
        if err != nil {
                return nil, err
        }

        ip := net.ParseIP(string(body))
        if ip == nil {
                return nil, errors.New("invalid ip address")
        }
        return ip.To4(), nil
}

func k8sPodIPFromEnvVariable() (net.IP, error) {
        podIpEnvVarKey := "UNIQUE_ID_POD_IP"
        podIpStr := os.Getenv(podIpEnvVarKey)
        if podIpStr == "" {
                return nil, errors.New("Env Variable Not Present")
        }
        ip := net.ParseIP(podIpStr)
        if ip == nil {
                return nil, errors.New("invalid ip address")
        }
        return ip.To4(), nil
}

func isPrivateIPv4(ip net.IP) bool {
        return ip != nil &&
                (ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func lower16BitPrivateIP() (uint16, error) {
        ip, err := k8sPodIPFromEnvVariable()
        if err != nil {
                ip, err = amazonEC2PrivateIPv4()
        }
        if err != nil {
                ip, err = privateIPv4()
        }
        if err != nil {
                return 0, err
        }

        return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

// Decompose returns a set of SnowFlake ID parts.
// Time-MachineID-Sequence
func decompose(id uint64) map[string]uint64 {
        // const maskSequence = uint64((1<<BitLenSequence - 1) << BitLenMachineID)
        const maskSequence = uint64((1<<BitLenSequence - 1))
        //const maskMachineID = uint64(1<<BitLenMachineID - 1)
        const maskMachineID = uint64((1<<BitLenMachineID - 1) << BitLenSequence )
        msb := id >> 63
        time := id >> (BitLenSequence + BitLenMachineID)
        //sequence := id & maskSequence >> BitLenMachineID
        sequence := id & maskSequence

        machineID := id & maskMachineID >> BitLenSequence

        return map[string]uint64{
                "id":         id,
                "msb":        msb,
                "time":       time,
                "sequence":   sequence,
                "machine-id": machineID,
        }
}

