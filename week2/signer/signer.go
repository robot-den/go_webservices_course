package main

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

const WithLog = false

func main() {
	jobs := []job{
		job(SimpleInput),
		job(SingleHash),
		job(MultiHash),
		job(CombineResults),
		job(SimpleOutput),
	}
	ExecutePipeline(jobs...)
}

func ExecutePipeline(jobs ...job) {
	in := make(chan interface{}, 100)
	out := make(chan interface{}, 100)

	for i, job := range jobs {
		if i == 0 {
			go withChannelClosing(in, out, job)
		} else {
			in = out
			out = make(chan interface{}, 100)
			go withChannelClosing(in, out, job)
		}
	}
	for _ = range out {}
}

func SimpleInput(_in, out chan interface{}) {
	log(      "Worker started: SimpleInput")
	defer log("Worker stopped: SimpleInput")

	for i := 0; i < 2; i++ {
		out <- i
		log("SimpleInput: ", i)
	}
}

func SimpleOutput(in, _out chan interface{}) {
	log(      "Worker started: SimpleOutput")
	defer log("Worker stopped: SimpleOutput")

	for value := range in {
		newValue := fmt.Sprint(value) + "(output)"
		log("SimpleOutput: ", newValue)
	}
}

// Calc crc32(value)+"~"+crc32(md5(value)) for each value and send it to out channel
func SingleHash(in, out chan interface{}) {
	log(      "Worker started: SingleHash")
	defer log("Worker stopped: SingleHash")

	var md5Mutex = &sync.Mutex{}
	var wg = &sync.WaitGroup{}

	for value := range in {
		wg.Add(1)
		stringValue := fmt.Sprint(value)
		go SingleHashWorker(stringValue, out, md5Mutex, wg)
	}
	wg.Wait()
}

func SingleHashWorker(value string, out chan interface{}, md5Mutex *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	md5Value := dataSignerMd5WithLock(value, md5Mutex)
	multiArr := [2]string{}
	var subMu = &sync.Mutex{}
	var subWg = &sync.WaitGroup{}

	subWg.Add(1)
	go func(str string, sMu *sync.Mutex, sWg *sync.WaitGroup) {
		defer sWg.Done()

		crc32v := DataSignerCrc32(str)
		sMu.Lock()
		multiArr[0] = crc32v
		sMu.Unlock()
	}(value, subMu, subWg)

	subWg.Add(1)
	go func(str string, sMu *sync.Mutex, sWg *sync.WaitGroup) {
		defer sWg.Done()
		crc32v := DataSignerCrc32(str)
		sMu.Lock()
		multiArr[1] = crc32v
		sMu.Unlock()
	}(md5Value, subMu, subWg)

	subWg.Wait()
	subMu.Lock()
	newValue := multiArr[0] + "~" + multiArr[1]
	subMu.Unlock()

	out <- newValue
	log("SingleHashWorker: ", newValue)
}

func MultiHash(in, out chan interface{})  {
	log(      "Worker started: MultiHash")
	defer log("Worker stopped: MultiHash")

	var wg = &sync.WaitGroup{}

	for value := range in {
		wg.Add(1)
		stringValue := fmt.Sprint(value)
		go MultiHashWorker(stringValue, out, wg)
	}
	wg.Wait()
}

func MultiHashWorker(value string, out chan interface{}, wg *sync.WaitGroup) {
	defer wg.Done()

	multiArr := [6]string{}
	var mu = &sync.Mutex{}
	var subWg = &sync.WaitGroup{}
	for i := 0; i <= 5; i++ {
		subWg.Add(1)
		go func(index int, val string, arr *[6]string, mu *sync.Mutex, subWg *sync.WaitGroup) {
			defer subWg.Done()

			crc32SubValue := DataSignerCrc32(fmt.Sprint(index) + val)
			mu.Lock()
			arr[index] = crc32SubValue
			mu.Unlock()
			log(value, "Multi: ", crc32SubValue)
		}(i, value, &multiArr, mu, subWg)
	}
	newValue := ""
	subWg.Wait()
	mu.Lock()
	for _, v := range multiArr {
		newValue = newValue + v
	}
	mu.Unlock()
	out <- newValue
	log("MultiHashWorker: ", newValue)
}

func CombineResults(in, out chan interface{}) {
	log(      "Worker started: CombineResults")
	defer log("Worker stopped: CombineResults")

	var accumulator []string

	for value := range in {
		stringValue := fmt.Sprint(value)
		accumulator = append(accumulator, stringValue)
	}

	sort.Strings(accumulator)
	result := strings.Join(accumulator, "_")

	log("CombineResults: ", result)
	out <- result
}

func withChannelClosing(in, out chan interface{}, j job) {
	defer close(out)

	j(in, out)
}

// TODO: check overheating
func dataSignerMd5WithLock(value string, mu *sync.Mutex) string {
	mu.Lock()
	result := DataSignerMd5(value)
	mu.Unlock()
	return result
}

func log(values ...interface{}) {
	if WithLog {
		fmt.Println(values...)
	}
}
