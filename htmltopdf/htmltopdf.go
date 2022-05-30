package htmltopdf

import (
	"context"
	"fmt"
	"io"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type HTMLToPDF interface {
	GenerateFromReader(ctx context.Context, page io.Reader) (pdf io.Reader, err error)
}

type htmlToPDF struct{}

func NewHTMLToPDF() HTMLToPDF {
	return &htmlToPDF{}
}

func (htp *htmlToPDF) GenerateFromReader(ctx context.Context, page io.Reader) (pdf io.Reader, err error) {
	generator, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return
	}
	generator.AddPage((wkhtmltopdf.NewPageReader(page)))
	generator.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	generator.Dpi.Set(300)
	fmt.Println(generator.Args())

	err = generator.CreateContext(ctx)
	if err != nil {
		return
	}

	pdf = generator.Buffer()

	return
}
