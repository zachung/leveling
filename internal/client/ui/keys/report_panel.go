package keys

import (
	"github.com/gdamore/tcell/v2"
	"leveling/internal/client/service"
)

type ReportPanel struct {
	*T
}

func NewReportPanel(next Func) *ReportPanel {
	t := &T{next: next}
	i := &ReportPanel{T: t}
	t.Func = i

	return i
}

func (c ReportPanel) handleEvent(event *tcell.EventKey) *tcell.EventKey {
	if event.Key() == tcell.KeyF1 {
		service.UI().Report().Focus()
		return nil
	}
	return event
}
