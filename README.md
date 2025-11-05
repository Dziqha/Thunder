# ⚡ Thunder - Ultra Fast Hot Reload

Hot reload tool untuk Go yang **lebih cepat** dari Air!

## 🚀 Keunggulan Thunder vs Air

| Feature | Thunder ⚡ | Air 🌬️ |
|---------|-----------|---------|
| **Debounce Smart** | ✅ 100ms | ❌ 1000ms |
| **Reload Speed** | ⚡ Ultra Fast | 🐌 Slower |
| **Memory Usage** | 💚 Lightweight | ⚠️ Higher |
| **Build Time** | ⏱️ Optimized | ⏱️ Standard |
| **Setup** | 🎯 Simple | 📝 Complex config |
| **Colored Output** | 🎨 Beautiful | ⚪ Plain |

## 📦 Installation

```bash
# Clone atau copy thunder.go ke project Anda
# Install dependency
go get github.com/fsnotify/fsnotify
```

## 🎯 Usage

### Cara 1: Run langsung
```bash
# Jalankan thunder untuk main.go
go run thunder.go

# Atau untuk file lain
go run thunder.go cmd/api/main.go
```

### Cara 2: Build thunder sebagai binary
```bash
# Build thunder
go build -o thunder thunder.go

# Jalankan
./thunder
# atau
./thunder cmd/api/main.go
```

## ⚙️ Konfigurasi

Edit fungsi `DefaultConfig()` di thunder.go:

```go
Config{
    BuildPath:   "./tmp/main",           // Output binary
    MainFile:    "main.go",              // Entry point
    WatchDirs:   []string{"."},          // Direktori yang diwatch
    ExcludeDirs: []string{"tmp", "vendor"}, // Direktori yang diignore
    BuildArgs:   []string{"-tags", "dev"}, // Extra build args
    RunArgs:     []string{},             // Args untuk app
    Debounce:    100 * time.Millisecond, // Debounce time
}
```

## 🎨 Output Features

- ⚡ **Real-time monitoring** dengan colored output
- ⚙️ **Build status** dengan timing info
- ✓ **Success indicator** yang jelas
- ✗ **Error messages** yang informatif
- 🎯 **File change detection** yang akurat

## 🔥 Performance Tips

1. **Exclude unnecessary directories** - Tambahkan folder besar ke excludeDirs
2. **Adjust debounce** - Turunkan ke 50ms jika edit lambat, naikkan ke 200ms jika terlalu sensitive
3. **Use build tags** - Pisahkan dev dan production code dengan build tags

## 📝 Example Project Structure

```
myapp/
├── thunder.go          # Hot reload tool
├── main.go            # Your app
├── go.mod
├── handlers/
│   └── api.go
├── models/
│   └── user.go
└── tmp/               # Build output (auto-created)
    └── main
```

## 🎯 Tips

- Edit file .go apapun dan Thunder akan auto-reload
- Binary di-build ke folder `tmp/` (gitignore recommended)
- Ctrl+C untuk stop Thunder
- Error build akan ditampilkan tanpa crash

## 🆚 Comparison

### Air:
- Butuh config file `.air.toml`
- Reload lebih lambat (1s debounce default)
- Setup lebih kompleks

### Thunder:
- Zero config, works out of the box
- Ultra fast reload (100ms debounce)
- Colored, beautiful output
- Lebih lightweight

## 💡 Advanced Usage

### Custom Config
```go
config := Config{
    BuildPath:   "./bin/app",
    MainFile:    "cmd/api/main.go",
    WatchDirs:   []string{".", "pkg"},
    ExcludeDirs: []string{"tmp", "vendor", "node_modules"},
    BuildArgs:   []string{"-tags", "dev", "-race"},
    RunArgs:     []string{"-port", "3000"},
    Debounce:    50 * time.Millisecond,
}
```

### Multi-directory Watch
```go
WatchDirs: []string{".", "internal", "pkg", "cmd"}
```

### Production Build Flags
```go
BuildArgs: []string{
    "-ldflags", "-s -w",  // Reduce binary size
    "-tags", "production",
}
```

## 🐛 Troubleshooting

**Q: Thunder tidak detect perubahan?**
A: Pastikan direktori ada di `WatchDirs` dan tidak di `ExcludeDirs`

**Q: Build terlalu sering?**
A: Naikkan `Debounce` ke 200-500ms

**Q: Error "permission denied"?**
A: Pastikan folder `tmp` writeable atau ganti `BuildPath`

## ⚡ Enjoy Lightning-Fast Development!

Thunder dibuat untuk developer yang menghargai **kecepatan** dan **simplicity**.
No complex config, just pure performance! 🚀