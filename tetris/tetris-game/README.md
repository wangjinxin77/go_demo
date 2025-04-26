# Tetris Game

## Overview
Tetris is a classic puzzle game where players manipulate falling tetrominoes to create complete lines on the game board. The game ends when the tetrominoes stack up to the top of the board.

## Installation
To install the Tetris game on Ubuntu 22, follow these steps:

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/tetris-game.git
   cd tetris-game
   ```

2. Run the installation script:
   ```bash
   ./scripts/install.sh
   ```

3. After installation, you can run the game using:
   ```bash
   go run cmd/tetris/main.go
   ```

## Gameplay Rules
- Use the arrow keys to move the tetrominoes left, right, and down.
- Press the spacebar to rotate the tetromino.
- The objective is to create complete horizontal lines on the board. When a line is completed, it will disappear, and you will earn points.
- The game ends when the tetrominoes reach the top of the board.

## Running the Game
To run the game, navigate to the project directory and execute the following command:
```bash
go run cmd/tetris/main.go
```

## Uninstallation
To uninstall the game, run the uninstall script:
```bash
./scripts/uninstall.sh
```

## Contributing
Feel free to contribute to the project by submitting issues or pull requests.

## License
This project is licensed under the MIT License. See the LICENSE file for details.