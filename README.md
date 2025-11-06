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
 go install github.com/Dziqha/Thunder/cmd/thunder@latest
```

## 🎯 Usage

### Cara 1: Run Project
```bash
# Jalankan thunder untuk file.go
Thunder run
```

### Cara 2: Init Project
```bash
# Build thunder
Thunder init
```

## 🎨 Output Features

- ⚡ **Real-time monitoring** dengan colored output
- ⚙️ **Build status** dengan timing info
- ✓ **Success indicator** yang jelas
- ✗ **Error messages** yang informatif
- 🎯 **File change detection** yang akurat

## 🎯 Tips

- Edit file .go apapun dan Thunder akan auto-reload
- Binary di-build ke folder `tmp/` (gitignore recommended)
- Ctrl+C untuk stop Thunder
- Error build akan ditampilkan tanpa crash


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