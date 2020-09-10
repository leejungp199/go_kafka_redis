package module

const (
	// SP
	Ss = "sp"
	Rs = "tv"

	// suffix for , 
	suffix_p  = "pl"
	suffix_id = "il"

	// ID
	idpt_Sig      = "ss"

	// PL
	pl_sig      = "ss"

	// Message Typ
	MSG_C_REQ = "1a"
	MSG_C_RES = "1"
	MSG_M_REQ = "1"
	MSG_M_RES = "1s"
	MSG_D_REQ = "1"
	MSG_D_RES = "2"
)

func binarySearch(a []int, search int) (index int, searchCount int) {
	mid := len(a) / 2

	switch {
	case len(a) == 0:
		index = -1 // not found
	case a[mid] > search:
		index, searchCount = binarySearch(a[:mid], search)
	case a[mid] < search:
		index, searchCount = binarySearch(a[mid+1:], search)
		if index != -1 {
			index += mid + 1
		}
	default: // a[mid] == search
		index = mid // found
	}
	searchCount++
	return
}

/*
return index number if found. -1 if not found.
*/
func binarySearchFloat(a []float64, search float64) (index int, searchCount int) {
	mid := len(a) / 2

	switch {
	case len(a) == 0:
		index = -1 // not found
	case a[mid] > search:
		index, searchCount = binarySearchFloat(a[:mid], search)
	case a[mid] < search:
		index, searchCount = binarySearchFloat(a[mid+1:], search)
		if index != -1 {
			index += mid + 1
		}
	default: // a[mid] == search
		index = mid // found
	}
	searchCount++
	return
}
