package mapreduce

import (
    "io"
    "os"
    "encoding/json"
    "log"
)

// doReduce manages one reduce task: it reads the intermediate
// key/value pairs (produced by the map phase) for this task, sorts the
// intermediate key/value pairs by key, calls the user-defined reduce function
// (reduceF) for each key, and writes the output to disk.
func doReduce(
	jobName string, // the name of the whole MapReduce job
	reduceTaskNumber int, // which reduce task this is
	outFile string, // write the output here
	nMap int, // the number of map tasks that were run ("M" in the paper)
	reduceF func(key string, values []string) string,
) {
    var mapReduceKeyValues map[string][]string
    mapReduceKeyValues = make(map[string][]string)
    for i := 0; i < nMap; i++ {
        reduceFileName := reduceName(jobName, i, reduceTaskNumber)
        file, err := os.Open(reduceFileName)
        if err != nil {
            log.Fatalf("Open %s err:%s", reduceFileName, err)
        }
        defer file.Close()
        var kv KeyValue
        dec := json.NewDecoder(file)
        for {
            if err = dec.Decode(&kv); err != nil {
                if err == io.EOF {
                    break
                } else {
                    // panic(err)
                    log.Fatalf("Decode %s failed err:%s", kv, err)
                }
            }
            mapReduceKeyValues[kv.Key] = append(mapReduceKeyValues[kv.Key], kv.Value)
        }
    }

    mergeFileName := mergeName(jobName, reduceTaskNumber)
    file, err := os.OpenFile(mergeFileName, os.O_RDWR|os.O_CREATE, 0644)
    if err != nil {
        // panic(err)
        log.Fatalf("OpenFile %s failed err:%s", mergeFileName, err)
    }
    defer file.Close()
    enc := json.NewEncoder(file)
    var curReduceResult string
    for reduceKey,reduceValue := range mapReduceKeyValues {
	    // reduceF func(key string, values []string) string,
        curReduceResult = reduceF(reduceKey, reduceValue)
        if err = enc.Encode(&KeyValue{reduceKey, curReduceResult}); err != nil {
            // panic(err)
            log.Fatalf("Encode key:%s value:%s failed err:%s", reduceKey, curReduceResult, err)
        }
    }

	//
	// You will need to write this function.
	//
	// You'll need to read one intermediate file from each map task;
	// reduceName(jobName, m, reduceTaskNumber) yields the file
	// name from map task m.
	//
	// Your doMap() encoded the key/value pairs in the intermediate
	// files, so you will need to decode them. If you used JSON, you can
	// read and decode by creating a decoder and repeatedly calling
	// .Decode(&kv) on it until it returns an error.
	//
	// You may find the first example in the golang sort package
	// documentation useful.
	//
	// reduceF() is the application's reduce function. You should
	// call it once per distinct key, with a slice of all the values
	// for that key. reduceF() returns the reduced value for that key.
	//
	// You should write the reduce output as JSON encoded KeyValue
	// objects to the file named outFile. We require you to use JSON
	// because that is what the merger than combines the output
	// from all the reduce tasks expects. There is nothing special about
	// JSON -- it is just the marshalling format we chose to use. Your
	// output code will look something like this:
	//
	// enc := json.NewEncoder(file)
	// for key := ... {
	// 	enc.Encode(KeyValue{key, reduceF(...)})
	// }
	// file.Close()
	//
}
