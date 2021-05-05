package analysis

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path"
)

type ClientData struct {
	Level string
	ClientId int
	ClientBatchSize int
	TotalSent int
	MinLat int
	MaxLat int
	MaxLatIdx int
	AvgLat int
	P50Lat int
	P95Lat int
	P99Lat int
	Start int `json:"sendStart"`
	End int `json:"SendEnd"`
	Mid80Start int
	Mid80End int
	Mid80Dur float64
	Mid80RecvTimeDur float64
	Mid80Requests int
	Mid80Throughput float64 `json:"mid80Throughput (cmd/sec)"`
	Mid80Throughput2 float64 `json:"mid80Throughput2 (cmd/sec)"`
}

type Output struct {
	AvgDur float64
	AvgLat float64
	P50Lat float64
	P95LAT float64
	P99Lat float64
	MaxDur float64
	Mid80Throughput float64
}

func LoadClientLogs(logDirPath string) *[]ClientData {
	var numClients int
	var allData []ClientData

	files, err := ioutil.ReadDir(logDirPath)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		jsonFile, err := os.Open(path.Join(logDirPath, f.Name()))
		if err != nil {
			panic(err)
		}

		byteValue, _ := ioutil.ReadAll(jsonFile)
		var data ClientData
		json.Unmarshal(byteValue, &data)

		allData = append(allData, data)
		numClients += 1

		jsonFile.Close()
	}

	return &allData
}

func RunAnalysis(logDirPath string)  {
	allData := LoadClientLogs(logDirPath)
	numData := len(*allData)

	var sumDur, maxDur float64
	var sumAvgLat, sumP50Lat, sumP95Lat, sumP99Lat, maxMid80RecvTime, sumMid80Requests int
	//var sumMid80RecvTime int
	for _, data := range *allData {
		sumDur += data.Mid80Dur
		if data.Mid80Dur > maxDur {
			maxDur = data.Mid80Dur
		}

		sumAvgLat += data.AvgLat
		sumP50Lat += data.P50Lat
		sumP95Lat += data.P95Lat
		sumP99Lat += data.P99Lat

		mid80RecvTime := data.Mid80End - data.Mid80Start
		//sumMid80RecvTime += mid80RecvTime // in ns
		if mid80RecvTime > maxMid80RecvTime {
			maxMid80RecvTime = mid80RecvTime // in ns
		}

		sumMid80Requests += data.Mid80Requests
	}

	outputAvgDur := round(sumDur / float64(numData)) // TODO: check unit
	outputAvgLat := round(float64(sumAvgLat) / float64(numData) / math.Pow10(3))
	outputP50Lat := round(float64(sumP50Lat) / float64(numData) / math.Pow10(3))
	outputP95Lat := round(float64(sumP95Lat) / float64(numData) / math.Pow10(3))
	outputP99Lat := round(float64(sumP99Lat) / float64(numData) / math.Pow10(3))
	//avgMid80RecvTime := round(float64(sumMid80RecvTime) / float64(numData) / math.Pow10(9))
	outputMax80RecvTime := round(float64(maxMid80RecvTime) / math.Pow10(9))
	outputSumMid80Requests := sumMid80Requests * (*allData)[0].ClientBatchSize
	outputMid80Throughput := round(float64(outputSumMid80Requests) / outputMax80RecvTime)
	output := Output{
		AvgDur: outputAvgDur,
		AvgLat: outputAvgLat,
		P50Lat: outputP50Lat,
		P95LAT: outputP95Lat,
		P99Lat: outputP99Lat,
		MaxDur: maxDur,
		Mid80Throughput: outputMid80Throughput,
	}
	fmt.Printf("%+v\n", output)
	fmt.Println(output)
}

func round(input float64) float64 {
	return math.Round(input*100)/100
}