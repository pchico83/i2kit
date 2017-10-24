package main

import (
	"github.com/spf13/viper"

	"net/http"

	"github.com/docker/notary"
	"github.com/docker/notary/client"
	"github.com/docker/notary/tuf/data"
)

const remoteConfigField = "api"

// RepoFactory takes a GUN and returns an initialized client.Repository, or an error.
type RepoFactory func(gun data.GUN) (client.Repository, error)

// ConfigureRepo takes in the configuration parameters and returns a repoFactory that can
// initialize new client.Repository objects with the correct upstreams and password
// retrieval mechanisms.
func ConfigureRepo(v *viper.Viper, retriever notary.PassRetriever, onlineOperation bool) RepoFactory {
	localRepo := func(gun data.GUN) (client.Repository, error) {
		var rt http.RoundTripper
		trustPin, err := getTrustPinning(v)
		if err != nil {
			return nil, err
		}
		if onlineOperation {
			rt, err = getTransport(v, gun, admin)
			if err != nil {
				return nil, err
			}
		}
		return client.NewFileCachedRepository(
			v.GetString("trust_dir"),
			gun,
			getRemoteTrustServer(v),
			rt,
			retriever,
			trustPin,
		)
	}

	return localRepo
}
