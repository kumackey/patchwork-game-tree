package main

import (
	"fmt"
	"math"
	"sort"
)

type PatchShape [3][5]bool

type RemainingIncomeTimes int

const (
	One RemainingIncomeTimes = iota + 1
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
)

type SortType string

const (
	Evaluation SortType = "evaluation"
	Profit     SortType = "profit"
)

type Patch struct {
	Shape      PatchShape
	ButtonCost int
	TimeCost   int
	Buttons    int
	Size       int
	Name       string
}

func NewPatch(shape PatchShape, buttonCost, timeCost, buttonsEarned int) *Patch {
	patch := &Patch{
		Shape:      shape,
		ButtonCost: buttonCost,
		TimeCost:   timeCost,
		Buttons:    buttonsEarned,
	}
	patch.Size = patch.calculateSize()
	patch.Name = patch.generateName()
	return patch
}

func (p *Patch) calculateSize() int {
	sum := 0
	for _, row := range p.Shape {
		for _, cell := range row {
			if cell {
				sum++
			}
		}
	}
	return sum
}

func (p *Patch) generateName() string {
	return fmt.Sprintf("%d-%d(%d)+%d", p.ButtonCost, p.TimeCost, p.Size, p.Buttons)
}

func (p *Patch) Profit(remainingIncomeTimes RemainingIncomeTimes) float64 {
	return float64(p.Buttons)*float64(remainingIncomeTimes) + 2*float64(p.Size) - float64(p.ButtonCost)
}

func (p *Patch) ButtonsPerCost() float64 {
	return float64(p.Buttons) / float64(p.ButtonCost+p.TimeCost)
}

func (p *Patch) ProfitPerTime(remainingIncomeTimes RemainingIncomeTimes) float64 {
	return p.Profit(remainingIncomeTimes) / float64(p.TimeCost)
}

func (p *Patch) Evaluation(remainingIncomeTimes RemainingIncomeTimes) float64 {
	buttonCostWeighting := buttonCostWeightings[remainingIncomeTimes]
	return p.buttonPerCostZScore()*buttonCostWeighting + p.profitPerTimeZScore(remainingIncomeTimes)*(1-buttonCostWeighting)
}

func (p *Patch) buttonPerCostZScore() float64 {
	return p.calculateZScore(func(patch *Patch) float64 { return patch.ButtonsPerCost() }, Patches)
}

func (p *Patch) profitPerTimeZScore(remainingIncomeTimes RemainingIncomeTimes) float64 {
	return p.calculateZScore(func(patch *Patch) float64 { return patch.ProfitPerTime(remainingIncomeTimes) }, Patches)
}

func (p *Patch) calculateZScore(valueFunction func(*Patch) float64, patches []*Patch) float64 {
	mean := calculateAverage(valueFunction, patches)
	stdDev := calculateStandardDeviation(valueFunction, patches)
	return (valueFunction(p) - mean) / stdDev
}

func calculateAverage(valueFunction func(*Patch) float64, patches []*Patch) float64 {
	sum := 0.0
	for _, patch := range patches {
		sum += valueFunction(patch)
	}
	return sum / float64(len(patches))
}

func calculateStandardDeviation(valueFunction func(*Patch) float64, patches []*Patch) float64 {
	mean := calculateAverage(valueFunction, patches)
	variance := 0.0
	for _, patch := range patches {
		value := valueFunction(patch)
		variance += math.Pow(value-mean, 2)
	}
	return math.Sqrt(variance / float64(len(patches)))
}

var buttonCostWeightings = map[RemainingIncomeTimes]float64{
	One:   0,
	Two:   0.0625,
	Three: 0.125,
	Four:  0.25,
	Five:  0.5,
	Six:   0.75,
	Seven: 0.875,
	Eight: 0.9375,
	Nine:  0.96875,
}

func sortPatchesByEvaluation(remainingIncomeTimes RemainingIncomeTimes, patches []*Patch) []*Patch {
	sortedPatches := make([]*Patch, len(patches))
	copy(sortedPatches, patches)
	sort.Slice(sortedPatches, func(i, j int) bool {
		return sortedPatches[i].Evaluation(remainingIncomeTimes) > sortedPatches[j].Evaluation(remainingIncomeTimes)
	})
	return sortedPatches
}

func sortPatchesByProfit(remainingIncomeTimes RemainingIncomeTimes, patches []*Patch) []*Patch {
	sortedPatches := make([]*Patch, len(patches))
	copy(sortedPatches, patches)
	sort.Slice(sortedPatches, func(i, j int) bool {
		return sortedPatches[i].Profit(remainingIncomeTimes) > sortedPatches[j].Profit(remainingIncomeTimes)
	})
	return sortedPatches
}

var Patches = []*Patch{
	NewPatch([3][5]bool{
		{true, false, false, false, false},
		{true, true, true, true, false},
		{false, false, false, true, false},
	}, 1, 2, 0),
	NewPatch([3][5]bool{
		{false, true, false, false, false},
		{true, true, true, true, false},
		{false, false, true, false, false},
	}, 2, 1, 0),
	NewPatch([3][5]bool{
		{true, true, true, true, true},
		{false, false, false, false, false},
		{false, false, false, false, false},
	}, 7, 1, 1),
	NewPatch([3][5]bool{
		{true, false, false, false, false},
		{true, true, true, true, false},
		{true, false, false, false, false},
	}, 7, 2, 2),
	NewPatch([3][5]bool{
		{false, true, true, false, false},
		{true, true, true, true, false},
		{false, true, true, false, false},
	}, 5, 3, 1),
	NewPatch([3][5]bool{
		{false, true, false, false, false},
		{true, true, true, true, false},
		{false, true, false, false, false},
	}, 0, 3, 1),
	NewPatch([3][5]bool{
		{true, true, false, false, false},
		{false, true, true, false, false},
		{false, false, false, false, false},
	}, 3, 2, 1),
	NewPatch([3][5]bool{
		{false, false, true, false, false},
		{true, true, true, true, true},
		{false, false, true, false, false},
	}, 1, 4, 1),
	NewPatch([3][5]bool{
		{true, true, true, false, false},
		{true, false, true, false, false},
		{false, false, false, false, false},
	}, 1, 2, 0),
	NewPatch([3][5]bool{
		{true, true, true, false, false},
		{true, false, false, false, false},
		{false, false, false, false, false},
	}, 4, 2, 1),
	NewPatch([3][5]bool{
		{true, true, true, false, false},
		{false, false, true, true, false},
		{false, false, false, false, false},
	}, 2, 3, 1),
	NewPatch([3][5]bool{
		{true, true, true, false, false},
		{false, true, false, false, false},
		{true, true, true, false, false},
	}, 2, 3, 0),
	NewPatch([3][5]bool{
		{true, true, true, false, false},
		{true, true, false, false, false},
		{false, false, false, false, false},
	}, 2, 2, 0),
	NewPatch([3][5]bool{
		{true, true, true, false, false},
		{false, true, true, true, false},
		{false, false, false, false, false},
	}, 4, 2, 0),
	NewPatch([3][5]bool{
		{false, true, false, false, false},
		{true, true, true, false, false},
		{false, true, false, false, false},
	}, 5, 4, 2),
	NewPatch([3][5]bool{
		{true, true, true, true, false},
		{false, true, true, false, false},
		{false, false, false, false, false},
	}, 7, 4, 2),
	NewPatch([3][5]bool{
		{true, true, false, false, false},
		{false, true, true, false, false},
		{false, false, true, false, false},
	}, 10, 4, 3),
	NewPatch([3][5]bool{
		{true, true, true, true, false},
		{true, true, false, false, false},
		{false, false, false, false, false},
	}, 10, 5, 3),
	NewPatch([3][5]bool{
		{true, true, true, true, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
	}, 3, 3, 1),
	NewPatch([3][5]bool{
		{true, true, true, true, false},
		{true, false, false, false, false},
		{false, false, false, false, false},
	}, 10, 3, 2),
	NewPatch([3][5]bool{
		{true, true, true, true, false},
		{true, false, false, true, false},
		{false, false, false, false, false},
	}, 1, 5, 1),
	NewPatch([3][5]bool{
		{true, false, true, false, false},
		{true, true, true, false, false},
		{false, true, false, false, false},
	}, 3, 6, 2),
	NewPatch([3][5]bool{
		{true, true, false, false, false},
		{true, true, true, false, false},
		{false, false, true, false, false},
	}, 8, 6, 3),
	NewPatch([3][5]bool{
		{true, true, true, false, false},
		{false, true, false, false, false},
		{false, true, false, false, false},
	}, 5, 5, 2),
	NewPatch([3][5]bool{
		{true, true, true, true, false},
		{false, true, false, false, false},
		{false, false, false, false, false},
	}, 3, 4, 1),
	NewPatch([3][5]bool{
		{true, true, true, false, false},
		{false, true, false, false, false},
		{false, false, false, false, false},
	}, 2, 2, 0),
	NewPatch([3][5]bool{
		{true, true, false, false, false},
		{true, false, false, false, false},
		{false, false, false, false, false},
	}, 3, 1, 0),
	NewPatch([3][5]bool{
		{true, true, false, false, false},
		{false, true, true, false, false},
		{false, false, false, false, false},
	}, 7, 6, 3),
	NewPatch([3][5]bool{
		{true, true, false, false, false},
		{true, true, false, false, false},
		{false, false, false, false, false},
	}, 6, 5, 2),
	NewPatch([3][5]bool{
		{true, true, true, false, false},
		{true, false, false, false, false},
		{false, false, false, false, false},
	}, 4, 6, 2),
	NewPatch([3][5]bool{
		{true, true, true, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
	}, 2, 2, 0),
	NewPatch([3][5]bool{
		{true, true, false, false, false},
		{false, false, false, false, false},
		{false, false, false, false, false},
	}, 2, 1, 0),
	NewPatch([3][5]bool{
		{true, true, false, false, false},
		{true, false, false, false, false},
		{false, false, false, false, false},
	}, 1, 3, 0),
}
