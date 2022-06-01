package grades

import (
	"fmt"
	"sync"
)

type Student struct {
	ID        int
	FirstName string
	LastName  string
	Grades    []Grade
}

func (s Student) Average() float32 {
	var result float32
	for _, grade := range s.Grades {
		result += grade.Score
	}

	return result / float32(len(s.Grades))
}

var (
	students      Students
	studentsMutex sync.Mutex
)

type Students []Student

func (ss Students) GetByID(id int) (*Student, error) {
	for _, s := range ss {
		if s.ID == id {
			return &s, nil
		}
	}
	return nil, fmt.Errorf(" Student with Id %d not found.", id)
}

type GradeType string

const (
	GradeQuiz = GradeType("Quiz")
	GradeTest = GradeType("Test")
	GradeExam = GradeType("Exam")
)

type Grade struct {
	Title string
	Type  GradeType
	Score float32
}
