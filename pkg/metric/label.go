// @Description 常用标签

package metric

const (
	labErr = "0"
	labOK  = "1"
)

// RetLabel RetLabel
func RetLabel(err error) string {
	if err == nil {
		return labOK
	}
	return labErr
}
