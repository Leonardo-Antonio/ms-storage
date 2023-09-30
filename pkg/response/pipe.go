package response

type Response struct {
	Success   bool   `json:"success" xml:"success"`
	Data      any    `json:"data" xml:"data"`
	ItemFound bool   `json:"itemFound" xml:"itemFound"`
	TimeStamp uint64 `json:"timestamp" xml:"timestamp"`
}
