package pdf

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ahmedalkabir/garabic"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/image"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/extension"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"github.com/johnfercher/maroto/v2/pkg/repository"
)

//go:embed assets/*
var assets embed.FS

type options struct {
	headerTitle           *string
	phoneNumber           *string
	webSite               *string
	customFont            *string
	customRegularFontPath *string
	customBoldFontPath    *string
	logoPath              *string
	logoByte              []byte
	customRegularFontByte []byte
	customBoldFontByte    []byte
}

type Option func(options *options) error

func WithInitailInfo(title, phone, website string) Option {
	return func(options *options) error {

		options.headerTitle = &title
		options.phoneNumber = &phone
		options.webSite = &website

		return nil
	}
}

func WithCustomRegularFontPath(font, path string) Option {
	return func(options *options) error {
		if font == "" || path == "" {
			return nil
		}

		// check if the path exist
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			return err
		}

		options.customFont = &font
		options.customRegularFontPath = &path
		return nil
	}
}

func WithCustomBoldFontPath(font, path string) Option {
	return func(options *options) error {
		if font == "" || path == "" {
			return nil
		}

		// check if the path exist
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			return err
		}

		options.customFont = &font
		options.customBoldFontPath = &path
		return nil
	}
}

func WithCustomRegularFontByte(font, file string) Option {
	return func(options *options) error {
		if font == "" || file == "" {
			return nil
		}

		// check if the file exist
		_byte, err := assets.ReadFile(strings.Join([]string{"assets", "/", "fonts", "/", file}, ""))
		if err != nil {
			return err
		}

		options.customFont = &font
		options.customRegularFontByte = _byte
		return nil
	}
}

func WithCustomBoldFontByte(font, file string) Option {
	return func(options *options) error {
		if font == "" {
			return nil
		}

		// check if the file exist
		_byte, err := assets.ReadFile(strings.Join([]string{"assets", "/", "fonts", "/", file}, ""))
		if err != nil {
			return err
		}

		options.customFont = &font
		options.customBoldFontByte = _byte
		return nil
	}
}

func WithLogoPath(path string) Option {
	return func(options *options) error {
		if path == "" {
			return nil
		}

		// check if the path exist
		if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
			return err
		}

		options.logoPath = &path
		return nil
	}
}

func WithLogoByte(name string) Option {
	return func(options *options) error {
		if name == "" {
			return nil
		}

		// check if the file exist
		_byte, err := assets.ReadFile(strings.Join([]string{"assets", "/", name}, ""))
		if err != nil {
			return err
		}

		options.logoByte = _byte
		return nil
	}
}

type PDF struct {
	MarotoPDF core.Maroto
	options   options
}

func GetDarkGrayColor() *props.Color {
	return &props.Color{
		Red:   55,
		Green: 55,
		Blue:  55,
	}
}

func GetGrayColor() *props.Color {
	return &props.Color{
		Red:   200,
		Green: 200,
		Blue:  200,
	}
}

func getBlueColor() *props.Color {
	return &props.Color{
		Red:   10,
		Green: 10,
		Blue:  150,
	}
}

func getRedColor() *props.Color {
	return &props.Color{
		Red:   150,
		Green: 10,
		Blue:  10,
	}
}

func (p *PDF) buildPDF() (core.Maroto, error) {
	cfg := config.NewBuilder().
		WithPageNumber(props.PageNumber{
			Pattern: "Page {current} of {total}",
		})

	if p.options.customFont != nil {
		repo := repository.New()

		if p.options.customRegularFontPath != nil {
			repo.AddUTF8Font(*p.options.customFont, fontstyle.Normal, *p.options.customRegularFontPath).
				AddUTF8Font(*p.options.customFont, fontstyle.Italic, *p.options.customRegularFontPath).
				AddUTF8Font(*p.options.customFont, fontstyle.BoldItalic, *p.options.customRegularFontPath)

		} else if len(p.options.customRegularFontByte) > 0 {
			repo.AddUTF8FontFromBytes(*p.options.customFont, fontstyle.Normal, p.options.customRegularFontByte).
				AddUTF8FontFromBytes(*p.options.customFont, fontstyle.Italic, p.options.customRegularFontByte).
				AddUTF8FontFromBytes(*p.options.customFont, fontstyle.BoldItalic, p.options.customRegularFontByte)
		}

		if p.options.customBoldFontPath != nil {
			repo.AddUTF8Font(*p.options.customFont, fontstyle.Bold, *p.options.customBoldFontPath)
		} else if len(p.options.customBoldFontByte) > 0 {
			repo.AddUTF8FontFromBytes(*p.options.customFont, fontstyle.Bold, p.options.customBoldFontByte)
		}

		customFonts, err := repo.Load()
		if err != nil {
			return nil, err
		}

		cfg.WithCustomFonts(customFonts).
			WithDefaultFont(&props.Font{Family: *p.options.customFont})
	}

	m := maroto.New(cfg.Build())

	headerRow := row.New(20)

	if p.options.headerTitle != nil {
		headerRow.Add(col.New(3).Add(
			text.New(garabic.Shape(*p.options.headerTitle), props.Text{
				Size:  12,
				Align: align.Center,
			}),
			text.New(garabic.Shape(fmt.Sprintf("رقم الهاتف %s : ", *p.options.phoneNumber)), props.Text{
				Top:   7,
				Style: fontstyle.BoldItalic,
				Size:  12,
				Align: align.Center,
				// Color: getBlueColor(),
			}),
			text.New(*p.options.webSite, props.Text{
				Top:   12,
				Style: fontstyle.BoldItalic,
				Size:  12,
				Align: align.Center,
				// Color: getBlueColor(),
			}),
		))

	}

	if p.options.logoPath != nil {
		headerRow.Add(
			col.New(6),
			image.NewFromFileCol(3, *p.options.logoPath, props.Rect{
				Center:  true,
				Percent: 100,
			}))
	} else if len(p.options.logoByte) > 0 {
		headerRow.Add(
			col.New(6),
			image.NewFromBytesCol(3, p.options.logoByte, extension.Png, props.Rect{
				Center:  true,
				Percent: 100,
			}))
	}

	err := m.RegisterHeader(headerRow)
	if err != nil {
		return nil, err
	}

	footerRow := row.New(20)

	if p.options.headerTitle != nil {
		footerRow.Add(
			col.New(12).Add(
				text.New(fmt.Sprintf("Tel: %s", *p.options.phoneNumber), props.Text{
					Top:   13,
					Style: fontstyle.BoldItalic,
					Size:  8,
					Align: align.Left,
					Color: getBlueColor(),
				}),
				text.New(*p.options.webSite, props.Text{
					Top:   16,
					Style: fontstyle.BoldItalic,
					Size:  8,
					Align: align.Left,
					Color: getBlueColor(),
				}),
			),
		)
	}

	err = m.RegisterFooter(footerRow)

	if err != nil {
		return nil, err
	}

	return m, nil
}

func NewPDF(opts ...Option) (*PDF, error) {
	var pdf PDF

	for _, opt := range opts {
		err := opt(&pdf.options)
		if err != nil {
			return nil, err
		}
	}

	return &pdf, nil
}

func (p *PDF) GenerateCustomFunc(fn func(pdf *PDF)) ([]byte, error) {

	m, err := p.buildPDF()
	if err != nil {
		return nil, err
	}
	p.MarotoPDF = m

	fn(p)

	document, err := p.MarotoPDF.Generate()
	if err != nil {
		return nil, err
	}

	return document.GetBytes(), err
}
