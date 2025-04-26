#!/bin/bash

# Install dependencies
echo "Installing dependencies..."
go mod tidy

# Create desktop entry for the game
echo "Creating desktop entry..."
cat <<EOL > ~/.local/share/applications/tetris.desktop
[Desktop Entry]
Name=Tetris Game
Exec=go run cmd/tetris/main.go
Type=Application
Terminal=true
EOL

echo "Installation complete! You can now run the game from your application menu."