package fitzimg

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"bytes"
	"io"
	"os"
)

type ArchiveWriter interface {
	StartEntry(name string) (io.Writer, error)
	FinishEntry() error
	Write(f *os.File, name string) (int64, error)
	Close() error
}

type tarWriter struct {
	buf  *bytes.Buffer
	name string
	tar  *tar.Writer
}

func (w *tarWriter) StartEntry(name string) (io.Writer, error) {
	w.buf.Reset()
	w.name = name
	return bufio.NewWriter(w.buf), nil
}

func (w tarWriter) FinishEntry() error {
	defer w.buf.Reset()
	hdr := &tar.Header{
		Name: w.name,
		Mode: 0600,
		Size: int64(w.buf.Len()),
	}
	if err := w.tar.WriteHeader(hdr); err != nil {
		return err
	} else if _, err := w.buf.WriteTo(w.tar); err != nil {
		return err
	}
	return nil
}

func (w tarWriter) Write(f *os.File, name string) (int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	fh, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		return 0, err
	}
	fh.Name = name
	if err := w.tar.WriteHeader(fh); err != nil {
		return 0, err
	}
	return io.Copy(w.tar, f)
}

func (w tarWriter) Close() error {
	return w.tar.Close()
}

type zipWriter struct {
	zip *zip.Writer
}

func (w zipWriter) StartEntry(name string) (io.Writer, error) {
	return w.zip.Create(name)
}

func (w zipWriter) FinishEntry() error {
	return nil
}

func (w zipWriter) Write(f *os.File, name string) (int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, err
	}
	fh, err := zip.FileInfoHeader(fi)
	if err != nil {
		return 0, err
	}
	fh.Method = zip.Store
	fh.Name = name
	hw, err := w.zip.CreateHeader(fh)
	if err != nil {
		return 0, err
	}
	return io.Copy(hw, f)
}

func (w zipWriter) Close() error {
	return w.zip.Close()
}

type rawWriter struct {
	raw io.Writer
}

func (w rawWriter) StartEntry(_ string) (io.Writer, error) {
	return w.raw, nil
}

func (w rawWriter) FinishEntry() error {
	return nil
}

func (w rawWriter) Write(f *os.File, _ string) (int64, error) {
	return io.Copy(w.raw, f)
}

func (w rawWriter) Close() error {
	return nil
}

func initArchive(archive Archive, dst io.Writer) ArchiveWriter {
	switch archive {
	case Tar:
		return &tarWriter{new(bytes.Buffer), "", tar.NewWriter(dst)}
	case Zip:
		return &zipWriter{zip.NewWriter(dst)}
	default:
		return &rawWriter{dst}
	}
}
