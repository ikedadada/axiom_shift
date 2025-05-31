# Axiom Shift

Axiom Shift is a strategic and intellectual battle game that provides a unique gaming experience through trial and error, allowing players to grow and learn from their actions. The game features a character that players can develop, alongside an enemy that evolves in parallel, culminating in a final battle.

## Project Structure

```
axiom_shift/
├── main.go               # Entry point (DI and Ebiten startup only)
├── internal/
│   ├── domain/           # Entities, value objects, domain logic
│   ├── usecase/          # Application use cases (battle flow, etc.)
│   ├── logic/            # Pure logic (rules, random, etc.)
│   ├── game/             # Game state management, main loop
│   └── ui/               # UI logic, rendering, input
├── assets/               # Fonts, images, etc.
├── docs/                 # Documentation
├── go.mod                # Go module definition
├── go.sum                # Go module checksums
└── README.md             # Project documentation
```

## Getting Started

### Prerequisites

- Go 1.24.3 or later
- Ebiten library (latest stable)

### Installation

1. Clone the repository:

   ```zsh
   git clone https://github.com/yourusername/axiom_shift.git
   cd axiom_shift
   ```

2. Install dependencies:
   ```zsh
   go mod tidy
   ```

### Running the Game

To run the game, execute the following command:

```zsh
go run main.go
```

### Gameplay

- Players input a floating-point value between 0 and 1 before each battle, influencing their character's matrix state.
- The game consists of multiple battles (default: 10), with the final battle determining the overall outcome.
- Player and enemy matrices are visualized in the UI using colored rectangles (blue for player, red for enemy).
- The rule matrix is generated from a visible seed and remains fixed during a session.
- All logic except Ebiten-dependent UI is fully unit tested with 100% coverage and parameterized tests.
- Players must observe and adapt their strategies based on previous inputs and results.

### Development Policy

- See `docs/dev_policy.md` for architecture, testing, and coverage requirements.
- All code (except UI/ebiten-dependent) must be covered by parameterized unit tests.
- UI rendering uses Ebiten's Draw methods and color visualization for matrices.

### Contributing

Contributions are welcome! Please open an issue or submit a pull request for any enhancements or bug fixes.

### License

This project is licensed under the MIT License. See the LICENSE file for more details.
