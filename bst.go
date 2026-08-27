package main

// Node predstavlja jedan čvor binarnog stabla pretrage.
//
// Ista struktura se koristi i za efemernu i za perzistentnu verziju stabla.
// Razlika nije u obliku čvora, već u načinu na koji se čvorovi menjaju:
// efemerno stablo menja postojeće pokazivače, dok perzistentno stablo
// kreira nove čvorove na putanji izmene.
type Node struct {
	Value int
	Left  *Node
	Right *Node
}

// InsertEphemeral umeće novu vrednost u efemerno binarno stablo pretrage.
//
// Ova funkcija je destruktivna: ona menja postojeću strukturu stabla tako što
// ažurira pokazivače u već postojećim čvorovima. Zbog toga prethodno stanje
// stabla nakon umetanja više nije dostupno.
//
// Duplikati se u ovoj implementaciji ignorišu, kako bi stablo sadržalo
// jedinstvene vrednosti.
func InsertEphemeral(root *Node, val int) *Node {
	// Granični slučaj: ako je stablo prazno, novi čvor postaje koren.
	if root == nil {
		return &Node{Value: val}
	}

	// Iterativna implementacija izbegava dodatno zauzeće resursa usled rekurzivnih poziva.
	current := root

	for {
		if val < current.Value {
			// Ako ne postoji levi potomak, novi čvor se dodaje na to mesto.
			if current.Left == nil {
				current.Left = &Node{Value: val}
				return root
			}

			current = current.Left
		} else if val > current.Value {
			// Ako ne postoji desni potomak, novi čvor se dodaje na to mesto.
			if current.Right == nil {
				current.Right = &Node{Value: val}
				return root
			}

			current = current.Right
		} else {
			// Granični slučaj: duplikat se ignoriše i stablo ostaje nepromenjeno.
			return root
		}
	}
}

// InsertPersistent umeće novu vrednost u perzistentno binarno stablo pretrage.
//
// Ova funkcija implementira tehniku kopiranja putanja (engl. path copying). Umesto da menja
// postojeće čvorove, ona kreira nove čvorove samo na putanji od korena do mesta umetanja.
// Sva podstabla koja nisu na toj putanji dele se između stare i nove verzije.
//
// Funkcija vraća novi koren stabla. Stari koren ostaje validan i predstavlja
// prethodnu verziju stabla.
//
// Duplikati se ignorišu. Ako vrednost već postoji, vraća se postojeći koren,
// bez nepotrebnog kreiranja novih čvorova.
func InsertPersistent(root *Node, val int) *Node {
	// Granični slučaj: ako je stablo prazno, formira se nova verzija sa jednim čvorom.
	if root == nil {
		return &Node{Value: val}
	}

	if val < root.Value {
		// Rekurzivno se kreira nova verzija levog podstabla.
		newLeft := InsertPersistent(root.Left, val)

		// Ako se levo podstablo nije promenilo, nema potrebe kopirati trenutni čvor.
		// Ovo se dešava, na primer, kada je vrednost duplikat.
		if newLeft == root.Left {
			return root
		}

		// Kopiranje putanja: kopira se trenutni čvor, menja se samo levi pokazivač,
		// dok se desno podstablo deli sa prethodnom verzijom.
		return &Node{
			Value: root.Value,
			Left:  newLeft,
			Right: root.Right,
		}
	}

	if val > root.Value {
		// Rekurzivno se kreira nova verzija desnog podstabla.
		newRight := InsertPersistent(root.Right, val)

		// Ako se desno podstablo nije promenilo, vraća se postojeći koren.
		if newRight == root.Right {
			return root
		}

		// Kopiranje putanja: kopira se trenutni čvor, menja se samo desni pokazivač,
		// dok se levo podstablo deli sa prethodnom verzijom.
		return &Node{
			Value: root.Value,
			Left:  root.Left,
			Right: newRight,
		}
	}

	// Granični slučaj: vrednost već postoji u stablu, pa se nova verzija ne kreira.
	return root
}

// Search pretražuje stablo za zadatom vrednošću.
// Pošto je ovo operacija čitanja, ista funkcija radi i za efemerno i za perzistentno stablo
// u vremenskoj složenosti O(log n).
func Search(root *Node, val int) bool {
	current := root
	for current != nil {
		if val < current.Value {
			current = current.Left
		} else if val > current.Value {
			current = current.Right
		} else {
			return true
		}
	}
	return false
}

// DeleteEphemeral destruktivno briše čvor iz efemernog stabla.
func DeleteEphemeral(root *Node, val int) *Node {
	if root == nil {
		return nil
	}

	if val < root.Value {
		root.Left = DeleteEphemeral(root.Left, val)
	} else if val > root.Value {
		root.Right = DeleteEphemeral(root.Right, val)
	} else {
		// Čvor sa jednim ili nijednim detetom
		if root.Left == nil {
			return root.Right
		}
		if root.Right == nil {
			return root.Left
		}

		// Čvor sa dvoje dece: pronalazi se najmanji u desnom podstablu
		minRight := root.Right
		for minRight.Left != nil {
			minRight = minRight.Left
		}
		root.Value = minRight.Value
		root.Right = DeleteEphemeral(root.Right, minRight.Value)
	}
	return root
}

// DeletePersistent briše čvor primenom tehnike kopiranja putanja.
// Vraća se novi koren, dok staro stablo ostaje netaknuto.
func DeletePersistent(root *Node, val int) *Node {
	if root == nil {
		return nil
	}

	if val < root.Value {
		newLeft := DeletePersistent(root.Left, val)
		if newLeft == root.Left {
			return root
		}
		return &Node{Value: root.Value, Left: newLeft, Right: root.Right}
	}

	if val > root.Value {
		newRight := DeletePersistent(root.Right, val)
		if newRight == root.Right {
			return root
		}
		return &Node{Value: root.Value, Left: root.Left, Right: newRight}
	}

	// Pronađen je čvor za brisanje
	if root.Left == nil {
		return root.Right
	}
	if root.Right == nil {
		return root.Left
	}

	// Čvor sa dvoje dece
	minRight := root.Right
	for minRight.Left != nil {
		minRight = minRight.Left
	}
	newRight := DeletePersistent(root.Right, minRight.Value)
	return &Node{Value: minRight.Value, Left: root.Left, Right: newRight}
}