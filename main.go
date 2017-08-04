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
	master  []byte
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

func RegistFiles(filenames ...string) (map[string]*File, error) {
	files := map[string]*File{}
	for _, v := range filenames {
		var err error
		files[v], err = RegistFile(v)
		if err != nil {
			return nil, err
		}
	}
	return files, nil
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
	f.master = make([]byte, len(f.Body))
	copy(f.master, f.Body)
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

// ファイル内容と更新日時の両方更新
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

// 登録されたファイルが最新版かどうか判定する
func (f *File) IsLatest() (bool, error) {
	fInfo, err := os.Stat(f.Name)
	if err != nil {
		return false, err
	}
	return f.ModTime == fInfo.ModTime(), nil
}

// ファイルの最新の中身をとる
func (f *File) GetBytes() ([]byte, error) {
	islatest, err := f.IsLatest()
	if err != nil {
		return nil, err
	}
	if islatest {
		return f.master, nil
	}
	err = f.Update()
	if err != nil {
		return nil, err
	}
	return f.master, nil
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
	f.master = make([]byte, len(f.Body))
	copy(f.master, f.Body)
	return nil
}

// File.Bodyを保存する
func (f *File) Save() error {
	file, err := os.Create(f.Name)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Write(f.master)
	return nil
}

func (f *File) LatestBody() bool {
	return string(f.Body) == string(f.master)
}

// bodyを確定して保存可能に
func (f *File) CommitBody() {
	f.master = make([]byte, len(f.Body))
	copy(f.master, f.Body)
}

// bodyを前のUpdateまで巻き戻す
func (f *File) RollBackBody() {
	f.Body = make([]byte, len(f.master))
	copy(f.Body, f.master)
}
