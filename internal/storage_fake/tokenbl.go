package storage_fake

type tokenbl map[string]bool

func (s *storage_fake) AddTokenToBL(token string) {
	s.mx.Lock()
	defer s.mx.Unlock()
	s.tokenbl[token] = true

}

func (s *storage_fake) IsBlacklisted(token string) bool {
	s.mx.Lock()
	defer s.mx.Unlock()
	_, found := s.tokenbl[token]
	return found
}
