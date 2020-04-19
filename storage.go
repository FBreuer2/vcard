package vcard

type VCardStore interface {
	GetEntity(id string) (*VCardEntity, error)
	AddEntity(newEntity *VCardEntity) error
}

type MemoryDB struct {
	entities map[string]*VCardEntity
}
