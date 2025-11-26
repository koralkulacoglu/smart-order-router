package ui

import (
	"fmt"
	"sync"
	"time"

	"github.com/koralkulacoglu/smart-order-router/internal/models"
)

type Dashboard struct {
	mu        sync.Mutex
	logs      []string
	portfolio *models.Portfolio
}

func NewDashboard(p *models.Portfolio) *Dashboard {
	return &Dashboard{
		logs:      make([]string, 0),
		portfolio: p,
	}
}

func (d *Dashboard) Log(format string, a ...any) {
	d.mu.Lock()
	defer d.mu.Unlock()

	msg := fmt.Sprintf(format, a...)
	timestamp := time.Now().Format("15:04:05.000")
	fullMsg := fmt.Sprintf("[%s] %s", timestamp, msg)

	d.logs = append(d.logs, fullMsg)

	if len(d.logs) > 15 {
		d.logs = d.logs[1:]
	}
}

func (d *Dashboard) Run(stopChan <-chan bool) {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-stopChan:
			return
		case <-ticker.C:
			d.render()
		}
	}
}

func (d *Dashboard) render() {
	d.mu.Lock()
	defer d.mu.Unlock()

	fmt.Print("\033[H\033[2J")

	status := d.portfolio.GetStatus()
	border := "══════════════════════════════════════════════════════════════════════"

	fmt.Println()
	fmt.Printf("╔%s╗\n", border)
	fmt.Printf("║ %-67s ║\n", status)
	fmt.Printf("╚%s╝\n", border)
	fmt.Println()

	fmt.Println("--- Live Event Log ---")
	for _, log := range d.logs {
		fmt.Println(log)
	}
}
