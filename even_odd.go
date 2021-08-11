package main

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "os"
    "strconv"
    "sync"
)

const number_of_threads int = 2
const file_name string = "test_data.csv"

type numData struct {
    number_str string
}

type objNumber struct {
     Sum int `json:"sum"`
     Count int `json:"count"`
}

type Result struct {
    FileName string `json:"fileName"`
    EvenNumber objNumber `json:"evenNumber"`
    OddNumber objNumber `json:"oddNumber"`
}

var result = Result {
    FileName: file_name,
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func read_from_file(file string, out chan<- int) {
    csvFile, err := os.Open(file)
    check(err)
    defer csvFile.Close()
    defer close(out)

    reader := csv.NewReader(csvFile)
    reader.Comma = ';'

    for {
        line, err := reader.Read()
        if err == io.EOF {
            fmt.Println("---END OF FILE---")
            break
        }
        check(err)
        num, err := strconv.Atoi(line[0])
        check(err)
        out <- num
    }
}

func handle_even_odd(wg *sync.WaitGroup, in <-chan int, n_thread int) {
    var counter = 0
    defer wg.Done()
    fmt.Println("START N_THREAD", n_thread)

    for {
    num, opened := <-in
    if opened {
        counter ++
    } else {
        fmt.Printf("N_THREAD %d PROCESSED %d NUMBERS\n", n_thread, counter)
        break
    }

    switch num % 2 == 0 {
    case true:
        result.EvenNumber.Sum += num
        result.EvenNumber.Count ++
    case false:
        result.OddNumber.Sum += num
        result.OddNumber.Count ++
    }
    }
}

func main() {
    num_chan := make(chan int)
    var wg sync.WaitGroup

    wg.Add(number_of_threads)

    go read_from_file(result.FileName, num_chan)

    for n := 0; n < number_of_threads; n++ {
    go handle_even_odd(&wg, num_chan, n)
    }

    wg.Wait()

    result_json, err := json.Marshal(result)
    check(err)
    fmt.Println(string(result_json))
}
