// ファイルの変更を確認するパッケージ
package checkmodfile

import (
	"bytes"
	"io"
	"os"
	"time"
)

type File struct {
	Name    string
	ModTime time.Time
	head    os.FileInfo
	Body    []byte
}

// 管理対象に登録
func RegistFile(filename string) (*File, error) {
	f := new(File)
	f.Name = filename
	err := f.Update()
	if err != nil {
		return nil, err
	}
	return f, nil
}

// 現在のファイル内容を取り出す
func (f *File) UpdateBody() error {
	fr, err := os.Open(f.Name)
	if err != nil {
		return err
	}
	defer fr.Close()
	b := new(bytes.Buffer)
	io.Copy(b, fr)
	f.Body = b.Bytes()
	return nil
}

// ModTimeを新しい物に更新
func (f *File) UpdateHead() error {
	fInfo, err := os.Stat(f.Name)
	if err != nil {
		return err
	}
	f.head = fInfo
	return nil
}

// ファイル内容と更新日時の両方更新
func (f *File) Update() error {
	err := f.UpdateHead()
	if err != nil {
		return err
	}
	err = f.UpdateBody()
	if err != nil {
		return err
	}

	return nil
}

// 登録されたファイルが最新版かどうか判定する
func (f *File) IsLatest() (bool, error) {
	fInfo, err := os.Stat(f.Name)
	if err != nil {
		return false, err
	}
	return f.head.ModTime() == fInfo.ModTime(), nil
}

// ファイルの最新の中身をとる
func (f *File) GetBytes() ([]byte, error) {
	islatest, err := f.IsLatest()
	if err != nil {
		return nil, err
	}
	if islatest {
		return f.Body, nil
	}
	err = f.Update()
	if err != nil {
		return nil, err
	}
	return f.Body, nil
}

func (f *File) WriteTo(w io.Writer) error {
	latest, err := f.IsLatest()
	if err != nil {
		return err
	}
	if latest {
		w.Write(f.Body)
		return nil
	}
	b := new(bytes.Buffer)
	mw := io.MultiWriter(w, b)
	fr, err := os.Open(f.Name)
	if err != nil {
		return err
	}
	defer fr.Close()
	io.Copy(mw, fr)
	f.Body = b.Bytes()
	return nil
}

func (f *File) Save(b []byte) error {
	file, err := os.Create(f.Name)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(b)
	return nil
}
