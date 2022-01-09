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
	for page := params.FirstPage; page <= params.LastPage; page++ {
		if err := process(src, bkg, tmp, out, page, params); err != nil {
			return err
		}
	}
	if closer != nil {
		return closer.Close()
	}
	return nil
}
