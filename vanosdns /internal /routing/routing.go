package routing

import (
	"errors"
	"sync"
	"time"

	"github.com/miekg/dns"

	"vanosdns/internal/antisibil"
	"vanosdns/internal/forwarder"
	"vanosdns/internal/mesh"
)

// Router est le cerveau de décision DNS
type Router struct {
	mesh       *mesh.Mesh
	forwarder  *forwarder.Forwarder
	antisibil  *antisibil.Engine

	// Hedged requests
	hedgeDelay time.Duration
}

// NewRouter initialise le router
func NewRouter(
	meshNet *mesh.Mesh,
	fwd *forwarder.Forwarder,
	as *antisibil.Engine,
) *Router {
	return &Router{
		mesh:       meshNet,
		forwarder: fwd,
		antisibil: as,
		hedgeDelay: 12 * time.Millisecond, // micro-delay v1
	}
}

// Resolve décide comment résoudre la requête DNS
func (r *Router) Resolve(req *dns.Msg) (*dns.Msg, error) {
	// Sélection des meilleurs noeuds mesh
	nodes := r.mesh.BestNodes(2)
	nodes = r.antisibil.Filter(nodes)

	// Si aucun noeud mesh fiable → upstream direct
	if len(nodes) == 0 {
		return r.resolveUpstream(req)
	}

	// Hedged requests : mesh d’abord, upstream en backup
	respCh := make(chan *dns.Msg, 2)
	errCh := make(chan error, 2)

	var once sync.Once
	sendResp := func(resp *dns.Msg, err error) {
		once.Do(func() {
			if err != nil {
				errCh <- err
				return
			}
			respCh <- resp
		})
	}

	// Mesh request
	go func() {
		resp, err := r.mesh.Resolve(req, nodes)
		sendResp(resp, err)
	}()

	// Hedged upstream (micro-delay)
	go func() {
		time.Sleep(r.hedgeDelay)
		resp, err := r.resolveUpstream(req)
		sendResp(resp, err)
	}()

	// Premier résultat gagne
	select {
	case resp := <-respCh:
		return resp, nil
	case err := <-errCh:
		return nil, err
	case <-time.After(1500 * time.Millisecond):
		return nil, errors.New("DNS resolution timeout")
	}
}

// resolver upstream appelle les DNS upstream classiques
func (r *Router) resolveUpstream(req *dns.Msg) (*dns.Msg, error) {
	resp, _, err := r.forwarder.Resolve(req)
	return resp, err
}
