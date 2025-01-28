package bitedscale

import (
	"fmt"
	"io"
	"slices"
	"strconv"
	"strings"
)

const (
	X = iota
	Prop
	Bm
)

type State struct {
	W     io.Writer
	Scale int
	Name  string
	Mode  int
	K     string
	V     string
	LUT   map[rune]string
}

func NewState(w io.Writer, scale int, name string) *State {
	return &State{
		W:     w,
		Scale: scale,
		Name:  name,
		Mode:  X,
		LUT:   make(map[rune]string),
	}
}

func (state *State) Next() error {
	switch state.Mode {
	case Bm:
		return state.ModeBM()
	case Prop:
		return state.ModeProp()
	default:
		return state.ModeX()
	}
}

func (state *State) ModeX() error {
	fmt.Fprint(state.W, state.K)
	switch state.K {

	case "STARTPROPERTIES":
		state.Mode = Prop
		fmt.Fprint(state.W, " ", state.V)

	case "BITMAP":
		state.Mode = Bm

	case "FONT":
		if err := state.XLFD(); err != nil {
			return err
		}

	case "SIZE", "SWIDTH", "DWIDTH":
		if err := state.Vtoi(); err != nil {
			return fmt.Errorf("bad field %s", state.K)
		}

	case "FONTBOUNDINGBOX", "BBX":
		if err := state.Vstoi(); err != nil {
			return fmt.Errorf("bad field %s", state.K)
		}

	default:
		fmt.Fprint(state.W, " ", state.V)

	}
	fmt.Fprintln(state.W)
	return nil
}

func (state *State) ModeProp() error {
	fmt.Fprint(state.W, state.K)
	switch state.K {

	case "ENDPROPERTIES":
		state.Mode = X

	case "FAMILY_NAME":
		fmt.Fprint(state.W, ` "`, strings.ReplaceAll(state.Name, `"`, `""`), `"`)

	case
		"PIXEL_SIZE",
		"POINT_SIZE",
		"AVERAGE_WIDTH",
		"FONT_ASCENT",
		"FONT_DESCENT",
		"CAP_HEIGHT",
		"X_HEIGHT",
		"BITED_DWIDTH",
		"BITED_EDITOR_GRID_SIZE":
		if err := state.Vtoi(); err != nil {
			return fmt.Errorf("bad prop %s", state.K)
		}

	default:
		fmt.Fprint(state.W, " ", state.V)

	}
	fmt.Fprintln(state.W)
	return nil
}

func (state *State) ModeBM() error {
	switch state.K {
	case "ENDCHAR":
		fmt.Fprintln(state.W, state.K)
		state.Mode = X

	default:
		if err := state.ScaleHex(); err != nil {
			return err
		}

	}
	return nil
}

func (state *State) XLFD() error {
	fmt.Fprint(state.W, " ")

	xlfd := strings.Split(state.V, "-")
	for i, v := range xlfd {
		if i == 0 {
			continue
		}
		fmt.Fprint(state.W, "-")

		if i == 2 {
			fmt.Fprint(state.W, state.Name)
			continue
		}

		if i == 7 || i == 8 || i == 12 {
			if n, err := strconv.Atoi(v); err == nil {
				fmt.Fprint(state.W, n*state.Scale)
				continue
			}
			return fmt.Errorf("bad field FONT index %d", i)
		}

		fmt.Fprint(state.W, v)
	}
	return nil
}

func (state *State) Vtoi() error {
	fmt.Fprint(state.W, " ")

	a, b, f := strings.Cut(state.V, " ")
	n, err := strconv.Atoi(a)
	if err != nil {
		return err
	}

	fmt.Fprint(state.W, n*state.Scale)
	if f {
		fmt.Fprint(state.W, " ", b)
	}

	return nil
}

func (state *State) Vstoi() error {
	vs := strings.Fields(state.V)
	if len(vs) == 0 {
		return fmt.Errorf("nothing to parse")
	}

	for _, v := range vs {
		n, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		fmt.Fprint(state.W, " ", n*state.Scale)
	}

	return nil
}

func (state *State) ScaleHex() error {
	var line strings.Builder
	for _, c := range state.K {
		if (c < '0' || c > '9') && (c < 'A' || c > 'F') && (c < 'a' || c > 'f') {
			return fmt.Errorf("not hex %s", state.K)
		}

		if _, ok := state.LUT[c]; !ok {
			n, err := strconv.ParseUint(string(c), 16, 4)
			if err != nil {
				return err
			}
			b := fmt.Sprintf("%04b", n)
			var b1 strings.Builder

			for _, c := range b {
				for i := 0; i < state.Scale; i++ {
					b1.WriteRune(c)
				}
			}

			var b2 strings.Builder
			for chunk := range slices.Chunk([]rune(b1.String()), 4) {
				n1, err := strconv.ParseUint(string(chunk), 2, 4)
				if err != nil {
					return err
				}
				b2.WriteString(fmt.Sprintf("%X", n1))
			}
			state.LUT[c] = b2.String()
		}

		line.WriteString(state.LUT[c])
	}

	lstr := line.String()
	for i := 0; i < state.Scale; i++ {
		fmt.Fprintln(state.W, lstr)
	}

	return nil
}
