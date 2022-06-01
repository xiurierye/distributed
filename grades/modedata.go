package grades

func init() {

	students = []Student{
		{
			ID:        1,
			FirstName: "Nick",
			LastName:  "Carter",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 85,
				},
				{
					Title: "Quiz 2",
					Type:  GradeExam,
					Score: 85,
				},
				{
					Title: "Quiz 3",
					Type:  GradeTest,
					Score: 90,
				},
			},
		},
		{
			ID:        2,
			FirstName: "Nick",
			LastName:  "Carter",
			Grades: []Grade{
				{
					Title: "Quiz 1",
					Type:  GradeQuiz,
					Score: 88,
				},
				{
					Title: "Quiz 2",
					Type:  GradeExam,
					Score: 66,
				},
				{
					Title: "Quiz 3",
					Type:  GradeTest,
					Score: 99,
				},
			},
		},
	}
}
