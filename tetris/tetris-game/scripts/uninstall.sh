#!/bin/bash

# Uninstall script for Tetris Game

echo "Uninstalling Tetris Game..."

# Remove the desktop entry
if [ -f /usr/share/applications/tetris.desktop ]; then
    sudo rm /usr/share/applications/tetris.desktop
    echo "Removed desktop entry."
else
    echo "Desktop entry not found."
fi

# Remove game files (assuming they are installed in /usr/local/bin)
if [ -f /usr/local/bin/tetris ]; then
    sudo rm /usr/local/bin/tetris
    echo "Removed game executable."
else
    echo "Game executable not found."
fi

echo "Tetris Game has been uninstalled."