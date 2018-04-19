package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func getSession() *mgo.Session {
	// Direccion del stefulset replicaset de mongo
	session, err := mgo.Dial("mongodb://ml-cluster-mongodb-replicaset-0.ml-cluster-mongodb-replicaset,ml-cluster-mongodb-replicaset-1.ml-cluster-mongodb-replicaset,ml-cluster-mongodb-replicaset-2.ml-cluster-mongodb-replicaset:27017")

	//session, err := mgo.Dial("mongo:27017")

	if err != nil {
		panic(err)
	}

	return session
}

// Close mongodb session
func Close() {
	collection.Database.Session.Close()
}

var collection = getSession().DB("ml").C("mutants")

// Ping se utiliza para determinar si el servicio esta disponible de forma trivial
func Ping(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "PONG")
}

// StatsMutants retorna estadisticas de Mutants
func StatsMutants(w http.ResponseWriter, r *http.Request) {

	//defer r.Body.Close()

	total, err := collection.Count()

	if err != nil {
		fmt.Println("Error al procesar Count general")
		panic(err)
	}
	fmt.Println("Total de elementos en la coleccion ", total)

	// Busco y cuento Mutants
	iCounterMutants, err := collection.Find(bson.M{"ismutant": true}).Count()
	if err != nil {
		panic(err)
	}
	fmt.Println("Total de Mutantes ", iCounterMutants)

	// Busco y cuento no Mutantes
	iCounterNotMutants, err := collection.Find(bson.M{"ismutant": false}).Count()
	if err != nil {
		panic(err)
	}
	fmt.Println("Total de NO Mutantes ", iCounterNotMutants)

	v := Stats{
		CountMutantDNA: strconv.Itoa(iCounterMutants),
		CountHumanDNA:  strconv.Itoa(iCounterNotMutants),
		Ratio:          strconv.FormatFloat(float64(iCounterMutants)/float64(iCounterNotMutants), 'f', -1, 32),
	}
	responseStats(w, 200, v)
}

// MutantCheck es
func MutantCheck(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)

	var adnChain AdnChain

	err := decoder.Decode(&adnChain)

	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// convierto los string en bytes, para poder trabajar de forma eficiente con los elementos
	var matrix [6][]byte

	for fila := 0; fila < len(adnChain.Chain); fila++ {
		matrix[fila] = []byte(adnChain.Chain[fila])
	}
	// Verifico si la cadena enviada esta formada por elementos correctos
	isValidChainElements := isValidChainElements(adnChain.Chain)

	// Si la cadena cuenta con elementos no validos
	if !isValidChainElements {
		fmt.Println("Cadena dna invalida")
		responseAdnCheck(w, 403, adnChain)
		return
	}
	isMutant := isMutant(matrix)

	if !isMutant {
		fmt.Println("El Humano no es mutante")

		// Inserto el registro dna en la base de datos
		err = collection.Insert(&AdnChainModel{Chain: adnChain.Chain, IsMutant: false, Timestamp: time.Now()})

		if err != nil {
			w.WriteHeader(500)
			return
		}

		responseAdnCheck(w, 403, adnChain)
	} else {
		fmt.Println("El Humano es Mutante")

		// Inserto el registro dna en la base de datos
		err = collection.Insert(&AdnChainModel{Chain: adnChain.Chain, IsMutant: true, Timestamp: time.Now()})

		if err != nil {
			w.WriteHeader(500)
			return
		}
		responseAdnCheck(w, 200, adnChain)
	}

}

func responseAdnCheck(w http.ResponseWriter, status int, results AdnChain) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(results)
}

func responseStats(w http.ResponseWriter, status int, results Stats) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(results)
}

// isMutant determina si una secuencia de dna es o no mutante
func isMutant(dna [6][]byte) bool {

	// para chequear que no existan elementos repetidos en dos hallazgos distintos

	var findingsMatrix = [6][]byte{
		{10, 10, 10, 10, 10, 10},
		{10, 10, 10, 10, 10, 10},
		{10, 10, 10, 10, 10, 10},
		{10, 10, 10, 10, 10, 10},
		{10, 10, 10, 10, 10, 10},
		{10, 10, 10, 10, 10, 10},
	}

	var isMut = false
	var iCountFindings int

	PrintMatrix(dna, " DNA Matrix")

	// Analizo horizontalmente
	iCountFindings = IsMutantByHorizontalAnalysis(dna, findingsMatrix)
	if iCountFindings > 1 {
		fmt.Println("Es Mutante, cantidad de hallazgos: ", iCountFindings)
		isMut = true
		//return isMut
	}

	// Analizo verticalmente
	iCountFindings += IsMutantByVerticalAnalysis(dna, findingsMatrix)
	if iCountFindings > 1 {
		fmt.Println("Es Mutante, cantidad de hallazgos: ", iCountFindings)
		isMut = true
		//return isMut
	}

	// Analizo oblicuamente
	iCountFindings += IsMutantByCentralObliqueAnalysis(dna, findingsMatrix)
	if iCountFindings > 1 {
		fmt.Println("Es Mutante, cantidad de hallazgos: ", iCountFindings)
		isMut = true
		//return isMut
	}

	PrintMatrix(dna, "transpuesta analisis oblicuo lateral")
	iCountFindings += IsMutantByLateralObliqueAnalysis(dna, findingsMatrix)
	if iCountFindings > 1 {
		fmt.Println("Es Mutante, cantidad de hallazgos: ", iCountFindings)
		isMut = true
		//return isMut
	}

	iCountFindings += IsMutantByExternalObliqueAnalysis(dna, findingsMatrix)
	if iCountFindings > 1 {
		fmt.Println("Es Mutante, cantidad de hallazgos: ", iCountFindings)
		isMut = true
		//return isMut
	}

	PrintMatrix(findingsMatrix, "findingsMatrix")

	return isMut
}

// IsMutantByExternalObliqueAnalysis analiza las coincidencias de forma oblicuia
func IsMutantByExternalObliqueAnalysis(dna [6][]byte, findingsMatrix [6][]byte) int {

	var anterior byte
	var iCounterSequence int
	var iCounterFindings int

	for fila, columna := 0, 2; columna < len(dna); fila, columna = fila+1, columna+1 {
		// Si es seleccionable
		if findingsMatrix[fila][columna] == 10 {

			if anterior == dna[fila][columna] {
				iCounterSequence++
				if iCounterSequence >= 3 {

					// Realice un hallazgo, lo marco en findingsMatrix
					for i, j := 0, 0; i <= 3; i, j = i+1, j+1 {
						findingsMatrix[fila-i][columna-i] = 20
					}
					iCounterFindings++
					break
				}
			} else {
				anterior = dna[fila][columna]
				iCounterSequence = 0
			}
		} else {
			anterior = 0
			iCounterSequence = 0
		}
	}

	iCounterSequence = 0
	anterior = 0

	for fila, columna := 2, 0; columna < len(dna)-2; fila, columna = fila+1, columna+1 {
		// Si es seleccionable
		if findingsMatrix[fila][columna] == 10 {

			if anterior == dna[fila][columna] {
				iCounterSequence++
				if iCounterSequence >= 3 {

					// Realice un hallazgo, lo marco en findingsMatrix
					for i, j := 0, 0; i <= 3; i, j = i+1, j+1 {
						findingsMatrix[fila-i][columna-i] = 20
					}
					iCounterFindings++
					break
				}
			} else {
				anterior = dna[fila][columna]
				iCounterSequence = 0
			}
		} else {
			anterior = 0
			iCounterSequence = 0
		}
	}

	return iCounterFindings
}

// IsMutantByLateralObliqueAnalysis analiza las coincidencias de forma oblicuia
func IsMutantByLateralObliqueAnalysis(dna [6][]byte, findingsMatrix [6][]byte) int {

	var anterior byte
	var iCounterSequence int
	var iCounterFindings int

	for fila, columna := 0, 1; columna < len(dna); fila, columna = fila+1, columna+1 {
		// Si es seleccionable
		if findingsMatrix[fila][columna] == 10 {
			if anterior == dna[fila][columna] {
				iCounterSequence++
				if iCounterSequence >= 3 {
					// Realice un hallazgo, lo marco en findingsMatrix
					for i, j := 0, 0; i <= 3; i, j = i+1, j+1 {
						findingsMatrix[fila-i][columna-i] = 20
					}
					iCounterFindings++
					break
				}
			} else {
				anterior = dna[fila][columna]
				iCounterSequence = 0
			}
		} else {
			anterior = 0
			iCounterSequence = 0
		}

	}

	iCounterSequence = 0
	anterior = 0

	for fila, columna := 1, 0; fila < len(dna); fila, columna = fila+1, columna+1 {
		// Si es seleccionable
		if findingsMatrix[fila][columna] == 10 {

			if anterior == dna[fila][columna] {
				iCounterSequence++
				if iCounterSequence >= 3 {
					// Realice un hallazgo, lo marco en findingsMatrix
					for i, j := 0, 0; i <= 3; i, j = i+1, j+1 {
						findingsMatrix[fila-i][columna-i] = 20
					}
					iCounterFindings++
					break
				}
			} else {
				anterior = dna[fila][columna]
				iCounterSequence = 0
			}
		} else {
			anterior = 0
			iCounterSequence = 0
		}

	}
	return iCounterFindings
}

// IsMutantByHorizontalAnalysis analiza horizontalmente las coincidencias
func IsMutantByHorizontalAnalysis(dna [6][]byte, findingsMatrix [6][]byte) int {

	var anterior byte
	var iCounterSequence int
	var iCounterFindings int

	for fila, h := range dna {
		anterior = 0
		for columna, _ := range h {
			// Solo si se trata de una celda elegible sigo adelante, esto evita elegir un celda mas de una vez
			if findingsMatrix[fila][columna] == 10 {
				// Determino si en la dimension horizonal se repite cualquier valor
				if anterior == dna[fila][columna] {
					iCounterSequence++

					if iCounterSequence >= 3 {
						// Realice un hallazgo, lo marco en findingsMatrix
						for i := 0; i <= 3; i++ {
							findingsMatrix[fila][columna-i] = 20
						}
						iCounterFindings++
						break
					}
				} else {
					anterior = dna[fila][columna]
					iCounterSequence = 0
				}
			} else {
				anterior = 0
				iCounterSequence = 0
			}
		}
	}

	return iCounterFindings
}

// IsMutantByVerticalAnalysis analiza horizontalmente las coincidencias
func IsMutantByVerticalAnalysis(dna [6][]byte, findingsMatrix [6][]byte) int {

	var anterior byte
	var iCounterSequence int
	var iCounterFindings int

	for columna, h := range dna {
		anterior = 0
		for fila, _ := range h {
			// Solo si se trata de una celda elegible sigo adelante, esto evita elegir un celda mas de una vez
			if findingsMatrix[fila][columna] == 10 {
				// Determino si en la dimension vertical se repite cualquier valor
				if anterior == dna[fila][columna] {
					iCounterSequence++

					if iCounterSequence >= 3 {
						// Realice un hallazgo vertical, lo marco en findingsMatrix
						for i := 0; i <= 3; i++ {
							findingsMatrix[fila-i][columna] = 20
						}
						iCounterFindings++
						break
					}
				} else {
					anterior = dna[fila][columna]
					iCounterSequence = 0
				}
			} else {
				anterior = 0
				iCounterSequence = 0
			}
		}
	}

	return iCounterFindings
}

// IsMutantByCentralObliqueAnalysis analiza las coincidencias de forma oblicuia
func IsMutantByCentralObliqueAnalysis(dna [6][]byte, findingsMatrix [6][]byte) int {

	var anterior byte
	var iCounterSequence int
	var iCounterFindings int

	for fila, columna := 0, 0; fila < len(dna); fila, columna = fila+1, columna+1 {

		// Solo si se trata de una celda elegible sigo adelante, esto evita elegir un celda mas de una vez
		if findingsMatrix[fila][columna] == 10 {
			if anterior == dna[fila][columna] {
				iCounterSequence++
				if iCounterSequence >= 3 {
					// Realice un hallazgo, lo marco en findingsMatrix
					for i, j := 0, 0; i <= 3; i, j = i+1, j+1 {
						findingsMatrix[fila-i][columna-i] = 20
					}
					iCounterFindings++
					break
				}
			} else {
				anterior = dna[fila][columna]
				iCounterSequence = 0
			}
		} else {
			anterior = 0
			iCounterSequence = 0
		}

	}
	return iCounterFindings
}

/*
// Transpose se utiliza para reutilizar el mismo algoritmo de busqueda horizontal
func Transpose(a [6][]byte, b [6][]byte) {
	n := len(a)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			a[i][j], a[j][i] = a[j][i], a[i][j]
			b[i][j], b[j][i] = b[j][i], b[i][j]
		}
	}
}
*/

// PrintMatrix inprime un cuadro de N x N con los elementos de la matrix
func PrintMatrix(a [6][]byte, detalle string) {

	println("Imprimiento Matrix", detalle)

	for _, h := range a {
		for _, cell := range h {
			fmt.Print(cell, " ")
		}
		fmt.Println()
	}
}

// isMutant determina si un dna es o no mutante
func isValidChainElements(dna []string) bool {

	var result = true

	// los elementos que se pasan como parametro deben estan en mayusculas, es case sensitive
	for _, h := range dna {
		for _, cell := range h {
			//fmt.Print(cell, " ")
			if !(cell == 65 || cell == 67 || cell == 71 || cell == 84) {
				result = false
			}
		}
		//fmt.Println()
	}
	return result
}
