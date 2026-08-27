package main

// HistoryManager upravlja istorijom verzija perzistentnog stabla.
//
// Svaka verzija stabla predstavljena je pokazivačem na koren. Pošto je stablo
// perzistentno, čuvanje korena je dovoljno da se očuva kompletno prethodno stanje.
// Operacije Undo i Redo ne rekonstruišu stablo, već samo pomeraju indeks trenutne
// verzije, zbog čega rade u O(1) vremenu.
type HistoryManager struct {
	history      []*Node
	currentIndex int
}

// NewHistoryManager kreira novi menadžer istorije.
//
// Početna verzija može biti nil, što predstavlja prazno stablo. Na taj način
// i prazno stablo ima svoju eksplicitnu verziju u istoriji.
func NewHistoryManager(initialRoot *Node) *HistoryManager {
	return &HistoryManager{
		history:      []*Node{initialRoot},
		currentIndex: 0,
	}
}

// AddVersion dodaje novu verziju stabla u istoriju.
//
// Ključni granični slučaj nastaje kada se nova verzija dodaje nakon jedne ili više
// Undo operacija. Tada currentIndex nije na kraju niza, što znači da postoji
// takozvana buduća istorija. Ta buduća istorija mora biti odsečena pre dodavanja
// nove verzije, jer više ne pripada aktuelnoj grani izvršavanja.
func (h *HistoryManager) AddVersion(root *Node) {
	// Granični slučaj: ako HistoryManager nije inicijalizovan konstruktorom,
	// prva dodata verzija postaje početna i trenutno aktivna verzija.
	if len(h.history) == 0 {
		h.history = []*Node{root}
		h.currentIndex = 0
		return
	}

	// Zaštita od divergentne istorije:
	// ako trenutna verzija nije poslednja, sve verzije posle nje se odbacuju.
	if h.currentIndex < len(h.history)-1 {
		// Pre skraćivanja niza uklanjaju se reference na buduće korenove.
		// Time se omogućava sakupljaču otpadaka da oslobodi čvorove do kojih
		// više ne vodi nijedna aktivna verzija.
		for i := h.currentIndex + 1; i < len(h.history); i++ {
			h.history[i] = nil
		}

		// Odseca se stara budućnost. Nakon ovoga Redo više nije moguć
		// prema staroj grani istorije.
		h.history = h.history[:h.currentIndex+1]
	}

	// Nova verzija se dodaje na kraj istorije i postaje trenutno aktivna.
	h.history = append(h.history, root)
	h.currentIndex = len(h.history) - 1
}

// Undo vraća koren prethodne verzije stabla.
//
// Ako prethodna verzija postoji, currentIndex se pomera unazad.
// Ako je trenutna verzija već prva verzija, funkcija bezbedno vraća postojeći
// koren i ne menja indeks.
func (h *HistoryManager) Undo() *Node {
	// Granični slučaj: neinicijalizovana ili prazna istorija.
	if len(h.history) == 0 {
		return nil
	}

	// Zaštita: nije moguće otići pre prve verzije.
	if h.currentIndex == 0 {
		return h.history[h.currentIndex]
	}

	h.currentIndex--
	return h.history[h.currentIndex]
}

// Redo vraća koren naredne verzije stabla.
//
// Ako naredna verzija postoji, currentIndex se pomera unapred.
// Ako je trenutna verzija već poslednja verzija, funkcija bezbedno vraća
// postojeći koren i ne menja indeks.
func (h *HistoryManager) Redo() *Node {
	// Granični slučaj: neinicijalizovana ili prazna istorija.
	if len(h.history) == 0 {
		return nil
	}

	// Zaštita: nije moguće otići posle poslednje verzije.
	if h.currentIndex == len(h.history)-1 {
		return h.history[h.currentIndex]
	}

	h.currentIndex++
	return h.history[h.currentIndex]
}

// Current vraća koren trenutno aktivne verzije stabla.
//
// Ova metoda nije neophodna za Undo/Redo mehanizam, ali je korisna u testovima
// i eksperimentima, jer omogućava čitanje trenutnog stanja bez promene istorije.
func (h *HistoryManager) Current() *Node {
	// Granični slučaj: neinicijalizovana ili prazna istorija.
	if len(h.history) == 0 {
		return nil
	}

	return h.history[h.currentIndex]
}