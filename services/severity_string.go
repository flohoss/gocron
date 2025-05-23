// Code generated by "stringer -type=Severity"; DO NOT EDIT.

package services

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Debug-1]
	_ = x[Info-2]
	_ = x[Warning-3]
	_ = x[Error-4]
}

const _Severity_name = "DebugInfoWarningError"

var _Severity_index = [...]uint8{0, 5, 9, 16, 21}

func (i Severity) String() string {
	i -= 1
	if i < 0 || i >= Severity(len(_Severity_index)-1) {
		return "Severity(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _Severity_name[_Severity_index[i]:_Severity_index[i+1]]
}
