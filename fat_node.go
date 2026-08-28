package main

// VersionedNode predstavlja zapis o pokazivaču u određenoj verziji.
// Umesto prepisivanja starog pokazivača, "debeli" čvor čuva niz ovakvih zapisa.
type VersionedNode struct {
	Version int
	Node    *FatNode
}

// FatNode predstavlja "debeli" čvor u perzistentnom stablu.
// Umesto samo jednog Left i Right pokazivača, ovaj čvor sadrži istoriju
// promena (nizove) za svoja pokazivačka polja.
type FatNode struct {
	Value int
	Left  []VersionedNode
	Right []VersionedNode
}

// getLatest pronalazi važeći pokazivač za traženu verziju.
// Pošto se verzije uvek dodaju hronološki, niz je sortiran, pa pretraga unazad
// brzo pronalazi najnoviju promenu koja nije mlađa od tražene verzije.
func getLatest(history []VersionedNode, version int) *FatNode {
	for i := len(history) - 1; i >= 0; i-- {
		if history[i].Version <= version {
			return history[i].Node
		}
	}
	return nil
}

// InsertFat umeće novu vrednost metodom debelih čvorova.
// Za razliku od kopiranja putanja, ova metoda mutira postojeće čvorove dodajući
// nove zapise u njihove Left/Right nizove, uz obeležavanje verzije.
func InsertFat(root *FatNode, val int, version int) *FatNode {
	if root == nil {
		return &FatNode{Value: val}
	}

	curr := root
	for {
		if val < curr.Value {
			latestLeft := getLatest(curr.Left, version)
			if latestLeft == nil {
				// Nema levog deteta u ovoj verziji, dodajemo ga u istoriju.
				curr.Left = append(curr.Left, VersionedNode{
					Version: version,
					Node:    &FatNode{Value: val},
				})
				return root
			}
			curr = latestLeft
		} else if val > curr.Value {
			latestRight := getLatest(curr.Right, version)
			if latestRight == nil {
				// Nema desnog deteta u ovoj verziji, dodajemo ga u istoriju.
				curr.Right = append(curr.Right, VersionedNode{
					Version: version,
					Node:    &FatNode{Value: val},
				})
				return root
			}
			curr = latestRight
		} else {
			// Ignorišemo duplikate
			return root
		}
	}
}

// SearchFat pretražuje stablo metodom debelih čvorova za zadatu verziju.
// Zbog pretrage po nizu pokazivača, vreme pretrage je asimptotski malo sporije
// u odnosu na O(1) pristup kod kopiranja putanja.
func SearchFat(root *FatNode, val int, version int) bool {
	curr := root
	for curr != nil {
		if val < curr.Value {
			curr = getLatest(curr.Left, version)
		} else if val > curr.Value {
			curr = getLatest(curr.Right, version)
		} else {
			return true
		}
	}
	return false
}