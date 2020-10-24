package figlet

import "github.com/mbndr/figlet4go"

// A Transformer applies a Figlet transformation to a given string
type Transformer interface {
	Figletize(string) (string, error)
}

// NewTransformer returns a new Figlet transformer
func NewTransformer() Transformer {
	return &libFigletTransformer{}
}

type libFigletTransformer struct{}

func (t *libFigletTransformer) Figletize(in string) (string, error) {
	r := figlet4go.NewAsciiRender()

	figletedTxt, err := r.Render(in)
	if err != nil {
		return "", nil
	}

	return figletedTxt, nil
}
