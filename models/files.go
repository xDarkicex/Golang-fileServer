package models

import (
	"log"

	"github.com/xDarkicex/fileServer/datastore"
)

type File struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type files struct{}

var Files = new(files)

func (files) Create(name string) (*File, error) {
	file := &File{Name: name}
	var insertID int
	tx := datastore.Postgre.MustBegin()
	db := datastore.Postgre
	db.QueryRow("INSERT INTO files (name) VALUES ($1) returning id;", name).Scan(&insertID)
	err := tx.Commit()

	if err != nil {
		log.Fatal(err)
	}
	file.ID = insertID

	return file, err
}

func (files) Index() ([]*File, error) {
	files := []*File{}
	err := datastore.Postgre.Select(&files, "SELECT id, name FROM files")
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (files) ByName(name string) (*File, error) {
	file := File{}
	err := datastore.Postgre.Select(&file, "SELECT id, name FROM files WHERE name=$1;", name)
	if err != nil {
		return nil, err
	}
	return &file, nil
}
