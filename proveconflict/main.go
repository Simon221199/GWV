package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"sort"
	"strings"
	"time"
)

type set map[string]bool

func (set1 set) Join(set2 set) (set, bool) {

	newSet := set{}

	for key := range set1 {
		newSet[key] = true
	}

	elemAdded := false

	for key := range set2 {
		if _, ok := newSet[key]; !ok {
			newSet[key] = true
			elemAdded = true
		}
	}

	return newSet, elemAdded
}

func (set1 set) toString() string {

	parts := make([]string, len(set1))

	inx := 0
	for atom := range set1 {
		parts[inx] = atom
		inx++
	}

	sort.Strings(parts)

	return strings.Join(parts, ", ")
}

type proofPart map[string]set

type hornClause struct {
	h   string
	set set
}

func proveBottomUp(knowledgeBase []hornClause, assumables set) {

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(knowledgeBase), func(i, j int) { knowledgeBase[i], knowledgeBase[j] = knowledgeBase[j], knowledgeBase[i] })

	c := make(proofPart)

	for assumable := range assumables {

		c[assumable] = set{
			assumable: true,
		}

	}

	conflictCount := 0

	for _, clause := range knowledgeBase {

		if clause.h == "false" {
			conflictCount++
		}
	}

	foundConflicts := 0

	falses := make([]set, 0)

while:
	for {

	outerloop:
		for _, clause := range knowledgeBase {

			fmt.Printf("Checking <%s, {%s}>\n", clause.h, clause.set.toString())

			// if _, ok := c[ clause.h ]; ok {
			// 	continue
			// }

			A := set{}

			for b := range clause.set {

				if _, ok := c[b]; !ok {
					continue outerloop
				}

				if val, ok := c[b]; ok {
					A, _ = A.Join(val)
				}

				if assumables[b] {
					A[b] = true
				}
			}

			if val, ok := c[clause.h]; ok {
				a, _ := A.Join(val)
				A = a
			}

			fmt.Printf("Added <%s, {%s}>\n", clause.h, A.toString())

			if clause.h == "false" {
				foundConflicts++
				// break while

				falses = append(falses, A)
			} else {
				c[clause.h] = A
			}

			if foundConflicts >= conflictCount {
				break while
			}
		}
	}

	bytes, _ := json.MarshalIndent(falses, "", "    ")
	fmt.Println(string(bytes))
}

func main() {

	fmt.Println("hallo")

	// knowledgeBase := []hornClause{
	//
	// 	//
	// 	// observations
	// 	//
	// 	{"up_s1", set{}},
	// 	{"up_s2", set{}},
	// 	{"up_s3", set{}},
	// 	{"dark_l1", set{}},
	// 	{"dark_l2", set{}},
	//
	// 	//
	// 	// knowledge
	// 	//
	// 	{"light_l1", set{}},
	// 	{"light_l2", set{}},
	// 	{"live_outside", set{}},
	// 	{"live_l1", set{"live_w0": true}},
	// 	{"live_w0", set{"live_w1": true, "up_s2": true, "ok_s2": true}},
	// 	{"live_w0", set{"live_w2": true, "down_s2": true, "ok_s2": true}},
	// 	{"live_w1", set{"live_w3": true, "up_s1": true, "ok_s1": true}},
	// 	{"live_w2", set{"live_w3": true, "down_s1": true, "ok_s1": true}},
	// 	{"live_l2", set{"live_w4": true}},
	// 	{"live_w4", set{"live_w3": true, "up_s3": true, "ok_s3": true}},
	// 	{"live_p1", set{"live_w3": true}},
	// 	{"live_w3", set{"live_w5": true, "ok_cb1": true}},
	// 	{"live_p2", set{"live_w6": true}},
	// 	{"live_w6", set{"live_w5": true, "ok_cb2": true}},
	// 	{"live_w5", set{"live_outside": true}},
	// 	{"lit_l1", set{"light_l1": true, "live_l1": true, "ok_l1": true}},
	// 	{"lit_l2", set{"light_l2": true, "live_l2": true, "ok_l2": true}},
	// 	{"false", set{"dark_l1": true, "lit_l1": true}},
	// 	{"false", set{"dark_l2": true, "lit_l2": true}},
	// }
	//
	// assumables := set{
	// 	"ok_cb1": true,
	// 	"ok_cb2": true,
	// 	"ok_s1":  true,
	// 	"ok_s2":  true,
	// 	"ok_s3":  true,
	// 	"ok_l1":  true,
	// 	"ok_l2":  true,
	// }

	// knowledgeBase := []hornClause{
	//
	// 	//
	// 	// observations
	// 	//
	//
	// 	{"gardener_dirty", set{}},
	// 	{"gardener_not_dirty", set{}},
	// 	{"gardener_working", set{}},
	// 	{"butler_working", set{}},
	//
	// 	//
	// 	// rules
	// 	//
	//
	// 	{"gardener_dirty", set{"gardener_working": true}},
	// 	{"butler_dirty", set{"butler_working": true}},
	// 	{"false", set{"gardener_dirty": true, "gardener_not_dirty": true}},
	// 	{"false", set{"butler_dirty": true, "butler_not_dirty": true}},
	// }
	//
	// assumables := set{
	// 	"gardener_working": true,
	// 	"butler_working":   true,
	// }

	// knowledgeBase := []hornClause{
	// 	{"false", set{"a": true}},
	// 	{"a", set{"b": true, "c": true}},
	// 	{"b", set{"d": true}},
	// 	{"b", set{"e": true}},
	// 	{"c", set{"f": true}},
	// 	{"c", set{"g": true}},
	// 	{"e", set{"h": true, "w": true}},
	// 	{"e", set{"g": true}},
	// 	{"w", set{"f": true}},
	// }
	//
	// assumables := set{
	// 	"d": true,
	// 	"f": true,
	// 	"g": true,
	// 	"h": true,
	// }

	knowledgeBase := []hornClause{

		// {"starter_noise", set{}},
		// {"fuel_pump_noise", set{}},
		// {"engine_noise", set{}},
		{"no_starter_noise", set{}},
		{"no_fuel_pump_noise", set{}},
		{"no_engine_noise", set{}},

		//
		// rules
		//

		{"battery", set{"battery_ok": true}},
		{"ignition_key", set{"ignition_key_ok": true, "battery": true}},
		{"efr", set{"efr_ok": true, "ignition_key": true, "battery": true}},
		{"starter", set{"starter_ok": true, "ignition_key": true}},
		{"starter_noise", set{"starter": true}},
		{"fuel_tank", set{"fuel_tank_ok": true}},
		{"fuel_pump", set{"fuel_pump_ok": true, "efr": true, "fuel_tank": true}},
		{"fuel_pump_noise", set{"fuel_pump": true}},
		{"filter", set{"filter_ok": true, "fuel_pump": true}},
		{"engine", set{"engine_ok": true, "starter": true, "filter": true}},
		{"engine_noise", set{"engine": true}},

		{"false", set{"starter_noise": true, "no_starter_noise": true}},
		{"false", set{"fuel_pump_noise": true, "no_fuel_pump_noise": true}},
		{"false", set{"engine_noise": true, "no_engine_noise": true}},
	}

	assumables := set{
		"battery_ok":      true,
		"ignition_key_ok": true,
		"efr_ok":          true,
		"starter_ok":      true,
		"engine_ok":       true,
		"filter_ok":       true,
		"fuel_pump_ok":    true,
		"fuel_tank_ok":    true,
	}

	proveBottomUp(knowledgeBase, assumables)
}
