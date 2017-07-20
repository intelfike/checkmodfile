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

// ファイル内容を取り出す
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
func (f *File) UpdateMod() error {
	fInfo, err := os.Stat(f.Name)
	if err != nil {
		return err
	}
	f.ModTime = fInfo.ModTime()
	return nil
}

// 両方アップデート
func (f *File) Update() error {
	err := f.UpdateMod()
	if err != nil {
		return err
	}
	err = f.UpdateBody()
	if err != nil {
		return err
	}

	return nil
}

//
func (f *File) IsLatest() (bool, error) {
	fInfo, err := os.Stat(f.Name)
	if err != nil {
		return false, err
	}
	return f.ModTime == fInfo.ModTime(), nil
}

// ファイルの最新の中身をとる
func (f *File) GetLatest() ([]byte, error) {
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
