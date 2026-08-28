package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"
)

// Globalne promenljive za sprečavanje "Dead Code Elimination" optimizacije Go kompajlera.
// Kompajler bi inače obrisao kod za pretragu jer se rezultati nigde ne koriste.
var sinkNode *Node
var sinkFatNode *FatNode
var sinkBool bool

// readMem očitava trenutnu količinu alocirane memorije na hipu (Heap).
func readMem() uint64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.TotalAlloc
}

// bytesToMB konvertuje bajtove u megabajte radi lakšeg prikaza.
func bytesToMB(b uint64) float64 {
	return float64(b) / 1024.0 / 1024.0
}

func main() {
	// Definišemo niz veličina za koje ćemo vršiti benchmark testiranje.
	testSizes := []int{10000, 50000, 100000, 500000, 1000000}

	fmt.Println("POČETAK BENCHMARK TESTIRANJA...")
	fmt.Println("--------------------------------------------------------------------------------")

	for _, n := range testSizes {
		// 1. Inicijalizacija: Generisanje niza nasumičnih brojeva za trenutnu veličinu N.
		rand.Seed(time.Now().UnixNano())
		values := make([]int, n)
		for i := 0; i < n; i++ {
			values[i] = rand.Intn(n * 10)
		}

		fmt.Printf("\n>>> TEST ZA N = %d <<<\n", n)

		// --- 1. EFEMERNO STABLO ---
		// Forsiramo Garbage Collector (GC) pre merenja kako bismo imali čistu sliku.
		runtime.GC()
		memBeforeEph := readMem()

		// Umetanje (Efemerno)
		startEphIns := time.Now()
		var ephRoot *Node
		for _, val := range values {
			ephRoot = InsertEphemeral(ephRoot, val)
		}
		durEphIns := time.Since(startEphIns)

		// Pretraga (Efemerno)
		startEphSearch := time.Now()
		for _, val := range values {
			sinkBool = Search(ephRoot, val)
		}
		durEphSearch := time.Since(startEphSearch)

		// Brisanje (Efemerno)
		startEphDel := time.Now()
		for _, val := range values {
			ephRoot = DeleteEphemeral(ephRoot, val)
		}
		durEphDel := time.Since(startEphDel)

		// Merenje ukupne memorije za efemerno stablo.
		runtime.GC()
		memEphMB := bytesToMB(readMem() - memBeforeEph)
		sinkNode = ephRoot

		// --- 2. PERZISTENTNO STABLO (KOPIRANJE PUTANJA) ---
		runtime.GC()
		memBeforePers := readMem()

		// Umetanje uz beleženje u menadžer istorije (Path Copying)
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

		// Brisanje (Path Copying)
		startPersDel := time.Now()
		for _, val := range values {
			persRoot = DeletePersistent(persRoot, val)
			hm.AddVersion(persRoot)
		}
		durPersDel := time.Since(startPersDel)

		// Merenje ukupne memorije za metodu kopiranja putanja.
		runtime.GC()
		memPersMB := bytesToMB(readMem() - memBeforePers)
		sinkNode = persRoot

		// --- 3. METODA DEBELIH ČVOROVA (FAT NODES) ---
		// Dodato kao uporedna implementacija radi analize različitih tehnika perzistentnosti.
		runtime.GC()
		memBeforeFat := readMem()

		// Umetanje (Fat Nodes) - Obeležavamo verzije rastućim redosledom.
		startFatIns := time.Now()
		var fatRoot *FatNode
		for i, val := range values {
			// Verzija kreće od 1 (indeks i + 1)
			fatRoot = InsertFat(fatRoot, val, i+1)
		}
		durFatIns := time.Since(startFatIns)

		// Pretraga (Fat Nodes) - Pretražujemo najnoviju verziju stabla (verzija n).
		startFatSearch := time.Now()
		for _, val := range values {
			sinkBool = SearchFat(fatRoot, val, n)
		}
		durFatSearch := time.Since(startFatSearch)

		// Merenje ukupne memorije za metodu debelih čvorova.
		runtime.GC()
		memFatMB := bytesToMB(readMem() - memBeforeFat)
		sinkFatNode = fatRoot

		// --- ISPIS REZULTATA ---
		// (Brisanje za Fat Nodes je preskočeno u ispisu jer je akcenat bio na poređenju umetanja i memorije)
		fmt.Printf("UMETANJE | Efemerno: %-10s | Path Copying: %-10s | Fat Nodes: %-10s\n", durEphIns, durPersIns, durFatIns)
		fmt.Printf("PRETRAGA | Efemerno: %-10s | Path Copying: %-10s | Fat Nodes: %-10s\n", durEphSearch, durPersSearch, durFatSearch)
		fmt.Printf("BRISANJE | Efemerno: %-10s | Path Copying: %-10s | Fat Nodes: N/A\n", durEphDel, durPersDel)
		fmt.Printf("MEMORIJA | Efemerno: %-7.2f MB | Path Copying: %-7.2f MB | Fat Nodes: %-7.2f MB\n", memEphMB, memPersMB, memFatMB)
	}

	// 4. Testiranje operacija Undo i Redo
	fmt.Println("\n--------------------------------------------------------------------------------")
	fmt.Println("Testiranje brzine Undo operacije za N = 1.000.000 (Path Copying)...")
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
	fmt.Printf("Vreme za povratak kroz 1.000.000 stanja (Undo): %v\n", durationUndo)
}