package repository

import "github.com/includingByMeAndMyself/urlshortening/internal/model"

type Repository interface {
	Save(pair *model.URLPair)
	Get(id string) (*model.URLPair, bool)
}
