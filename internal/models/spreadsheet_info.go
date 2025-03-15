package models

type StudentInfo struct {
	Name string
	Mail string
	Link string
}

type TableInfo struct {
	Name     string
	Students []StudentInfo
}
