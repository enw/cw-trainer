package midi

import (
	"fmt"
	"sync"
)

// Manager owns the port even while an asynchronous connection result is in
// flight. Closing the app cannot strand a port that finishes connecting later.
type Manager struct {
	mu     sync.Mutex
	driver Driver
	port   Port
	closed bool
}

func NewManager(driver Driver) *Manager       { return &Manager{driver: driver} }
func (m *Manager) Devices() ([]Device, error) { return m.driver.Devices() }
func (m *Manager) Connect(id int32) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.closed {
		return fmt.Errorf("MIDI manager closed")
	}
	if m.port != nil {
		if err := m.port.Close(); err != nil {
			m.port = nil
			return err
		}
		m.port = nil
	}
	port, err := m.driver.Connect(id)
	if err != nil {
		return err
	}
	m.port = port
	return nil
}
func (m *Manager) Poll() ([]Event, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.port == nil {
		return nil, fmt.Errorf("MIDI is not connected")
	}
	events, err := m.port.Poll()
	if err != nil {
		_ = m.port.Close()
		m.port = nil
	}
	return events, err
}
func (m *Manager) Disconnect() error { m.mu.Lock(); defer m.mu.Unlock(); return m.disconnect() }
func (m *Manager) disconnect() error {
	if m.port == nil {
		return nil
	}
	err := m.port.Close()
	m.port = nil
	return err
}
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closed = true
	return m.disconnect()
}
