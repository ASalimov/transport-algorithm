package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"runtime/debug"

	"github.com/olekukonko/tablewriter"
)

type link struct {
	i, j int
}

// FactsWrapper is an struct
type FactsWrapper struct {
	Facts
	decision Decision
}

const (
	empty = math.MaxFloat64
)

// ErrorDLength is an error that indicates an incorrect length of demands. It should be equal to the length of Costs
var ErrorDLength = errors.New("wrong length of demands")

// ErrorSLength is an error that indicates an incorrect length of supplies. It should be equal to the length of Costs
var ErrorSLength = errors.New("wrong length of supplies")

// NewFactsWrapper is an function
func NewFactsWrapper(facts Facts) (*FactsWrapper, error) {
	if len(facts.Demands) != len(facts.Costs) {
		return nil, ErrorDLength
	}
	if len(facts.Supplies) != len(facts.Costs[0]) {
		return nil, ErrorSLength
	}
	return &FactsWrapper{Facts: facts}, nil
}

// Find is an function for finding optimal decision
func (t *FactsWrapper) Find() error {
	t.balance()
	t.MinimumRate()
	fmt.Println("First solution by the minimum tariff method:")
	fmt.Println(t)
	fmt.Println()
	for {
		if i, j, ok := t.IsOptimal(); !ok {
			t.BetterOptimal(i, j)
		} else {
			break
		}
	}
	for i := range t.decision.Volume {
		for j := range t.decision.Volume[i] {
			if t.decision.Volume[i][j] == empty {
				t.decision.Volume[i][j] = 0
			}
		}
	}
	return nil
}

func (t *FactsWrapper) balance() {
	delta := sum(t.Supplies) - sum(t.Demands)
	if delta > 0 {
		t.Demands = append(t.Demands, delta)
		cost := make([]float64, len(t.Supplies))
		for i := range cost {
			cost[i] = 0
		}
		t.Costs = append(t.Costs, cost)

	}
	if delta < 0 {
		t.Supplies = append(t.Supplies, -delta)
		for i := range t.Costs {
			t.Costs[i] = append(t.Costs[i], 0)
		}
	}
}

// MinimumRate is a function that provide first decision using minimum rate algorithm
func (t *FactsWrapper) MinimumRate() [][]float64 {

	lD := len(t.Demands)
	lS := len(t.Supplies)
	resp := makeTable(lD, lS)
	remD, remS := make([]float64, lD), make([]float64, lS)
	copy(remD, t.Demands)
	copy(remS, t.Supplies)
	costs := make([][]float64, lD)
	for i := range t.Costs {
		costs[i] = make([]float64, len(t.Costs[i]))
		copy(costs[i], t.Costs[i])
	}

	basises := 0
	for true {
		iD, iS, price := min(costs)
		if price == empty {
			break
		}
		basises++
		switch true {
		case remD[iD] < remS[iS]:
			resp[iD][iS] = remD[iD]
			remS[iS] -= remD[iD]
			remD[iD] = 0
			for k := 0; k < lS; k++ {
				costs[iD][k] = empty
			}
		case remD[iD] > remS[iS]:
			resp[iD][iS] = remS[iS]
			remD[iD] -= remS[iS]
			remS[iS] = 0
			for k := 0; k < lD; k++ {
				costs[k][iS] = empty
			}
		case remD[iD] == remS[iS]:
			resp[iD][iS] = remD[iD]
			remD[iD] = 0
			remS[iS] = 0
			for k := 0; k < lD; k++ {
				costs[k][iS] = empty
			}
			for k := 0; k < lS; k++ {
				costs[iD][k] = empty
			}
		}

	}
	t.decision = Decision{
		Volume:     resp,
		PotentialV: makeSlice(len(t.Demands)),
		PotentialU: makeSlice(len(t.Supplies)),
	}
	t.GetTotal()
	return resp
}

func (t FactsWrapper) String() string {
	defer func() {
		if r := recover(); r != nil {
			Error.Println("error: ", r, "\nstacktrace from panic: \n"+string(debug.Stack()))
		}
	}()
	header := make([]string, len(t.Supplies)+1)
	header[0] = "Demands\\Supplies"
	bulk := [][]string{}
	for i, supply := range t.Supplies {
		header[i+1] = fmt.Sprintf("Supply %d [%.0f]", i+1, supply)
	}
	for i, demand := range t.Demands {
		bulkItem := make([]string, len(t.Supplies)+1)
		bulkItem[0] = fmt.Sprintf("Demand %d [%.0f]", i+1, demand)
		for j := range t.Supplies {
			dC := t.Costs[i][j]
			if t.Costs[i][j] == empty {
				dC = 0
			}
			if t.decision.Volume[i][j] == empty {
				bulkItem[j+1] = fmt.Sprintf("%2.0f []", dC)
			} else {
				bulkItem[j+1] = fmt.Sprintf("%2.0f [%.0f]", dC, t.decision.Volume[i][j])
			}
		}
		bulk = append(bulk, bulkItem)
	}
	buf := bytes.NewBufferString("")
	tw := tablewriter.NewWriter(buf)
	tw.SetHeader(header)
	tw.SetAutoFormatHeaders(false)
	footer := make([]string, len(t.Supplies)+1)
	footer[len(t.Supplies)-1] = "Total"
	footer[len(t.Supplies)] = fmt.Sprintf("%.2f", t.decision.Total)
	tw.SetFooter(footer)
	tw.SetBorder(false)
	tw.AppendBulk(bulk)

	tw.Render()

	st := buf.String()
	return st
}

// GetTotal is a function provide total amount of the current decision
func (t *FactsWrapper) GetTotal() (total float64) {
	for i, row := range t.Costs {
		for j, price := range row {
			if t.decision.Volume[i][j] != empty {
				total += t.decision.Volume[i][j] * price
			}
		}
	}
	t.decision.Total = total
	return
}

func (t *FactsWrapper) preventDegeneracy() {
	basises := 0
	for i := range t.decision.Volume {
		for j := range t.decision.Volume[i] {
			if t.decision.Volume[i][j] == 0 {
				t.decision.Volume[i][j] = empty
			} else if t.decision.Volume[i][j] != empty {
				basises++
			}
		}
	}
	needToAdd := len(t.Demands) + len(t.Supplies) - 1 - basises
	for needToAdd > 0 {
		randI := rand.Intn(len(t.Demands))
		randJ := rand.Intn(len(t.Supplies))
		if t.decision.Volume[randI][randJ] == empty {
			t.decision.Volume[randI][randJ] = 0
			needToAdd--
		}

	}
}

// IsOptimal is a function, returns true if decision is optimal
func (t *FactsWrapper) IsOptimal() (int, int, bool) {
	t.preventDegeneracy()
	t.decision.PotentialU = makeSlice(len(t.Supplies))
	t.decision.PotentialV = makeSlice(len(t.Demands))
	i := 0
	for {

		if err := t.findPotentials(); err == nil {
			break
		}
		t.preventDegeneracy()
		i++
		if i > 10000 {
			log.Fatal("failed to findPotentials")
		}
	}
	PV := t.decision.PotentialV
	PU := t.decision.PotentialU
	notOptimal := false
	for i := range t.Costs {
		for j := range t.Costs[i] {
			if PV[i]+PU[j] > t.Costs[i][j] && t.decision.Volume[i][j] == empty && t.Costs[i][j] != empty {
				notOptimal = true
				Info.Printf("not omptimal, cause: u+v = %.0f  costs[%d][%d]=%.0f\n", PV[i]+PU[j], i, j, t.Costs[i][j])
				_, err := t.goByChainHorizontaly(i, j, i, j)
				if err != nil {
					Info.Printf("wrong chain for costs[%d][%d]=%.0f", i, j, t.Costs[i][j])
					continue
				}
				return i, j, false
			}
		}
	}
	return 0, 0, !notOptimal
}

func (t *FactsWrapper) findPotentials() error {

	V := t.decision.Volume
	PV := t.decision.PotentialV
	PU := t.decision.PotentialU
	reset(PV)
	reset(PU)

	iMax, jMax, _ := max(V)
	price := t.Costs[iMax][jMax]
	PV[iMax] = price - 1
	PU[jMax] = 1
	found := 0
	for {
		found1 := found
		for i := 0; i < len(t.Demands); i++ {
			for j := 0; j < len(t.Supplies); j++ {
				cost := t.Costs[i][j]
				if V[i][j] != empty {
					if PV[i] != empty && PU[j] == empty {
						PU[j] = cost - PV[i]
						found1++
					} else if PU[j] != empty && PV[i] == empty {
						PV[i] = cost - PU[j]
						found1++
					}
				}
			}
		}
		if found1 >= len(t.Demands)+len(t.Supplies)-2 {
			return nil
		}
		if found1 == found {
			return fmt.Errorf("potential values not found")
		}
		found = found1
	}
}

// BetterOptimal is a function that finds a more optimal solution than the previous one
func (t *FactsWrapper) BetterOptimal(pivotI, pivotJ int) {
	repeat := 100
	var chain []link
	var err error
	for {
		chain, err = t.goByChainHorizontaly(pivotI, pivotJ, pivotI, pivotJ)
		if err == nil {
			break
		}
		t.preventDegeneracy()
		repeat--
		if repeat == 0 {
			log.Fatal(err)
		}

	}

	min := math.MaxFloat64
	for i, link := range chain {
		if i%2 == 0 && t.decision.Volume[link.i][link.j] < min {
			min = t.decision.Volume[link.i][link.j]
		}
	}

	t.decision.Volume[pivotI][pivotJ] = min
	for i, link := range chain {
		min1 := min
		if i%2 == 0 {
			min1 = -min
		}
		if t.decision.Volume[link.i][link.j] += min1; t.decision.Volume[link.i][link.j] == 0 {
			t.decision.Volume[link.i][link.j] = empty
		}
	}
	t.GetTotal()

}

func (t *FactsWrapper) goByChainHorizontaly(pivotI, pivotJ, iC, jC int) ([]link, error) {

	for j := 0; j < len(t.Supplies); j++ {
		// horizontaly
		if jC != j && t.decision.Volume[iC][j] != empty {
			if pivotI == iC && pivotJ == j {
				return []link{}, nil
			}
			if chain1, err := t.goByChainVertically(pivotI, pivotJ, iC, j); err == nil {
				chain := []link{link{iC, j}}
				chain = append(chain, chain1...)
				return chain, nil
			}
		}
	}
	return nil, fmt.Errorf("chain not found")
}

func (t *FactsWrapper) goByChainVertically(pivotI, pivotJ, iC, jC int) ([]link, error) {
	for i := 0; i < len(t.Demands); i++ {
		// vertically
		if pivotI == i && pivotJ == jC {
			return []link{}, nil
		}
		if iC != i && t.decision.Volume[i][jC] != empty {

			if chain1, err := t.goByChainHorizontaly(pivotI, pivotJ, i, jC); err == nil {
				chain := []link{link{i, jC}}
				chain = append(chain, chain1...)
				return chain, nil
			}
		}
	}
	return nil, fmt.Errorf("chain not found")
}
