package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
)

func loadImageAsBase64(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}

func (a *App) GetImageData(ensOrAddr string) string {
	addr := ensOrAddr
	// if addr, _ := a.conn.GetEnsAddress(ensOrAddr); len(addr) < 42 { // base.HexToAddress(addr) == base.ZeroAddr || !base.IsValidAddress(addr) {
	// 	logger.Error(fmt.Errorf("ENS not registered: %s", ensOrAddr))
	// 	return ""
	// } else {
	folder := "./output/generated/"
	fn := filepath.Join(folder, fmt.Sprintf("%s.png", addr))
	if file.FileExists(fn) {
		base64Image, err := loadImageAsBase64("path/to/your/image.png")
		if err != nil {
			return ""
		}
		return "data:image/png;base64," + base64Image
	} else {
		logger.Fatal(fn + " not found")
	}
	// }
	return ""
}
