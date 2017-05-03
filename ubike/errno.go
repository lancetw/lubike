package ubike

// StationErrno : errno for Station
type StationErrno int

// StationErrno code
const (
	Full                         StationErrno = 1
	Ok                           StationErrno = 0
	InvalidLatLng                StationErrno = -1
	GivenLocationNotInTaipeiCity StationErrno = -2
	SystemError                  StationErrno = -3
)
