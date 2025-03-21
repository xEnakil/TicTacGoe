@echo off

:: Create dist directory
mkdir dist 2>nul

:: Build for Windows
echo Building for Windows...
go build -ldflags="-H windowsgui" -o dist/TicTacToe.exe

:: Build for macOS
echo Building for macOS...
set GOOS=darwin
set GOARCH=amd64
go build -o dist/TicTacToe_mac
set GOOS=windows
set GOARCH=amd64

echo Build complete!
echo Windows executable: dist\TicTacToe.exe
echo macOS executable: dist\TicTacToe_mac
