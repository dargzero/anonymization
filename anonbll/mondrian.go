package anonbll

import (
	"fmt"
	"github.com/dargzero/anonymization/anondb"
	"github.com/dargzero/anonymization/anonmodel"
	"sort"
)

type mondrianDimension interface {
	initialize(anonCollectionName string, fieldName string)
	getInitialBoundaries() anonmodel.Boundary
	getDimensionForStatistics(interface{}, bool) mondrianDimension
	prepare(partition anonmodel.Partition, count int)
	getNormalizedRange() float64
	tryGetAllowableCut(int, anonmodel.Partition, int) (bool, anonmodel.Partition, anonmodel.Partition, error)
}

type mondrianDimensions []mondrianDimension

func (d mondrianDimensions) Len() int      { return len(d) }
func (d mondrianDimensions) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d mondrianDimensions) Less(i, j int) bool {
	return d[i].getNormalizedRange() > d[j].getNormalizedRange()
}

type mondrian struct {
	dataset            *anonmodel.Dataset
	anonCollectionName string
	fields             []anonmodel.FieldAnonymizationInfo
	dimensions         map[string]mondrianDimension
}

func (m *mondrian) initialize(dataset *anonmodel.Dataset, anonCollectionName string, fields []anonmodel.FieldAnonymizationInfo) {
	m.dataset = dataset
	m.anonCollectionName = anonCollectionName
	m.fields = fields
}

func (m *mondrian) anonymize() error {
	if err := anondb.CreateIndices(m.anonCollectionName, m.fields); err != nil {

		return err
	}
	defer anondb.DropIndices(m.anonCollectionName, m.fields)

	if err := m.createDimensions(); err != nil {

		return err
	}

	partition := m.getInitialPartition()
	return m.anonymizePartition(partition, true)
}

func (m *mondrian) createDimensions() error {
	m.dimensions = make(map[string]mondrianDimension)

	for _, field := range m.fields {
		var dimension mondrianDimension
		switch field.Type {
		case "numeric":
			dimension = &NumericDimension{}
		case "prefix":
			dimension = &prefixDimension{}
		case "coords":
			dimension = &gpsDimension{}
		default:
			return fmt.Errorf("Not supprted field type: %v", field.Type)
		}

		dimension.initialize(m.anonCollectionName, field.Name)
		m.dimensions[field.Name] = dimension
	}

	return nil
}

func (m *mondrian) getInitialPartition() anonmodel.Partition {
	result := make(map[string]anonmodel.Boundary)
	for fieldName, dimension := range m.dimensions {
		result[fieldName] = dimension.getInitialBoundaries()
	}
	return result
}

func (m *mondrian) anonymizePartition(partition anonmodel.Partition, firstRun bool) error {
	count, statistics, err := anondb.GetDimensionStatistics(m.anonCollectionName, partition)
	if err != nil {

		return err
	}

	if count < m.dataset.Settings.K {
		return fmt.Errorf("Not enough documents in db to satisfy %v-anonymity (have only %v documents)", m.dataset.Settings.K, count)
	}

	var currentDimensions mondrianDimensions
	for fieldName, stat := range statistics {
		dimension := m.dimensions[fieldName].getDimensionForStatistics(stat, firstRun)
		dimension.prepare(partition, count)
		currentDimensions = append(currentDimensions, dimension)
	}

	if count < m.dataset.Settings.K*2 {
		return m.generalize(partition)
	}

	allowable, left, right, err := m.chooseDimension(currentDimensions, partition, count)
	if err != nil {

		return err
	}
	if !allowable {
		return m.generalize(partition)
	}

	if err := m.anonymizePartition(left, false); err != nil {

		return err
	}
	return m.anonymizePartition(right, false)
}

func (m *mondrian) chooseDimension(currentDimensions mondrianDimensions, partition anonmodel.Partition, count int) (bool, anonmodel.Partition, anonmodel.Partition, error) {
	sort.Sort(currentDimensions)
	for _, dimension := range currentDimensions {
		if allowable, left, right, err := dimension.tryGetAllowableCut(m.dataset.Settings.K, partition, count); err != nil || allowable {
			return allowable, left, right, err
		}
	}
	return false, nil, nil, nil
}

func (m *mondrian) generalize(partition anonmodel.Partition) error {
	return anondb.Generalize(m.anonCollectionName, partition)
}
