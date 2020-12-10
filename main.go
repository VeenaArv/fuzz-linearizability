package main

import (
	"fmt"
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
			stats := fuzzing.CheckLinearizability(string(input), false /*strongReadConsistency*/, false /*delays */, run, i)
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
	for i := 8; i < 15; i++ {
		fuzzing.RandomizedTesting(20, i, false, i-7)
	}
	// for i := 8; i < 15; i++ {
	// 	fuzzing.RandomizedTesting(20, i, false, i)
	// }
	// fuzzing.RandomizedTestingWithDelays(15, false, 2)
	// data, _ := ioutil.ReadFile("data/custom_test.txt")
	// content := string(data)
	// testCaseStats := fuzzing.CheckLinearizability(content, false, 3, 1)
	// fuzzing.WriteStats(testCaseStats, 3, 1)
}
