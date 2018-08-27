package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os/exec"
	"os/user"
	"path"
	"strings"

	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgb/xtest"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/ewmh"
	"github.com/BurntSushi/xgbutil/icccm"
	"github.com/BurntSushi/xgbutil/keybind"
	"github.com/BurntSushi/xgbutil/mousebind"
	"github.com/BurntSushi/xgbutil/xevent"
	shellwords "github.com/mattn/go-shellwords"
)

const (
	versionStr = "Gosture v1.0"
	authorStr  = "By AyuanX, 22-Aug-2018"
	cfgFile    = ".Gosture.cfg"
)

var (
	x      *xgbutil.XUtil
	mg     *mgT
	cfg    *cfgT
	gesMap map[string]func()
	dirMap = [8]byte{'4', '7', '8', '9', '6', '3', '2', '1'}
	rdy    = false
	ena    = false
)

type mgT struct {
	button xproto.Button
	mods   uint16
	x, y   int
	v, c   int
	dir    []byte
}

func (m *mgT) begin(ex, ey int) {
	m.x, m.y, m.c = ex, ey, 0
	m.dir = nil
}

func (m *mgT) step(ex, ey int) {
	a := angle(m.x, m.y, ex, ey)
	m.x, m.y = ex, ey
	if m.c == 0 {
		m.v, m.c = a, 1
	} else if m.v != a {
		m.c--
	} else if m.c < 4 {
		m.c++
	} else {
		m.c = 0
		if len(m.dir) == 0 || m.dir[len(m.dir)-1] != dirMap[a] {
			m.dir = append(m.dir, dirMap[a])
		}
	}
}

func (m *mgT) end(ex, ey int) {
	if m.dir != nil {
		if f, ok := gesMap[string(m.dir)]; ok {
			f()
		}
	} else if mg.mods == 0 {
		mousebind.UngrabPointer(x)
		mousebind.Ungrab(x, x.RootWin(), mg.mods, mg.button)
		// TODO: proper device id
		xtest.FakeInput(x.Conn(), xproto.ButtonPress, byte(mg.button), xproto.TimeCurrentTime, x.RootWin(), int16(ex), int16(ey), 0)
		xtest.FakeInput(x.Conn(), xproto.ButtonRelease, byte(mg.button), xproto.TimeCurrentTime, x.RootWin(), int16(ex), int16(ey), 0)
		mousebind.Grab(x, x.RootWin(), mg.mods, mg.button, false)
	}
}

func angle(x1, y1, x2, y2 int) int {
	a := math.Atan2(float64(y2-y1), float64(x2-x1))
	if a >= math.Pi*7/8 {
		a -= math.Pi * 7 / 8
	} else {
		a += math.Pi * 9 / 8
	}
	return int(a * 4 / math.Pi)
}

type cfgT struct {
	EnMouse bool       `json:"mouse-gesture-enable"`
	TrMouse string     `json:"mouse-gesture-trigger"`
	GeList  []*gesActT `json:"gesture-list"`
	KeList  []*keyActT `json:"hotkey-list"`
}
type gesActT struct {
	Ges string `json:"gesture"`
	Act string `json:"action"`
}
type keyActT struct {
	Key string `json:"hotkey"`
	Act string `json:"action"`
}

func (cfg *cfgT) readCFG() error {
	// TODO: validate cfg
	user, err := user.Current()
	if err != nil {
		return err
	}
	buf, err := ioutil.ReadFile(path.Join(user.HomeDir, cfgFile))
	if err != nil {
		return err
	}
	if err = json.Unmarshal(buf, cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *cfgT) applyCFG() error {
	// TODO: more error handling
	if cfg.EnMouse == true {
		var err error
		if mg.mods, mg.button, err = mousebind.ParseString(x, cfg.TrMouse); err != nil {
			// Welp, to draw text using ximg is a pain!!! And I am lazy
			exec.Command("zenity", "--error", "--no-markup", "--no-wrap", "--title=Gosture", "--text="+fmt.Sprintf("Failed to parse mouse gesture trigger!\n\n%v", err)).Run()
			log.Fatalf("Filed to parse mouse gesture trigger! %v", err)
		}
		gesMap = make(map[string]func(), len(cfg.GeList))
		for _, ge := range cfg.GeList {
			strs := strings.Split(ge.Act, ",")
			switch strings.ToLower(strs[0]) {
			case "cmd":
				cmd := buildCMD(strs[1:])
				gesMap[ge.Ges] = func() { runCMD(cmd) }
			case "key":
				keys := buildKey(strs[1:])
				gesMap[ge.Ges] = func() { runKey(keys) }
			case "minwin":
				gesMap[ge.Ges] = minWin
			case "maxwin":
				gesMap[ge.Ges] = maxWin
			case "closewin":
				gesMap[ge.Ges] = closeWin
			}
		}
		mousebind.Drag(x, x.RootWin(), x.RootWin(), cfg.TrMouse, true,
			func(x *xgbutil.XUtil, rx, ry, ex, ey int) (bool, xproto.Cursor) {
				mg.begin(ex, ey)
				return true, 0
			},
			func(x *xgbutil.XUtil, rx, ry, ex, ey int) {
				mg.step(ex, ey)
			},
			func(x *xgbutil.XUtil, rx, ry, ex, ey int) {
				mg.end(ex, ey)
			})
	}
	// TODO: "Super + <key>" does not always work, could be Start Menu conflict?
	for _, ke := range cfg.KeList {
		strs := strings.Split(ke.Act, ",")
		switch strings.ToLower(strs[0]) {
		case "cmd":
			cmd := buildCMD(strs[1:])
			keybind.KeyPressFun(func(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
				runCMD(cmd)
			}).Connect(x, x.RootWin(), ke.Key, true)
		case "key":
			keys := buildKey(strs[1:])
			keybind.KeyPressFun(func(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
				runKey(keys)
			}).Connect(x, x.RootWin(), ke.Key, true)
		case "minwin":
			keybind.KeyPressFun(func(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
				minWin()
			}).Connect(x, x.RootWin(), ke.Key, true)
		case "maxwin":
			keybind.KeyPressFun(func(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
				maxWin()
			}).Connect(x, x.RootWin(), ke.Key, true)
		case "closewin":
			keybind.KeyPressFun(func(x *xgbutil.XUtil, e xevent.KeyPressEvent) {
				closeWin()
			}).Connect(x, x.RootWin(), ke.Key, true)
		}
	}
	return nil
}

func buildCMD(strs []string) *exec.Cmd {
	shellwords.ParseBacktick = true
	shellwords.ParseEnv = true
	if args, err := shellwords.Parse(strs[0]); err != nil {
		exec.Command("zenity", "--error", "--no-markup", "--no-wrap", "--title=Gosture", "--text="+fmt.Sprintf("Failed to parse cmd!\n\n%v", err)).Run()
		log.Fatalf("Failed to parse cmd! %v", err)
		return nil
	} else {
		cmd := exec.Command(args[0], args[1:]...)
		if len(strs) > 1 {
			cmd.Dir = strs[1]
		}
		return cmd
	}
}

func runCMD(cmd *exec.Cmd) {
	// exec.Cmd does not support reuse
	cmdl := *cmd
	go cmdl.Run()
}

func buildKey(strs []string) []xproto.Keycode {
	keys := make([]xproto.Keycode, len(strs))
	for i, s := range strs {
		codes := keybind.StrToKeycodes(x, s)
		if len(codes) == 0 {
			exec.Command("zenity", "--error", "--no-markup", "--no-wrap", "--title=Gosture", "--text="+fmt.Sprintf("Failed to parse keycode name: %q!", s)).Run()
			log.Fatalf("Failed to parse keycode name: %q!", s)
		}
		keys[i] = codes[0]
	}
	return keys
}

func runKey(keys []xproto.Keycode) {
	for _, k := range keys {
		// TODO: proper device id
		xtest.FakeInput(x.Conn(), xproto.KeyPress, byte(k), xproto.TimeCurrentTime, x.RootWin(), 0, 0, 0)
		defer xtest.FakeInput(x.Conn(), xproto.KeyRelease, byte(k), xproto.TimeCurrentTime, x.RootWin(), 0, 0, 0)
	}
}

func minWin() {
	//TODO: get window under cursor instead of active window
	if w, err := ewmh.ActiveWindowGet(x); err == nil && w != 0 {
		ewmh.ClientEvent(x, w, "WM_CHANGE_STATE", icccm.StateIconic)
	}
}

func maxWin() {
	//TODO: get window under cursor instead of active window
	if w, err := ewmh.ActiveWindowGet(x); err == nil {
		ewmh.WmStateReqExtra(x, w, ewmh.StateToggle, "_NET_WM_STATE_MAXIMIZED_HORZ", "_NET_WM_STATE_MAXIMIZED_VERT", 2)
	}
}

func closeWin() {
	//TODO: get window under cursor instead of active window
	if w, err := ewmh.ActiveWindowGet(x); err == nil {
		ewmh.CloseWindow(x, w)
	}
}

func gosture() {
	var err error
	if x, err = xgbutil.NewConn(); err != nil {
		log.Panic(err)
	}
	if err = xtest.Init(x.Conn()); err != nil {
		log.Panic(err)
	}
	keybind.Initialize(x)
	mousebind.Initialize(x)
	defer enable(false)

	mg = new(mgT)
	cfg = new(cfgT)

	reload()
	xevent.Main(x)
}

func reload() {
	enable(false)
	if err := cfg.readCFG(); err != nil {
		rdy = false
		go exec.Command("zenity", "--error", "--no-markup", "--no-wrap", "--title=Gosture", "--text="+fmt.Sprintf("Failed to parse cfg file: %q!\n\n%v", cfgFile, err)).Run()
		//log.Errorf("Failed to parse cfg file: %v!\n%v", cfgFile, err)
	} else {
		rdy = true
		enable(true)
	}
	onToggle()
}

func enable(en bool) {
	if !rdy || en == ena {
		return
	}
	if en == false {
		mousebind.Detach(x, x.RootWin())
		keybind.Detach(x, x.RootWin())
		xevent.Detach(x, x.RootWin())
	} else {
		cfg.applyCFG()
	}
	ena = en
}
