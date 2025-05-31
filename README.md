# Axiom Shift

Axiom Shift is a strategic and intellectual battle game that provides a unique gaming experience through trial and error, allowing players to grow and learn from their actions. The game features a character that players can develop, alongside an enemy that evolves in parallel, culminating in a final battle.

## Project Structure

```
axiom_shift
├── cmd
│   └── main.go          # Entry point of the application
├── internal
│   ├── game
│   │   ├── game.go      # Main game loop and state management
│   │   ├── matrix.go    # Matrix operations and definitions
│   │   ├── player.go    # Player character logic
│   │   ├── enemy.go     # Enemy logic
│   │   └── battle.go    # Battle processing and outcome determination
│   ├── ui
│   │   ├── ui.go        # UI rendering
│   │   └── input.go     # User input management
│   └── logic
│       ├── rules.go     # Rule matrix definitions
│       └── seed.go      # Seed management for reproducibility
├── go.mod                # Module definition
├── go.sum                # Module checksums
└── README.md             # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.16 or later
- Ebiten library

### Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/axiom_shift.git
   cd axiom_shift
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

### Running the Game

To run the game, execute the following command:
```
go run cmd/main.go
```

### Gameplay

- Players input a floating-point value between 0 and 1 before each battle, influencing their character's matrix state.
- The game consists of multiple battles, with the final battle determining the overall outcome.
- Players must observe and adapt their strategies based on previous inputs and results.

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

### License

This project is licensed under the MIT License. See the LICENSE file for more details.