package data

type Employee struct {
	Name   string
	Office string
}

func GetData() []Employee {
	emplist := []Employee{
		Employee{
			Name:   "dumbledore",
			Office: "headmaster",
		},
		Employee{
			Name:   "harry",
			Office: "auror",
		},
		Employee{
			Name:   "hermionie",
			Office: "minister",
		},
		Employee{
			Name:   "ginny",
			Office: "quidditch",
		},
		Employee{
			Name:   "ron",
			Office: "auror",
		},
		Employee{
			Name:   "snape",
			Office: "professor",
		},
	}
	return emplist
}
