package main

import (
	"bytes"
	_ "embed"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

//go:embed odd.png
var logoData []byte

//go:embed a.png
var iconData []byte

//go:embed b.mp3
var musicData []byte

func rot13(s string) string {
	result := make([]rune, 0, len(s))
	for _, v := range s {
		switch {
		case v >= 'a' && v <= 'z':
			result = append(result, rune((int(v-'a')+13)%26)+'a')
		case v >= 'A' && v <= 'Z':
			result = append(result, rune((int(v-'A')+13)%26)+'A')
		case v >= '0' && v <= '9':
			result = append(result, rune((int(v-'0')+5)%10)+'0')
		case v == 'æ':
			result = append(result, 'ø')
		case v == 'ø':
			result = append(result, 'å')
		case v == 'å':
			result = append(result, 'æ')
		case v == 'Æ':
			result = append(result, 'Ø')
		case v == 'Ø':
			result = append(result, 'Å')
		case v == 'Å':
			result = append(result, 'Æ')
		default:
			result = append(result, v)
		}
	}
	return string(result)
}

type readSeekCloser struct{ *bytes.Reader }

func (r *readSeekCloser) Close() error { return nil }

func playMusic() {
	rsc := &readSeekCloser{bytes.NewReader(musicData)}
	streamer, format, err := mp3.Decode(rsc)
	if err != nil {
		return
	}
	if err := speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10)); err != nil {
		return
	}
	speaker.Play(beep.Loop(-1, streamer))
	select {} // block forever so the goroutine (and audio) keeps running
}

func main() {
	a := app.New()
	a.SetIcon(fyne.NewStaticResource("a.png", iconData))

	w := a.NewWindow("4F 44 44")
	w.Resize(fyne.NewSize(540, 420))
	w.CenterOnScreen()

	logo := canvas.NewImageFromResource(fyne.NewStaticResource("odd.png", logoData))
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(500, 148))

	keyLabel := widget.NewLabel("State your Key:")
	keyLabel.TextStyle = fyne.TextStyle{Bold: true}

	inputBox := widget.NewEntry()
	inputBox.SetPlaceHolder("Top Secret Message")

	cipherLabel := widget.NewLabel("Get your Code:")
	cipherLabel.TextStyle = fyne.TextStyle{Bold: true}

	outputBox := widget.NewEntry()
	outputBox.SetPlaceHolder("Code comes here...")
	outputBox.Disable()

	generateBtn := widget.NewButton("Generate", func() {
		outputBox.Enable()
		outputBox.SetText(rot13(inputBox.Text))
		outputBox.Disable()
	})
	clearBtn := widget.NewButton("Clear", func() {
		inputBox.SetText("")
		outputBox.Enable()
		outputBox.SetText("")
		outputBox.Disable()
	})

	btnRow := container.NewHBox(
		generateBtn,
		clearBtn,
		layout.NewSpacer(),
		widget.NewLabel("kthxbye :)"),
	)

	statusLabel := widget.NewLabel("Encoded with the top level cryptography algorithm ROT13.")
	statusLabel.Wrapping = fyne.TextWrapWord

	w.SetContent(container.NewPadded(container.NewVBox(
		logo,
		keyLabel,
		inputBox,
		cipherLabel,
		outputBox,
		btnRow,
		widget.NewSeparator(),
		statusLabel,
	)))

	go playMusic()

	w.ShowAndRun()
}
