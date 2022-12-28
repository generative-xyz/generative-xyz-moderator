package response

//user Devices server has prefix with "8"
const (
	//Success has prefix with 8"0"
	Success                       = 1
	

	//Error has prefix with 8"1"
	Error = -1
)

//Message
var ResponseMessage = map[int]struct {
	Code    int
	Message string
}{
	Success:                       {Success, "Success"},
	Error: {Error, "Failed."},
}
