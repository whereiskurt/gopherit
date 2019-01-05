package adapter

import (
	"00-newapp-template/pkg"
	"00-newapp-template/pkg/acme"
	"path/filepath"
)

// Unmarshal holds the config - needed for Service.... TODO: Remove config and take Service
type Unmarshal struct {
	Config *pkg.Config
}

// NewUnmarshal calls the ACME Services and returns ACME JSONs to the adapter
func NewUnmarshal(config *pkg.Config) (u Unmarshal) {
	u.Config = config
	return
}

func (u *Unmarshal) service() (s acme.Service) {
	s = acme.NewService(u.Config.Client.BaseURL, u.Config.Client.SecretKey, u.Config.Client.AccessKey)
	if u.Config.Client.CacheResponse {
		serviceCacheFolder := filepath.Join(".", u.Config.Client.CacheFolder, "service/")
		s.EnableCache(serviceCacheFolder, u.Config.Client.CacheKey)
	}
	s.SetLogger(u.Config.Log)

	return
}

func (u *Unmarshal) gophers() (gg []acme.Gopher) {
	s := u.service()
	gg = s.GetGophers()
	return
}

func (u *Unmarshal) things(gopherID string) (tt []acme.Thing) {
	s := u.service()
	tt = s.GetThings(gopherID)
	return
}

// DeleteGopher returns all gophers are aren't deleted.
func (u *Unmarshal) deleteGopher(gopherID string) (gg []acme.Gopher) {
	s := u.service()
	gg = s.DeleteGopher(gopherID)
	return
}

// DeleteThings returns all things for a gopher that aren't deleted.
func (u *Unmarshal) deleteThing(gopherID string, thingID string) (tt []acme.Thing) {
	s := u.service()
	tt = s.DeleteThing(gopherID, thingID)
	return
}
