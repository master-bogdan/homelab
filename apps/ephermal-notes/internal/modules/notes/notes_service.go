package notes

import "github.com/master-bogdan/ephermal-notes/internal/infra/db/redis"

type NotesService interface {
	CreateNote(note *memory_db.NotesModel) (*memory_db.NotesModel, error)
	GetNote(ID string) (*memory_db.NotesModel, error)
	DeleteNote(ID string) error
}

type notesService struct {
	repo memory_db.NotesRepository
}

func NewNotesService(repository memory_db.NotesRepository) NotesService {
	return &notesService{
		repo: repository,
	}
}

func (s *notesService) CreateNote(note *memory_db.NotesModel) (*memory_db.NotesModel, error) {
	val, err := s.repo.Create(note)
	if err != nil {
		return nil, err
	}

	return val, err
}

func (s *notesService) GetNote(ID string) (*memory_db.NotesModel, error) {
	note, err := s.repo.GetById(ID)
	if err != nil {
		return nil, err
	}

	return note, err
}

func (s *notesService) DeleteNote(ID string) error {
	err := s.repo.Delete(ID)
	if err != nil {
		return err
	}

	return nil
}
