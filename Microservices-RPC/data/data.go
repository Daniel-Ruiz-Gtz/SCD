package data

import "errors"

// Server struct.
type Server struct {
	Subjects map[string]map[string]float64
	Students map[string]map[string]float64
}

// Data struct.
type Data struct {
	Student string
	Subject string
	Grade   float64
}

/* Function AddRegister
	@param data Data
	@param reply *string
*/
func (s *Server) AddRegister(data Data, reply *string) error {
	// New register.
	if _, ok := s.Students[data.Student]; !ok {
		// Create grade map.
		grade := make(map[string]float64)
		// Assignment of grade for subject.
		grade[data.Subject] = data.Grade
		// Assignment of grade for students.
		s.Students[data.Student] = grade
	} else {
		// The record already exists, if student and subject exists.
		if _, ok := s.Students[data.Student][data.Subject]; ok {
			return errors.New("\nA grade already exists for this student")
		}
		// Add student to existing subject.
		s.Students[data.Student][data.Subject] = data.Grade
	}

	// Other register.
	if _, ok := s.Subjects[data.Subject]; !ok {
		// Create grade map.
		grade := make(map[string]float64)
		// Assignment of grade for student.
		grade[data.Student] = data.Grade
		// Assignment of grade for subjects.
		s.Subjects[data.Subject] = grade
	} else {
		// Add subject to an existing student.
		s.Subjects[data.Subject][data.Student] = data.Grade
	}

	// Submit results.
	*reply = "Successful registration"

	// Return error.
	return nil
}

/* Function GetStudentAverage
	@param student string
	@param s Server
*/
func GetStudentAverage(student string, s *Server) float64 {
	// Variable average.
	var average = float64(0)

	// Add grades.
	for subject := range s.Students[student] {
		average += s.Students[student][subject]
	}

	// Division grades.
	average /= float64(len(s.Students[student]))

	// return average.
	return average
}

/* Function StudentAverage
	@param data Data
	@param reply *float64
*/
func (s *Server) StudentAverage(data Data, reply *float64) error {
	// Get student average.
	var average = GetStudentAverage(data.Student, s)

	// Submit results.
	*reply = average

	// Return error.
	return nil
}

/* Function GeneralAverage
	@param data Data
	@param reply *float64
*/
func (s *Server) GeneralAverage(_ Data, reply *float64) error {
	// Variable average.
	var average = float64(0)

	// Add all grades students.
	for student := range s.Students {
		average += GetStudentAverage(student, s)
	}

	// Division all grades students.
	average /= float64(len(s.Students))

	// Submit results.
	*reply = average

	// Return error.
	return nil
}

/* Function SubjectAverage
	@param data Data
	@param reply *float64
*/
func (s *Server) SubjectAverage(data Data, reply *float64) error {
	// Variable average.
	var average = float64(0)

	// Add grades for subject.
	for student := range s.Subjects[data.Subject] {
		average += s.Subjects[data.Subject][student]
	}

	// Division grades for subject.
	average /= float64(len(s.Subjects[data.Subject]))

	// Submit results.
	*reply = average

	// Return error.
	return nil
}
