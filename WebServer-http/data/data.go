package data

// Import library
import "strconv"

// Data struct
type Data struct {
	Student string
	Subject string
	Grade   string
}

// AllData struct
type AllData struct {
	Data []Data
}

/* Add Function
 * @param d Data
 */
func (dt *AllData) Add(d Data) {
	dt.Data = append(dt.Data, d)
}

/* StudentAVG Function
 * @param stu string
 * @return avg string
 */
func (dt *AllData) StudentAVG(stu string) string {
	var avg float64
	var count = 0

	for _, d := range dt.Data {
		if d.Student == stu {
			grade, _ := strconv.ParseFloat(d.Grade, 64)
			avg += grade
			count += 1
		}
	}

	avg /= float64(count)
	return strconv.FormatFloat(avg, 'f', 2, 64)
}

/* StudentAVG Function
 * @param sub string
 * @return avg string
 */
func (dt *AllData) SubjectAVG(sub string) string {
	var avg float64
	var count = 0

	for _, d := range dt.Data {
		if d.Subject == sub {
			grade, _ := strconv.ParseFloat(d.Grade, 64)
			avg += grade
			count += 1
		}
	}

	avg /= float64(count)
	return strconv.FormatFloat(avg, 'f', 2, 64)
}

/* GeneralAVG Function
 * @return string
 */
func (dt *AllData) GeneralAVG() string {
	var avg float64

	for _, d := range dt.Data {
		grade, _ := strconv.ParseFloat(d.Grade, 64)
		avg += grade
	}

	avg /= float64(len(dt.Data))
	return strconv.FormatFloat(avg, 'f', 2, 64)
}

/* String Function
 * @return string
 */
func (dt *AllData) String() string {
	var html string

	for _, d := range dt.Data {
		html +=
			"<tr>" +
				"<td>" + d.Student + "</td>" +
				"<td>" + d.Subject + "</td>" +
				"<td>" + d.Grade + "</td>" +
				"</tr>"
	}

	return html
}