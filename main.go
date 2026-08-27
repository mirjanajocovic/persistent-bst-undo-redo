package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

var sinkNode *Node
var sinkBool bool

func readMem() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.TotalAlloc
}

func bytesToMB(b uint64) float64 {
	return float64(b) / 1024.0 / 1024.0
}

func main() {
	testSizes := []int{10000, 50000, 100000, 500000, 1000000}

	fmt.Println("POČETAK BENCHMARK TESTIRANJA...")
	fmt.Println("---------------------------------------------------------")

	for _, n := range testSizes {
		rand.Seed(time.Now().UnixNano())
		values := make([]int, n)
		for i := 0; i < n; i++ {
			values[i] = rand.Intn(n * 10)
		}

		fmt.Printf("\n>>> TEST ZA N = %d <<<\n", n)

		// --- 1. EFEMERNO STABLO ---
		runtime.GC()
		memBeforeEph := readMem()

		// Umetanje
		startEphIns := time.Now()
		var ephRoot *Node
		for _, val := range values {
			ephRoot = InsertEphemeral(ephRoot, val)
		}
		durEphIns := time.Since(startEphIns)

		// Pretraga
		startEphSearch := time.Now()
		for _, val := range values {
			sinkBool = Search(ephRoot, val)
		}
		durEphSearch := time.Since(startEphSearch)

		// Brisanje
		startEphDel := time.Now()
		for _, val := range values {
			ephRoot = DeleteEphemeral(ephRoot, val)
		}
		durEphDel := time.Since(startEphDel)

		runtime.GC()
		memEphMB := bytesToMB(readMem() - memBeforeEph)
		sinkNode = ephRoot

		// --- 2. PERZISTENTNO STABLO ---
		runtime.GC()
		memBeforePers := readMem()

		// Umetanje
		startPersIns := time.Now()
		var persRoot *Node
		hm := NewHistoryManager(persRoot)
		for _, val := range values {
			persRoot = InsertPersistent(persRoot, val)
			hm.AddVersion(persRoot)
		}
		durPersIns := time.Since(startPersIns)

		// Pretraga (ista funkcija, nad perzistentnim korenom)
		startPersSearch := time.Now()
		for _, val := range values {
			sinkBool = Search(persRoot, val)
		}
		durPersSearch := time.Since(startPersSearch)

		// Brisanje
		startPersDel := time.Now()
		for _, val := range values {
			persRoot = DeletePersistent(persRoot, val)
			hm.AddVersion(persRoot)
		}
		durPersDel := time.Since(startPersDel)

		runtime.GC()
		memPersMB := bytesToMB(readMem() - memBeforePers)
		sinkNode = persRoot

		// --- ISPIS ---
		fmt.Printf("UMETANJE | Efemerno: %-10s | Perzistentno: %-10s\n", durEphIns, durPersIns)
		fmt.Printf("PRETRAGA | Efemerno: %-10s | Perzistentno: %-10s\n", durEphSearch, durPersSearch)
		fmt.Printf("BRISANJE | Efemerno: %-10s | Perzistentno: %-10s\n", durEphDel, durPersDel)
		fmt.Printf("MEMORIJA | Efemerno: %-7.2f MB | Perzistentno: %-7.2f MB\n", memEphMB, memPersMB)
	}

	fmt.Println("\n---------------------------------------------------------")
	fmt.Println("Testiranje Undo operacije za N = 1.000.000...")
	var root *Node
	hm := NewHistoryManager(root)
	for i := 0; i < 1000000; i++ {
		root = InsertPersistent(root, rand.Intn(1000000))
		hm.AddVersion(root)
	}

	startUndo := time.Now()
	for i := 0; i < 999999; i++ {
		_ = hm.Undo()
	}
	durationUndo := time.Since(startUndo)

	fmt.Printf("Vreme za 1.000.000 Undo operacija: %v\n", durationUndo)
}