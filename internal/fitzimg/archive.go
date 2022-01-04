package fitzimg

import (
	"archive/tar"
	"archive/zip"
	"io"
	"os"
)

func initArchive(archive Archive, dst io.Writer) (interface{}, io.Closer) {
	switch archive {
	case Tar:
		t := tar.NewWriter(dst)
		return t, t
	case Zip:
		z := zip.NewWriter(dst)
		return z, z
	default:
		return dst, nil
	}
}

func write(archive Archive, w interface{}, f *os.File, name string) error {
	switch archive {
	case Tar:
		return writeTarEntry(w.(*tar.Writer), f, name)
	case Zip:
		return writeZipEntry(w.(*zip.Writer), f, name)
	default:
		return writeRawEntry(w.(io.Writer), f)
	}
}

func writeRawEntry(w io.Writer, f *os.File) error {
	_, err := io.Copy(w, f)
	return err
}

func writeTarEntry(w *tar.Writer, f *os.File, name string) error {
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	fh, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return err
	}
	fh.Name = name
	if err := w.WriteHeader(fh); err != nil {
		return err
	}
	return writeRawEntry(w, f)
}

func writeZipEntry(w *zip.Writer, f *os.File, name string) error {
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	fh, err := zip.FileInfoHeader(fi)
	if err != nil {
		return err
	}
	fh.Method = zip.Store
	fh.Name = name
	hw, err := w.CreateHeader(fh)
	if err != nil {
		return err
	}
	return writeRawEntry(hw, f)
}
