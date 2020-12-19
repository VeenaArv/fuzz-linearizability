package main

import (
	"fmt"

	// "fuzz-linearizability/rqlite"
	"fuzz-linearizability/fuzzing"
	"io/ioutil"
)

func runSampleInput() {
	files := []string{"go_fuzz_integration/corpus/input1.txt",
		"go_fuzz_integration/corpus/input2.txt",
		"go_fuzz_integration/corpus/input3.txt"}
	var runs []fuzzing.RunStats

	for run := 0; run < 10; run++ {
		var testcases []fuzzing.TestCaseStats
		for i, file := range files {
			input, _ := ioutil.ReadFile(file)
			stats := fuzzing.CheckLinearizability(string(input),
				fuzzing.AlgoRunParams{0 /*NumEvents*/, 0 /*NumTests*/, run, /*Run*/
					false /*StrongReadConsistency*/, false, /*Delays*/
					"sample" /*Version*/}, i)
			testcases = append(testcases, stats)
		}
		runs = append(runs, fuzzing.CalcRunStats(testcases))
		fmt.Println("runs")
		fmt.Println(runs)
	}
	stats := fuzzing.CalcAvgStats(runs)
	fmt.Println("avg stats")
	fmt.Println(stats)
	// fmt.Println(rqlite.TestHistory())

}

func main() {
	// fmt.Println(rqlite.TestHistory())
	// runSampleInput()
	// for i := 10; i < 30; i++ {
	// fuzzing.RandomizedTesting(fuzzing.AlgoRunParams{i /*NumEvents*/, 10 /*NumTests*/, 1 /*Run*/, false /*StrongReadConsistency*/, false /*Delays*/, "random" /*Version*/})
	// }

	fuzzing.RandomizedTesting(fuzzing.AlgoRunParams{20 /*NumEvents*/, 1 /*NumTests*/, 1 /*Run*/, false /*StrongReadConsistency*/, false /*Delays*/, "random" /*Version*/})

	// fuzzing.GeneticAlgoWithIncreasingTestCases(fuzzing.AlgoRunParams{5 /*NumEvents*/, 10 /*NumTests*/, 2 /*Run*/, false /*StrongReadConsistency*/, false /*Delays*/, "genetic" /*Version*/})
	// fuzzing.GeneticAlgo(fuzzing.AlgoRunParams{70 /*NumEvents*/, 10 /*NumTests*/, 3 /*Run*/, false /*StrongReadConsistency*/, false /*Delays*/, "genetic_2" /*Version*/})
	// fuzzing.RandomizedTestingWithDelays(15, false, 2)
	// data, _ := ioutil.ReadFile("output/histories/random/history_0.txt")
	// content := string(data)
	// testCaseStats := fuzzing.CheckLinearizability(content, false, 3, 1)
	// historyFilePath := "1.txt"
	// linearizable := rqlite.CheckHistory(historyFilePath, false /*delFile*/)
	// fmt.Println(linearizable)
	// fuzzing.WriteStats(testCaseStats, 3, 1)
}
