package managed

import (
	"context"
	"fmt"

	"github.com/jaxsax/projects/tapeworm/botv2/internal/logging"
)

type Service interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}

type serviceWrapper struct {
	Service Service
	Name    string
}
type Manager struct {
	services []serviceWrapper

	started bool
}

func New() *Manager {
	return &Manager{}
}

func (m *Manager) Add(service Service, name string) {
	m.services = append(m.services, serviceWrapper{
		Service: service,
		Name:    name,
	})
}

type serviceStatus struct {
	err error
}

func (m *Manager) Start(ctx context.Context) error {
	if m.started {
		return fmt.Errorf("cannot start an already started manager")
	}

	for _, s := range m.services {
		serviceWrapper := s
		go func() {
			err := serviceWrapper.Service.Start(ctx)
			if err != nil {
				logging.FromContext(ctx).Error(err, "service service failed", "name", serviceWrapper.Name)
			}
		}()
	}

	m.started = true

	return nil
}

func (m *Manager) Stop(ctx context.Context) error {
	for i := len(m.services) - 1; i >= 0; i-- {
		serviceWrapper := m.services[i]

		status := make(chan error, 1)
		go func() {
			if err := serviceWrapper.Service.Stop(ctx); err != nil {
				status <- err
				return
			}

			status <- nil
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case stopErr := <-status:
			if stopErr != nil {
				logging.FromContext(ctx).Error(stopErr, "stop service failed")
				return stopErr
			}
		}
	}

	m.started = false

	return nil
}
