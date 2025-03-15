package service

import "EduCommentSync/internal/processor"

func (s *Service) process() error {
	rawComments, err := s.repo.GetRawComments()
	if err != nil {
		return err
	}

	teachers, err := s.repo.GetTeachers()
	if err != nil {
		return err
	}

	comments := processor.ProcessComments(rawComments, teachers)

	err = s.repo.AddComments(comments)
	return err
}
