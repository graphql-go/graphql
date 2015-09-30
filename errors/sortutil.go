package graphqlerrors

import "bytes"

type GQLFormattedErrorSlice []GraphQLFormattedError

func (errs GQLFormattedErrorSlice) Len() int {
	return len(errs)
}

func (errs GQLFormattedErrorSlice) Swap(i, j int) {
	errs[i], errs[j] = errs[j], errs[i]
}

func (errs GQLFormattedErrorSlice) Less(i, j int) bool {
	mCompare := bytes.Compare([]byte(errs[i].Message), []byte(errs[j].Message))
	lesserLine := errs[i].Locations[0].Line < errs[j].Locations[0].Line
	eqLine := errs[i].Locations[0].Line == errs[j].Locations[0].Line
	lesserColumn := errs[i].Locations[0].Column < errs[j].Locations[0].Column
	if mCompare < 0 {
		return true
	}
	if mCompare == 0 && lesserLine {
		return true
	}
	if mCompare == 0 && eqLine && lesserColumn {
		return true
	}
	return false
}
