# Copilot Instructions for Free-Hands-Onmyoji

## Project Overview

**Free-Hands-Onmyoji** is a Go-based game automation tool for the mobile game "Onmyoji" (阴阳师). It simulates mouse clicks and screen recognition to automate repetitive in-game tasks like exploration, dungeon runs, and event battles in the BlueStacks Android emulator on macOS.

### Tech Stack

- **Language:** Go 1.24+
- **Core Libraries:**
  - `gocv.io/x/gocv` - OpenCV bindings for template matching
  - `robotgo` - Cross-platform mouse/keyboard control
  - `toml` - Configuration file parsing

## Architecture Overview

The codebase follows a **State Machine + Task Registration** pattern for extensible automation.

### Core Components

1. **State Machine** (`internal/statemachine/`) - Orchestrates task execution flow

   - Tasks registered as `NamedTask` interface implementations
   - Tasks transition via `Next(key)` or `NextIndex()` methods
   - Shared state via `attributes` map for inter-task communication

2. **Task Registrator** (`internal/onmyoji/registrator.go`) - Registers mode-specific tasks

   - Each game mode (k28, breaker, guren, limitedEvents) has its own `Registrator` implementing `TaskRegistration` interface
   - Loads image templates (PNG files) for template matching
   - Registers tasks in execution order

3. **Window Platform Abstraction** (`internal/onmyoji/window/`) - Cross-platform window control

   - Abstract `Platform` interface with macOS/Windows implementations
   - Handles BlueStacks window positioning on primary/secondary displays
   - Mouse clicks translated to screen coordinates via `robotgo`

4. **Image-Based Detection** - Core recognition mechanism
   - Template images stored in `./k28/`, `./breaker/`, `./guren/`, `./limitedEvents/` directories
   - OpenCV template matching finds UI elements in screenshots
   - `ImgInfo` struct stores image path, dimensions, and pre-loaded Mat

### Task Execution Flow

```
main.go (parse args & initialize)
  ↓
onmyoji.Registrator (load mode-specific tasks & templates)
  ↓
StateMachine.Run() → Task.Execute()
  ↓
Task updates state attributes & transitions via TaskController
  ↓
Repeat until exit signal (Cmd+Shift+O) or timeout
```

## Key Patterns & Conventions

### 1. Adding a New Game Mode

Each mode needs:

- **Registrator** in `internal/onmyoji/{mode}/registr.go` implementing:
  - `LoadImageTemplates()` - reads PNG templates from `.//{mode}/` directory
  - `Registration()` - registers task sequence to state machine
- **Task files** defining `NamedTask` (e.g., `chapter_detector.go`)
- **Template images** in `./{mode}/` directory following existing naming convention

Example from k28:

```go
// internal/onmyoji/k28/registr.go
func (r Registrator) Registration(machine *statemachine.StateMachine, w window.Window, config onmyoji.Config, imgMap map[string]onmyoji.ImgInfo) error {
    onmyoji.Registration(machine, newChapterDetectorTask(config, w, imgMap[string(tasks.ZhangJie)]))
    // ... register more tasks in order
}
```

### 2. Task Implementation Pattern

All tasks implement:

```go
type Task interface {
    Name() tasks.TaskType
    Execute(controller TaskController) error
}
```

Common patterns:

- **Detection tasks**: Capture screenshot → template match → click if found → transition
- **Move tasks**: Execute predefined click sequences (see `k28/move.go`)
- **Retry logic**: Use config thresholds (e.g., `LevelCompletionThreshold`) before failing

### 3. Cross-Task Communication

Use `TaskController.SetAttribute()` / `GetAttribute()`:

```go
// Set in one task
controller.SetAttribute(tasks.BossFound, true)

// Retrieve in another
found, _ := controller.GetAttribute(tasks.BossFound)
```

Defined in `internal/tasks/type.go` as `TaskState` enum.

### 4. Configuration Management

- Config loaded from `./config.toml` via `LoadConfig()` in `internal/onmyoji/config.go`
- Creates default config if missing
- Mode-specific settings (K28Config, BreakerConfig) passed to task constructors
- Timing values in milliseconds (ChestWaitTime, LevelCompletionWaitTime)

### 5. Image Template Matching

- Templates loaded once at startup, reused across runs
- Filename → TaskState enum mapping (e.g., `ZhangJie.png` → `tasks.ZhangJie`)
- Missing templates cause panic in `LoadImageTemplates()` validation
- Match confidence determined by OpenCV's default threshold

## Build & Development Workflow

### Build Commands

```bash
make build          # Compile with OpenCV linking
make run            # Build + run with default k28 task
make test           # Run all tests with proper environment
make clean          # Remove binary
make package        # Create distribution zip
make deps           # Update Go modules
```

### Environment Requirements

- macOS with M-series chip (ARM64)
- OpenCV 4.11.0_1 via Homebrew: `brew install opencv`
- CGO configuration in Makefile handles include/lib paths
- BlueStacks emulator running with game open

### Running Custom Tasks

```bash
./free-hands-onmyoji -task k28 -display -1 -timeout 0
./free-hands-onmyoji -task breaker -display 1 -timeout 60 -closeBlueStacks
```

Task types: `k28`, `breaker`, `guren`, `limitedEvents` (see `onmyoji.GetModeNames()`)

## Critical Implementation Details

### Platform-Specific Code

- `internal/onmyoji/window/platform_darwin.go` - macOS-only (uses Mack for app activation)
- `internal/onmyoji/window/platform_windows.go` - Windows-only
- Select platform implementation in `init.go` via build tags or runtime detection

### Image Loading Strategy

- Opens image files in loops but must close handles (`defer file.Close()`)
- Uses `image.DecodeConfig()` for dimensions only (fast)
- Full image loaded separately via `utils.ReadPic()` for template matching
- Skips `.DS_Store` and invalid images without crashing

### Window Positioning

- Primary display: `platform.GetWindowPosition()`
- Secondary display: `platform.GetWindowPositionOnSecondDisplay(displayID)`
- Window coordinates passed to all click operations
- Timeout and graceful exit via global event listener (Cmd+Shift+O)

### Logging

- Structured logging via `internal/logger` using `go.uber.org/zap`
- Levels: Info (major flow), Debug (task transitions), Error (failures)
- All task transitions and template matches logged

## Common Debugging Scenarios

### Template Match Failing

1. Check image exists in mode directory and filename matches `TaskState` enum
2. Verify image is PNG, not corrupted
3. Game UI may have changed—capture fresh screenshot and update template

### Task Not Transitioning

1. Inspect state machine attribute usage in task
2. Verify `Next()` is called with correct `TaskType`
3. Check task order in `Registrator.Registration()`

### Window/Click Issues

1. Confirm BlueStacks active and game visible
2. Test `-display 1` for secondary monitor setup
3. Check macOS accessibility permissions for robotgo

## Adding Tests

- Unit tests use naming convention `*_test.go`
- CGO tests may require OpenCV library available
- Run via `make test` (handles environment setup)

## Key Files Reference

- [main.go](main.go) - Entry point, arg parsing, mode selection
- [internal/statemachine/state_machine.go](internal/statemachine/state_machine.go) - Core orchestration
- [internal/onmyoji/registrator.go](internal/onmyoji/registrator.go) - Task registration pattern
- [internal/onmyoji/window/platform.go](internal/onmyoji/window/platform.go) - Window abstraction
- [internal/tasks/type.go](internal/tasks/type.go) - TaskState enums
- [internal/onmyoji/k28/](internal/onmyoji/k28/) - Complete k28 mode example
- [Makefile](Makefile) - Build environment setup
