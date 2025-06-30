package lib

import "fmt"

// MapDataUpload contains information on zone maps
type MapDataUpload struct {
	ZoneID          int      `json:"ZoneID"`
	BuildingType    []int    `json:"BuildingType"`
	AvailableFood   []int    `json:"AvailableFood"`
	Reward          []int    `json:"Reward"`
	AvailableSilver []int    `json:"AvailableSilver"`
	Owners          []string `json:"Owners"`
	PublicFee       []int    `json:"PublicFee"`
	AssociateFee    []int    `json:"AssociateFee"`
	Coordinates     [][]int  `json:"Coordinates"`
	Durability      []int    `json:"Durability"`
	Permission      []int    `json:"Permission"`
}

func (m *MapDataUpload) StringArrays() [][]string {
	result := [][]string{}
	for i := range m.BuildingType {
		result = append(result, []string{
			fmt.Sprintf("%d", m.ZoneID),
			fmt.Sprintf("%d", m.BuildingType[i]),
			fmt.Sprintf("%d", m.AvailableFood[i]),
			fmt.Sprintf("%d", m.Reward[i]),
			fmt.Sprintf("%d", m.AvailableSilver[i]),
			m.Owners[i],
		})
	}

	return result
}
