//go:build windows

package e2e

import (
	"sync/atomic"
	"testing"

	"github.com/egoist/mygo"
	win "github.com/egoist/mygo/internal/windows"
)

// TestPanelAltF4 asks to close a panel on Alt+F4, which has no window menu
// (TestPanelStyles), with a material behind its page.
func TestPanelAltF4(t *testing.T) {
	panel := newWindow(t, mygo.WindowOptions{Panel: true, Frameless: true, Transparent: true, Vibrancy: mygo.VibrancyPopover, Width: 400, Height: 200})
	panel.Page().LoadHTML(`<body style="background:transparent">panel</body>`, "")
	eventually(t, "the panel focused", panel.IsFocused)
	var closes atomic.Int32
	panel.OnClose(func(e *mygo.CloseEvent) { closes.Add(1); e.PreventDefault() })
	mygo.RunOnMain(func() { win.TestPressKeys(0x12, 0x73) }) // VK_MENU, VK_F4
	eventually(t, "Alt+F4 asks to close", func() bool { return closes.Load() > 0 })
}

// A panel the app maximizes fills the work area of its display, as other
// windows do, rather than the whole display: Windows gives a window
// without a maximize box the whole monitor, and the panel, topmost, would
// hide the taskbar. Its invisible borders overhang the work area alike on
// every side.
func TestPanelMaximizedInWorkArea(t *testing.T) {
	panel := newWindow(t, mygo.WindowOptions{Panel: true, Frameless: true, Transparent: true, Vibrancy: mygo.VibrancyAcrylic, X: 100, Y: 100, Width: 400, Height: 200})
	panel.Page().LoadHTML(`<body style="background:transparent">panel</body>`, "")
	eventually(t, "the panel shown", panel.IsVisible)
	b := panel.Bounds()
	d := mygo.Screen.DisplayNearestPoint(mygo.Point{X: b.X + b.Width/2, Y: b.Y + b.Height/2})
	if d.WorkArea == d.Bounds {
		t.Skip("the display has no work area apart from its bounds: no taskbar to keep clear")
	}
	panel.Maximize()
	eventually(t, "the panel maximized", panel.IsMaximized)
	wa := d.WorkArea
	near := func(a, b int) bool { return a-b <= 1 && b-a <= 1 } // DIPs, rounded
	eventually(t, "the panel in the work area", func() bool {
		b = panel.Bounds()
		left, top := wa.X-b.X, wa.Y-b.Y
		right, bottom := b.X+b.Width-(wa.X+wa.Width), b.Y+b.Height-(wa.Y+wa.Height)
		return near(left, right) && near(top, bottom) && near(left, top)
	})
	if c := panel.ContentBounds(); c != wa {
		t.Errorf("content bounds = %+v, want the work area %+v", c, wa)
	}
}
