package fuzzing

import (
	"fmt"
	"fuzz-linearizability/rqlite"
	"os"
	"strings"
	"time"
)

type TestCaseStats struct {
	numOperations int
	time          time.Duration
	linearizable  bool
}

type RunStats struct {
	testCases            []TestCaseStats
	totalTime            time.Duration
	nonLinearizableTests int
	tests                int

	// uniqueHistories int
}

type AvgStats struct {
	avgRunTime              time.Duration
	avgNonLinearizableTests float32
	avgTests                float32
	runs                    int
}

func (stats TestCaseStats) String() string {
	return fmt.Sprintf("(numOperations %d time %s linearizable %t)", stats.numOperations, stats.time, stats.linearizable)
}
func (stats RunStats) String() string {
	return fmt.Sprintf("(totalTime %s nonLinearizableTests %d tests %d)", stats.totalTime, stats.nonLinearizableTests, stats.tests)
}
func (stats AvgStats) String() string {
	return fmt.Sprintf("(avgRunTime %s avgNonLinearizableTests %f avgTests %f runs %d)",
		stats.avgRunTime, stats.avgNonLinearizableTests, stats.avgTests, stats.runs)
}

func CalcAvgStats(runStats []RunStats) AvgStats {
	totalTotalTime := new(time.Duration)
	totalNonLinearizableTests := 0
	totalTests := 0
	runs := len(runStats)
	for _, runStat := range runStats {
		*totalTotalTime += runStat.totalTime
		totalNonLinearizableTests += runStat.nonLinearizableTests
		totalTests += runStat.tests
	}
	return AvgStats{time.Duration(int64(*totalTotalTime) / int64(runs)),
		float32(totalNonLinearizableTests) / float32(runs),
		float32(totalTests) / float32(runs), runs}
}

func CalcRunStats(testCases []TestCaseStats) RunStats {
	totaltime := new(time.Duration)
	nonLinearizableTests := 0
	tests := len(testCases)
	for _, testCase := range testCases {
		if !testCase.linearizable {
			nonLinearizableTests++
		}
		*totaltime += testCase.time
	}
	return RunStats{testCases, *totaltime, nonLinearizableTests, tests}
}

func WriteStats(stats fmt.Stringer, run int, id int) {
	dirPath := fmt.Sprintf("output/stats/%T", stats)
	filePath := fmt.Sprintf("%s/run_%d", dirPath, run)
	_ = os.MkdirAll(dirPath, os.ModePerm)
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.WriteString(stats.String() + "\n")
	if err != nil {
		panic(err)
	}
}

func CheckLinearizability(input string, strongReadConsistency bool, delays bool, run int, id int) TestCaseStats {
	time := new(time.Duration)
	dirPath := fmt.Sprintf("output/histories/run_%d", run)
	filePath := fmt.Sprintf("%s/history_%d.txt", dirPath, id)
	_ = os.MkdirAll(dirPath, os.ModePerm)
	numOperations := strings.Count(input, "\n")
	// fmt.Println(input)
	// fmt.Println(numOperations)
	linearizable := checkLinearizability(input, filePath, strongReadConsistency, delays, time)
	return TestCaseStats{numOperations, *time, linearizable}
}

func checkLinearizability(input string, historyFilePath string, strongReadConsistency bool, delays bool, timeElasped *time.Duration) bool {
	defer timeTrack(time.Now(), "linearizability checking", timeElasped)
	// This applies operations to rqlite and writes history to filePath.
	rqlite.RunOperations(input, historyFilePath, strongReadConsistency /*strongReadConsistency*/, delays /*delays*/)
	// This uses porcupine to check the history in filePath and returns
	// true if linearizable.
	linearizable := rqlite.CheckHistory(historyFilePath, false /*delFile*/)
	// fmt.Println(linearizable)
	return linearizable
}

func timeTrack(start time.Time, name string, timeElasped *time.Duration) {
	elapsed := time.Since(start)
	*timeElasped = elapsed
	// fmt.Printf("%s took %s", name, elapsed)
}
