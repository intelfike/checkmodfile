# checkmodfile is golang package

# Type, Function and Method
- 型<br>
type File struct{ ... }<br><br>
- 新しいファイルを管理対象に登録<br>
func RegistFile(filename string) (*File, error)<br><br>
- ファイル内容を取得。更新されてない場合はメモリ上に保存されたものを取得<br>
func (f *File) GetBytes() ([]byte, error)<br><br>
- 登録したファイルが最新か判定する<br>
func (f *File) IsLatest() (bool, error)<br><br>
- 登録されたファイルの内容と更新日付を更新する<br>
func (f *File) Update() error<br><br>
- 登録されたファイルの内容を更新する<br>
func (f *File) UpdateBytes() error<br><br>
- 登録されたファイルの更新日時を更新する<br>
func (f *File) UpdateMod() error<br><br>
- 渡されたwriterに書き出す。更新されていない場合はメモリ上に保存されたものを書き出す<br>
func (f *File) WriteTo(io.Writer) error <br><br>

# Example
HTTPサーバーをサンプルに用います。<br>
ベンチマークコマンドは以下のものを用います
```
ab -n 10000 -c 10 localhost
```

## キャッシュしない場合 (約14000 [#/sec])
下２つと比べると2割ほど低速です。
```
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// HTMLを返す
	file, err := os.Open("data/index.html")
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	io.Copy(w, file)
})
```
## 自力でキャッシュを実装する場合 (約 17000 [#/sec])
高速ですが、エラー処理を毎度書くとなかなか行数が長くなります。
```
// キャッシュしたいファイル
fileName := "data/index.html"
// ファイルを開く
ff, err := os.Open(fileName)
if err != nil {
	fmt.Println(err)
	return
}
fileBuffer := new(bytes.Buffer)
io.Copy(fileBuffer, ff)
// ファイル内容をｍ
fileText := fileBuffer.Bytes()
ff.Close()
fInfo, err := os.Stat(fileName)
if err != nil {
	fmt.Println(err)
	return
}
// 時間を保存
modTime := fInfo.ModTime()
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// HTMLを返す
	fInfo, err = os.Stat(fileName)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
	// ファイルが更新されていなければ、メモリ上のものを返す
	if modTime == fInfo.ModTime() {
		w.Write(fileText)
		return
	}
	// ファイルが更新されていれば、ファイルを読み取って返す
	f, err := os.Open("data/index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	mw := io.MultiWriter(f, fileBuffer)
	io.Copy(mw, f)
	fileText = fileBuffer.Bytes()
})
```

## このパッケージを使った場合 (約 17000 [#/sec])
単純さと高速性を兼ね備えています。
```
// ファイルを登録
f, err := checkmodfile.RegistFile("data/index.html")
if err != nil {
	fmt.Println(err)
	return
}
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// HTMLを返す
	err := f.WriteTo(w)
	if err != nil {
		fmt.Fprint(w, err)
		return
	}
})
```

# License
MIT

&copy;2017- intelfike<br><br><br>