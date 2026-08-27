# Perzistentno binarno stablo pretrage u Go-u

Ovaj repozitorijum sadrži izvorni kôd praktične implementacije perzistentnog binarnog stabla pretrage (engl. _Persistent Binary Search Tree_), razvijen kao deo master rada na temu perzistentnih struktura podataka.

Struktura je implementirana u programskom jeziku **Go** primenom tehnike **kopiranja putanja** (engl. _path copying_). Pored same strukture, repozitorijum sadrži i `HistoryManager` koji omogućava efikasno upravljanje istorijom verzija.

## Ključne funkcionalnosti

- **Perzistentnost:** Operacije umetanja i brisanja (`InsertPersistent`, `DeletePersistent`) ne uništavaju prethodno stanje strukture, već kreiraju nove čvorove na putanji izmene i dele neizmenjena podstabla (_structural sharing_).
- **Efemerna alternativa:** Repozitorijum uključuje i klasičnu, destruktivnu verziju stabla (`InsertEphemeral`, `DeleteEphemeral`) radi preciznog uporednog testiranja i analize performansi.
- **Undo/Redo mehanizam:** `HistoryManager` čuva korenove svih verzija, omogućavajući prelazak kroz istoriju stanja (Undo/Redo) u strogoj **$O(1)$ vremenskoj složenosti**, bez dubokog kopiranja ili rekonstrukcije stabla.
- **Zaštita od divergentne istorije:** Sistem automatski prepoznaje nove izmene nakon `Undo` operacija, odseca staru budućnost i omogućava sakupljaču otpadaka (_Garbage Collector_) da oslobodi napuštene čvorove.

## Struktura koda

- `bst.go` - Definicija čvora i implementacija efemernih i perzistentnih operacija (umetanje, brisanje, pretraga).
- `history.go` - Implementacija `HistoryManager` strukture za $O(1)$ kretanje kroz verzije.
- `main.go` - Skripta za _benchmark_ testiranje. Generiše nizove do milion elemenata, meri vreme izvršavanja i utrošak memorije.

## Pokretanje

Za pokretanje benchmark testa i analizu performansi, potrebno je imati instaliran [Go](https://golang.org/dl/). Klonirajte repozitorijum i pokrenite projekat iz terminala:

```bash
git clone [https://github.com/TVOJ_USERNAME/go-persistent-bst.git](https://github.com/TVOJ_USERNAME/go-persistent-bst.git)
cd go-persistent-bst
go run .
```
