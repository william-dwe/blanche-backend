package usecase

import (
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/dto"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/entity"
	"git-garena.com/sea-labs-id/batch-04/stage-02/blanche/blanche-be/repository"
)

type ExampleUsecase interface {
	ExampleProcess(dto.ExampleReqDTO) (*dto.ExampleResDTO, error)
	CachedExampleProcess(dto.ExampleReqDTO) (*dto.ExampleResDTO, error)
}

type ExampleUsecaseConfig struct {
	ExampleRepository repository.ExampleRepository
}

type exampleUsecaseImpl struct {
	exampleRepository repository.ExampleRepository
}

func NewExampleUsecase(c ExampleUsecaseConfig) ExampleUsecase {
	return &exampleUsecaseImpl{
		exampleRepository: c.ExampleRepository,
	}
}

func (u *exampleUsecaseImpl) CachedExampleProcess(input dto.ExampleReqDTO) (*dto.ExampleResDTO, error) {
	createExample := entity.Example{
		Qty:  10,
		Name: input.ExampleField,
	}

	err := u.exampleRepository.CachedStore(createExample)
	if err != nil {
		return nil, err
	}

	return &dto.ExampleResDTO{}, nil
}

func (u *exampleUsecaseImpl) ExampleProcess(input dto.ExampleReqDTO) (*dto.ExampleResDTO, error) {
	//do something
	createExample := entity.Example{
		//set value
		Qty:  10,
		Name: input.ExampleField,
	}

	err := u.exampleRepository.Store(createExample)
	if err != nil {
		return nil, err
	}

	return &dto.ExampleResDTO{}, nil
}
