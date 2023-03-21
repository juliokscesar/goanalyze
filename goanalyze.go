package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
)

func getMeanF64(data []float64) float64 {
	mean := float64(0)
	for _, val := range data {
		mean += val
	}

	mean /= float64(len(data))
	return mean
}

func getMedianF64(data []float64) float64 {
	var median float64

	dataLen := len(data)
	if (dataLen % 2) == 0 {
		median = (data[dataLen/2] + data[dataLen/2-1]) / 2
	} else {
		median = data[(dataLen-1)/2]
	}

	return median
}

func getModeF64(data []float64) float64 {
	valRepeatCount := make(map[float64]float64)

	for _, val := range data {
		if _, ok := valRepeatCount[val]; !ok {
			valRepeatCount[val] = float64(0)
		} else {
			valRepeatCount[val]++
		}
	}

	mode := float64(data[0])
	for val, count := range valRepeatCount {
		if valRepeatCount[mode] < count {
			mode = val
		}
	}

	return mode
}

func getStdDeviationF64(data []float64, mean float64) float64 {
	stdDeviation := float64(0)
	for _, val := range data {
		stdDeviation += (val - mean) * (val - mean)
	}

	stdDeviation /= float64(len(data))
	stdDeviation = math.Sqrt(stdDeviation)

	return stdDeviation
}

func getMeanDeviationF64(data []float64, mean float64) float64 {
	meanDeviation := float64(0)

	for _, val := range data {
		meanDeviation += math.Abs(val - mean)
	}

	meanDeviation /= float64(len(data))

	return meanDeviation
}

func convertArrStrToF64(arr []string) []float64 {
	var numbers []float64

	for _, strNum := range arr {
		n, err := strconv.ParseFloat(strNum, 64)
		if err == nil {
			numbers = append(numbers, n)
		} else {
			log.Println(err)
		}
	}

	return numbers
}

func analyzeRawData(data [][]string) []map[string]float64 {
	allAnalysis := make([]map[string]float64, len(data))

	for lineIndex, line := range data {
		curLineVals := make(map[string]float64)
		curLineVals["LineIndex"] = float64(lineIndex)

		numValues := convertArrStrToF64(line)
		curLineVals["Amount"] = float64(len(numValues))
		sort.Float64s(numValues)

		curLineVals["Min"] = numValues[0]
		curLineVals["Max"] = numValues[len(numValues)-1]
		curLineVals["Amplitude"] = numValues[len(numValues)-1] - numValues[0]

		curLineVals["Mean"] = getMeanF64(numValues)
		curLineVals["Median"] = getMedianF64(numValues)
		curLineVals["Mode"] = getModeF64(numValues)

		curLineVals["StandardDeviation"] = getStdDeviationF64(numValues, curLineVals["Mean"])
		curLineVals["MeanDeviation"] = getMeanDeviationF64(numValues, curLineVals["Mean"])

		allAnalysis = append(allAnalysis, curLineVals)
	}

	return allAnalysis
}

func printAnalysis(filePath string) {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("%s does not exist.", filePath)
		}
	}

	file, err := os.OpenFile(filePath, os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	data, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	allAnalysis := analyzeRawData(data)
	for _, analysis := range allAnalysis {
		fmt.Printf("Line Index Analysis: %.0f\nAmount of Values: %.0f\nMaximum: %f\nMinimum: %f\nMean: %f\nMode: %f\nMedian: %f\nMean Deviation: %f\nStandard Deviation: %f\nAmplitude: %f\n\n",
			analysis["LineIndex"],
			analysis["Amount"],
			analysis["Max"],
			analysis["Min"],
			analysis["Mean"],
			analysis["Mode"],
			analysis["Median"],
			analysis["MeanDeviation"],
			analysis["StandardDeviation"],
			analysis["Amplitude"],
		)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: goanalyze <CSV file>")
	}

	dataFile := os.Args[1]
	printAnalysis(dataFile)
}
