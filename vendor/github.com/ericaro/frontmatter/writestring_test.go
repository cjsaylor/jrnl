package frontmatter

import (
	"fmt"
	"strings"
	"testing"
)

type pseudo string

const (
	//v,V for not value Value (pointer or not)
	//x,X for not exported exported
	//p,P string or pseudo string
	xvp = iota
	xVp
	Xvp
	XVp
	xvP
	xVP
	XvP
	XVP
)

func nameOf(mode int) string {
	switch mode {
	case xvp:
		return "xvp"
	case xVp:
		return "xVp"
	case Xvp:
		return "Xvp"
	case XVp:
		return "XVp"
	case xvP:
		return "xvP"
	case xVP:
		return "xVP"
	case XvP:
		return "XvP"
	case XVP:
		return "XVP"
	}
	return "<>"
}

func xvpOf(v string) *string { return &v }
func xVpOf(v string) string  { return v }
func XvpOf(v string) *string { return &v }
func XVpOf(v string) string  { return v }
func xvPOf(v string) *pseudo { return (*pseudo)(&v) }
func xVPOf(v string) pseudo  { return pseudo(v) }
func XvPOf(v string) *pseudo { return (*pseudo)(&v) }
func XVPOf(v string) pseudo  { return pseudo(v) }

func Ofxvp(v *string) string { return *v }
func OfxVp(v string) string  { return v }
func OfXvp(v *string) string { return *v }
func OfXVp(v string) string  { return v }
func OfxvP(v *pseudo) string { return string(*v) }
func OfxVP(v pseudo) string  { return string(v) }
func OfXvP(v *pseudo) string { return string(*v) }
func OfXVP(v pseudo) string  { return string(v) }

var bench = tester{
	xvp: xvpOf("xvp"),
	xVp: xVpOf("xVp"),
	Xvp: XvpOf("Xvp"),
	XVp: XVpOf("XVp"),
	xvP: xvPOf("xvP"),
	xVP: xVPOf("xVP"),
	XvP: XvPOf("XvP"),
	XVP: XVPOf("XVP"),
}

//tester has all the cases: (string or *string) x (exported, non exported) x (convertible, non convertible)
type tester struct {
	xvp *string `xvp:"-"`
	xVp string  `xVp:"-"`
	Xvp *string `Xvp:"-"`
	XVp string  `XVp:"-"`
	xvP *pseudo `xvP:"-"`
	xVP pseudo  `xVP:"-"`
	XvP *pseudo `XvP:"-"`
	XVP pseudo  `XVP:"-"`
}

func (t *tester) Get(mode int) string {
	switch mode {
	case xvp:
		return Ofxvp(t.xvp)
	case xVp:
		return OfxVp(t.xVp)
	case Xvp:
		return OfXvp(t.Xvp)
	case XVp:
		return OfXVp(t.XVp)
	case xvP:
		return OfxvP(t.xvP)
	case xVP:
		return OfxVP(t.xVP)
	case XvP:
		return OfXvP(t.XvP)
	case XVP:
		return OfXVP(t.XVP)
	}
	return "<>"
}

func ExampleTester() {
	fmt.Println(bench.Get(xvp))
	fmt.Println(bench.Get(xVp))
	fmt.Println(bench.Get(Xvp))
	fmt.Println(bench.Get(XVp))
	fmt.Println(bench.Get(xvP))
	fmt.Println(bench.Get(xVP))
	fmt.Println(bench.Get(XvP))
	fmt.Println(bench.Get(XVP))
	//Output
	// xvp
	// xVp
	// Xvp
	// XVp
	// xvP
	// xVP
	// XvP
	// XVP
}

func TestWriteString(t *testing.T) {

	probeMode(t,
		xvp,
		xVp,
		Xvp,
		XVp,
		xvP,
		xVP,
		XvP,
		XVP,
	)
}
func probeMode(t *testing.T, modes ...int) {
	// we have
	for _, mode := range modes {
		n := nameOf(mode) //mode's name (xvp ...)
		//attempt a write
		err := WriteString(&bench, n, "-", n+"changed")

		//check the written result
		t.Logf("%s = %s", n, bench.Get(mode))
		written := (bench.Get(mode) == n+"changed")

		if strings.HasPrefix(n, "x") { //it was unexported, should have failed correctly
			if written {
				t.Errorf("error: written into %s", n)
			}
			if err != ErrUnexported {
				t.Errorf("writing to unexported field %s has the wrong error  %s", n, err.Error())
			}
		} else if !written {
			t.Errorf("error writing %s : got %s instead: err: %s", n, bench.Get(mode), err.Error())
		}
	}
}
func TestReadString(t *testing.T) {

}
func probeRead(t *testing.T, modes ...int) {
	for _, mode := range modes {

		txt, err := ReadString(bench, nameOf(mode), "-")
		if err != nil {
			t.Errorf("read string err: %s", err.Error())
		}
		if txt != nameOf(mode) {
			t.Errorf("read string invalid result: %q != %q", txt, nameOf(mode))
		}

	}
}
