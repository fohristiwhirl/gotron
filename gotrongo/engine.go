package gotrongo

import (
	"fmt"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"strings"
	"sync"
)

var eng engine
var sync_chan chan bool = make(chan bool)

func init() {
	eng.sprites = make(map[string]string)
	eng.sounds = make(map[string]string)
	eng.keyboard = make(map[string]bool)
}

type click struct {
	Button          int
	X               int
	Y               int
}

type engine struct {

	mutex           sync.Mutex

	// The following are written once only...

	started         bool
	fps             float64
	width           int
	height          int

	// The following are written several times at the beginning, then only read from...

	sprites         map[string]string       // filename -> JS varname
	sounds          map[string]string       // filename -> JS varname

	// Written often...

	keyboard        map[string]bool
	click           *click
}

func RegisterSprite(filename string) {

	eng.mutex.Lock()
	defer eng.mutex.Unlock()

	if eng.started {
		panic("RegisterSprite(): already started")
	}

	sprite_id := len(eng.sprites)

	eng.sprites[filename] = fmt.Sprintf("sprite%d", sprite_id)
	fmt.Printf("r\x1e%s\x1fsprite%d\n", filename, sprite_id)
}

func RegisterSound(filename string) {

	eng.mutex.Lock()
	defer eng.mutex.Unlock()

	if eng.started {
		panic("RegisterSound(): already started")
	}

	sound_id := len(eng.sounds)

	eng.sounds[filename] = fmt.Sprintf("sound%d", sound_id)
	fmt.Printf("s\x1e%s\x1fsound%d\n", filename, sound_id)
}

func Start(fps float64) (width int, height int) {

	eng.mutex.Lock()            // Really just for the .started var
	defer eng.mutex.Unlock()

	config := load_config("gotron.cfg")

	eng.width = config.Width
	eng.height = config.Height

	if eng.started {
		panic("wsengine.Start(): already started")
	}

	eng.started = true
	eng.fps = fps

	go stdin_reader()

	return config.Width, config.Height
}

func KeyDown(key string) bool {
	return _keydown(key, false)
}

func KeyDownClear(key string) bool {       // Clears the key after (sets it to false)
	return _keydown(key, true)
}

func _keydown(key string, clear bool) bool {
	eng.mutex.Lock()
	defer eng.mutex.Unlock()

	ret := eng.keyboard[key]

	if clear {
		eng.keyboard[key] = false
	}

	return ret
}

func GetWidthHeight() (int, int) {

	eng.mutex.Lock()
	defer eng.mutex.Unlock()

	return eng.width, eng.height
}

func GetWidthHeightFloats() (float64, float64) {

	eng.mutex.Lock()
	defer eng.mutex.Unlock()

	return float64(eng.width), float64(eng.height)
}

func GetClick() *click {

	// Return mouse click then delete it from memory.

	eng.mutex.Lock()
	defer eng.mutex.Unlock()

	ret := eng.click
	eng.click = nil
	return ret
}

func SendDebug(format_string string, args ...interface{}) {

	msg := fmt.Sprintf(format_string, args...)

	msg = strings.Replace(msg, "\x1e", " ", -1)       // Replace meaningful characters in our protocol
	msg = strings.Replace(msg, "\x1f", " ", -1)
	msg = strings.Replace(msg, "\n", " ", -1)

	final := "d\x1e" + template.HTMLEscapeString(msg) + "\n"

	eng.mutex.Lock()
	defer eng.mutex.Unlock()

	fmt.Printf(final)
}

func Sync() {		// Waits until the next frontend draw.
	<- sync_chan
}

type Config struct {
	Executable     string      `json:"executable"`
	Width          int         `json:"width"`
	Height         int         `json:"height"`
}

func load_config(filename string) Config {

	file, err := ioutil.ReadFile(filename)
	if err != nil {
		panic("Couldn't load config file")
	}

	var result Config
	err = json.Unmarshal(file, &result)

	if err != nil {
		panic("Couldn't parse config file (bad JSON?)")
	}

	return result
}
