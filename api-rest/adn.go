package main

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// AdnChain se utiliza para guardar los valores del adn enviados a la API
type AdnChain struct {
	Chain []string `json:"dna"`
}

// AdnChainModel se utiliza para persistir los valores del adn enviados a la API
type AdnChainModel struct {
	ID        bson.ObjectId `json:"id" bson:"_id,omitempty"`
	Chain     []string      `json:"dna"`
	IsMutant  bool          `json:"ismutant"`
	Timestamp time.Time     `json:"time"`
}

// Stats es la estructura para gestionar el json  {“count_mutant_dna”:40, “count_human_dna”:100: “ratio”:0.4}
type Stats struct {
	CountMutantDNA string `json:"count_mutant_dna"`
	CountHumanDNA  string `json:"count_human_dna"`
	Ratio          string `json:"ratio"`
}
