package fitzimg

import "io"

func serial(src []byte, dst io.Writer, params *Params) error {
	tmp, err := initTmp()
	if err != nil {
		return err
	}
	defer removeTmp(tmp)
	out, closer := initArchive(params.Archive, dst)
	bkg := background(params.Width, params.Height, params.Resize)
	from := params.FirstPage
	to := params.LastPage
	for page := from; page <= to; page++ {
		if err := process(src, bkg, tmp, out, page, params); err != nil {
			return err
		}
	}
	if closer != nil {
		return closer.Close()
	}
	return nil
}
