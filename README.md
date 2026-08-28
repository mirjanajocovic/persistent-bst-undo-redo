# Perzistentno binarno stablo pretrage u Go-u

Ovaj repozitorijum sadrži izvorni kôd praktične implementacije perzistentnog binarnog stabla pretrage (engl. _Persistent Binary Search Tree_), razvijen kao deo master rada na temu perzistentnih struktura podataka.

Struktura je inicijalno implementirana primenom tehnike **kopiranja putanja** (engl. _path copying_). Kao dodatak i referenca za uporednu analizu u master radu, implementirana je i **metoda debelih čvorova** (engl. _fat nodes_).

## Ključne funkcionalnosti

- **Kopiranje putanja (Path Copying):** Operacije umetanja i brisanja (`InsertPersistent`, `DeletePersistent`) ne uništavaju prethodno stanje strukture, već kreiraju nove čvorove na putanji izmene i dele neizmenjena podstabla (_structural sharing_).
- **Metoda debelih čvorova (Fat Nodes):** Umetanje (`InsertFat`) čuva prethodna stanja mutiranjem postojećih čvorova na način da beleži niz verzionisanih pokazivača umesto samo jednog važećeg. Pretraga u istoriji vrši se hronološkim skeniranjem važećih pokazivača.
- **Efemerna alternativa:** Repozitorijum uključuje i klasičnu, destruktivnu verziju stabla (`InsertEphemeral`, `DeleteEphemeral`) radi preciznog uporednog testiranja i analize performansi.
- **Undo/Redo mehanizam:** `HistoryManager` čuva korenove svih verzija, omogućavajući prelazak kroz istoriju stanja (Undo/Redo) u strogoj **$O(1)$ vremenskoj složenosti**, bez dubokog kopiranja ili rekonstrukcije stabla (važi za Path Copying tehniku).

## Struktura koda

- `bst.go` - Definicija čvora i implementacija efemernih i perzistentnih operacija tehnikom kopiranja putanja.
- `fat_node.go` - Strukture i operacije nad stablom zasnovane na metodi debelih čvorova.
- `history.go` - Implementacija `HistoryManager` strukture za $O(1)$ kretanje kroz verzije.
- `main.go` - Skripta za _benchmark_ testiranje. Generiše nizove, meri vreme izvršavanja i utrošak memorije za sve tri metode (Efemerno, Path Copying, Fat Nodes).

## Pokretanje

Za pokretanje benchmark testa i analizu performansi, potrebno je imati instaliran [Go](https://golang.org/dl/). Klonirajte repozitorijum i pokrenite projekat iz terminala:

```bash
git clone [https://github.com/TVOJ_USERNAME/go-persistent-bst.git](https://github.com/TVOJ_USERNAME/go-persistent-bst.git)
cd go-persistent-bst
go run .
```
