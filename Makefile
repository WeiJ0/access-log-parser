# Makefile for Apache Log Analyzer
# 提供統一的建置、測試和執行命令

.PHONY: all build test clean run dev install lint fmt help generate-testdata

# 預設目標
all: fmt lint test build

# 建置應用程式
build:
	@echo "建置應用程式..."
	wails build

# 建置（帶壓縮）
build-compressed:
	@echo "建置應用程式（啟用 UPX 壓縮）..."
	wails build -clean -upx

# 建置安裝程式
build-installer:
	@echo "建置 Windows 安裝程式..."
	wails build -nsis

# 開發模式
dev:
	@echo "啟動開發模式..."
	wails dev

# 執行測試
test:
	@echo "執行單元測試..."
	go test -v ./...

# 執行測試（含覆蓋率）
test-coverage:
	@echo "執行測試並產生覆蓋率報告..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆蓋率報告已產生: coverage.html"

# 執行基準測試
bench:
	@echo "執行基準測試..."
	go test -bench=. -benchmem ./...

# 執行特定套件的基準測試
bench-parser:
	@echo "執行解析器基準測試..."
	go test -bench=. -benchmem ./internal/parser

# 安裝依賴
install:
	@echo "安裝 Go 依賴..."
	go mod download
	@echo "安裝前端依賴..."
	cd frontend && npm install

# 程式碼檢查
lint:
	@echo "執行程式碼檢查..."
	go vet ./...
	gofmt -l .

# 格式化程式碼
fmt:
	@echo "格式化程式碼..."
	gofmt -w .

# 清理建置產物
clean:
	@echo "清理建置產物..."
	rm -rf build/bin
	rm -rf frontend/dist
	rm -f coverage.out coverage.html
	@echo "清理完成"

# 產生測試資料
generate-testdata:
	@echo "產生測試資料..."
	go run scripts/generate_test_log.go -lines 100 -output testdata/valid.log -error-rate 0.05 -invalid-rate 0
	go run scripts/generate_test_log.go -lines 100 -output testdata/invalid.log -error-rate 0.05 -invalid-rate 0.2
	go run scripts/generate_test_log.go -lines 1000000 -output testdata/100mb.log -error-rate 0.05 -invalid-rate 0.01
	@echo "測試資料產生完成"

# 執行所有檢查（格式、檢查、測試）
check: fmt lint test
	@echo "所有檢查通過！"

# 顯示幫助
help:
	@echo "Apache Log Analyzer - Make 命令"
	@echo ""
	@echo "可用命令:"
	@echo "  make build              - 建置應用程式"
	@echo "  make build-compressed   - 建置應用程式（啟用壓縮）"
	@echo "  make build-installer    - 建置 Windows 安裝程式"
	@echo "  make dev                - 啟動開發模式"
	@echo "  make test               - 執行單元測試"
	@echo "  make test-coverage      - 執行測試並產生覆蓋率報告"
	@echo "  make bench              - 執行基準測試"
	@echo "  make bench-parser       - 執行解析器基準測試"
	@echo "  make install            - 安裝所有依賴"
	@echo "  make lint               - 執行程式碼檢查"
	@echo "  make fmt                - 格式化程式碼"
	@echo "  make clean              - 清理建置產物"
	@echo "  make generate-testdata  - 產生測試資料"
	@echo "  make check              - 執行所有檢查"
	@echo "  make help               - 顯示此幫助訊息"
