package fitzimg

import (
	"io"
	"os"
	"sync/atomic"
)

func interleave(src []byte, dst io.Writer, params *Params) error {
	if params.FirstPage == params.LastPage {
		// Fallback to serial for single pages
		return serial(src, dst, params)
	}
	duo, err := initTmpDuo()
	if err != nil {
		return err
	}
	defer removeTmpDuo(duo)
	out, closer := initArchive(params.Archive, dst)
	receive := make(chan error, 1)
	cancel := int32(0)
	go work(receive, &cancel, src, duo, params)
	for page := params.FirstPage; page <= params.LastPage; page++ {
		if err := <-receive; err != nil {
			return err
		}
		n := name(page, params.Format)
		if err := transfer(duo[page&1], out, n, params); err != nil {
			atomic.StoreInt32(&cancel, 1)
			return err
		}
	}
	if closer != nil {
		closer.Close()
	}
	return nil
}

func work(send chan error, cancel *int32, src []byte, duo []string, params *Params) {
	defer close(send)
	bkg := background(params.Width, params.Height, params.Resize)
	from := params.FirstPage
	to := params.LastPage
	for page := from; page <= to && atomic.LoadInt32(cancel) == 0; page++ {
		err := dump(src, bkg, duo[page&1], page, params)
		if send <- err; err != nil {
			return
		}
	}
}

func initTmpDuo() ([]string, error) {
	tmp1, err := initTmp()
	if err != nil {
		return nil, err
	}
	tmp2, err := initTmp()
	if err != nil {
		os.Remove(tmp1)
		return nil, err
	}
	return []string{tmp1, tmp2}, nil
}

func removeTmpDuo(duo []string) {
	removeTmp(duo[0])
	removeTmp(duo[1])
}
