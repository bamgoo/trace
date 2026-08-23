package trace

import (
	"github.com/infrago/base"
	"github.com/infrago/infra"
)

func (m *Module) Ready() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.started && len(m.instances) > 0
}

func (m *Module) Health() infra.ModuleHealth {
	m.mutex.RLock()
	started := m.started
	connections := len(m.instances)
	m.mutex.RUnlock()
	return infra.NewModuleHealth("trace", started && connections > 0, nil, base.Map{"connections": connections})
}

func (m *Module) Stats() infra.ModuleStats {
	ready := m.Ready()
	return infra.NewModuleStats("trace", ready, m.Snapshot())
}
