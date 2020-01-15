package netconf

import (
	"sort"
)

// DBCType defines a DBC type.
type DBCType struct {
	Currency string // the DBC currency
	Amount   uint64 // the amount per DBC, last 8 digits are decimal places
}

func DBCTypeMapToSortedArray(m map[DBCType]bool) []DBCType {
	var dbcTypes []DBCType
	for t := range m {
		dbcTypes = append(dbcTypes, t)
	}
	sort.Slice(dbcTypes, func(i, j int) bool {
		if dbcTypes[i].Currency < dbcTypes[j].Currency ||
			(dbcTypes[i].Currency == dbcTypes[j].Currency &&
				dbcTypes[i].Amount < dbcTypes[j].Amount) {
			return true
		}
		return false
	})
	return dbcTypes
}
