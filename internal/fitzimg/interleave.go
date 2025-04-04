package fitzimg

import (
	"fmt"
	"io"
	"os"
	"sync/atomic"

	"github.com/gen2brain/go-fitz"
)

func interleave(doc *fitz.Document, dst io.Writer, params *Params) error {
	from := params.FirstPage
	to := params.LastPage
	if from == to {
		// Fallback to serial for single pages
		return serial(doc, dst, params)
	}
	duo, err := initTmpDuo()
	if err != nil {
		return fmt.Errorf("fitzimg.interleave initTmpDuo: %w", err)
	}
	defer removeTmpDuo(duo)
	out := initArchive(params.Archive, dst)
	receive := make(chan error, 1)
	cancel := new(atomic.Bool)
	go work(receive, cancel, doc, duo, params)
	for page := from; page <= to; page++ {
		if err := <-receive; err != nil {
			return fmt.Errorf("fitzimg.interleave work page=%d: %w", page, err)
		}
		n := name(page, params.Format)
		if err := transfer(duo[page&1], out, n); err != nil {
			cancel.Store(true)
			<-receive
			return fmt.Errorf("fitzimg.interleave transfer page=%d: %w", page, err)
		}
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("fitzimg.interleave out.Close: %w", err)
	}
	return nil
}

func work(send chan error, cancel *atomic.Bool, doc *fitz.Document, duo [2]string, params *Params) {
	defer close(send)
	bkg := background(params.Width, params.Height, params.Resize)
	from := params.FirstPage
	to := params.LastPage
	for page := from; page <= to && !cancel.Load(); page++ {
		err := dump(doc, bkg, duo[page&1], page, params)
		if send <- err; err != nil {
			return
		}
	}
}

func initTmpDuo() ([2]string, error) {
	var duo [2]string
	tmp1, err := initTmp()
	if err != nil {
		return duo, err
	}
	tmp2, err := initTmp()
	if err != nil {
		os.Remove(tmp1)
		return duo, err
	}
	duo[0] = tmp1
	duo[1] = tmp2
	return duo, nil
}

func removeTmpDuo(duo [2]string) {
	removeTmp(duo[0])
	removeTmp(duo[1])
}
