# Classic Tetris Console Game Development Checklist

# ✅ **Phase 1: Core Game Data Structures**
- [ ] **Define Game Grid**
  - [ ] Create a **10x20 grid** (2D array).
  - [ ] Implement a function to **initialize the grid**.
  - [ ] Implement a function to **print the grid** to the console.
  - [ ] **Test**: Verify that an empty grid is displayed correctly.

- [ ] **Define Tetrominoes**
  - [ ] Create a **Tetromino struct** with:
    - [ ] Shape (2D array or coordinate offsets).
    - [ ] Position (row, column).
    - [ ] Rotation state (0, 90, 180, 270 degrees).
  - [ ] Define all **7 tetrominoes**: **I, O, T, L, J, S, Z**.
  - [ ] **Test**: Ensure each piece initializes with the correct shape and rotation states.

- [ ] **Implement Collision Detection**
  - [ ] Ensure tetrominoes **don’t move out of bounds**.
  - [ ] Ensure tetrominoes **don’t move through settled blocks**.
  - [ ] **Test**: Moving left, right, or down should be **blocked at edges or settled blocks**.

---

# ✅ **Phase 2: Gameplay Mechanics**
- [ ] **Implement Gravity (Auto Fall)**
  - [ ] Create a **game loop** that moves the tetromino **down every X milliseconds**.
  - [ ] Ensure the piece **stops when it reaches the bottom or another block**.
  - [ ] **Test**: A piece above the bottom should move **down each frame**.

- [ ] **Implement Line Clearing**
  - [ ] Detect when a row is **fully filled**.
  - [ ] Remove full rows and **shift all above rows downward**.
  - [ ] **Test**: Ensure line clearing **correctly removes rows and shifts blocks**.

- [ ] **Implement Piece Locking**
  - [ ] When a piece **lands**, lock it into the grid.
  - [ ] Spawn a **new piece at the top**.
  - [ ] **Test**: A landed piece should be **converted into a permanent block**.

- [ ] **Implement Game Over Condition**
  - [ ] If a new piece **cannot spawn**, trigger **Game Over**.
  - [ ] Display a **"Game Over" message**.
  - [ ] Ask if the player wants to **restart or quit**.
  - [ ] **Test**: Game should **end when the board is full**.

---

# ✅ **Phase 3: Player Input**
- [ ] **Implement Left & Right Movement**
  - [ ] Capture **left/right arrow keys**.
  - [ ] Move tetromino **left or right**.
  - [ ] Ensure movement **respects collision detection**.
  - [ ] **Test**: Movement should be **blocked at walls or other blocks**.

- [ ] **Implement Rotation**
  - [ ] Capture **Up Arrow key** for rotation.
  - [ ] Implement **clockwise rotation**.
  - [ ] Implement **wall kicks** to allow rotation near obstacles.
  - [ ] **Test**: Rotating near **walls or blocks should respect collision detection**.

- [ ] **Implement Soft Drop**
  - [ ] Capture **Down Arrow key**.
  - [ ] Increase **fall speed while held**.
  - [ ] **Test**: Soft drop should **move the piece down faster**.

- [ ] **Implement Hard Drop**
  - [ ] Capture **Spacebar key**.
  - [ ] Instantly **move piece to lowest possible position**.
  - [ ] Instantly **lock the piece in place**.
  - [ ] **Test**: Hard drop should **skip the lock delay and place piece immediately**.

---

# ✅ **Phase 4: Scoring & Leveling**
- [ ] **Implement Score Calculation**
  - [ ] Award **100 points for 1 line**.
  - [ ] Award **300 points for 2 lines**.
  - [ ] Award **500 points for 3 lines**.
  - [ ] Award **800 points for 4 lines (Tetris)**.
  - [ ] **Test**: Score updates correctly after line clears.

- [ ] **Implement Level Progression**
  - [ ] Increase **level every 10 lines** cleared.
  - [ ] Increase **fall speed** as level increases.
  - [ ] **Test**: Clearing **10 lines should increase the level**.

- [ ] **Display Score & Level**
  - [ ] Update the **console UI** with:
    - [ ] **Score**.
    - [ ] **Level**.
    - [ ] **Lines Cleared**.
  - [ ] **Test**: UI updates correctly after line clears.

---

# ✅ **Phase 5: Game Loop & UX Improvements**
- [ ] **Implement Pause Feature**
  - [ ] Capture **"P" key** for pause.
  - [ ] Display **"Paused" message**.
  - [ ] **Test**: Game should **freeze when paused**.

- [ ] **Implement Configurable Settings**
  - [ ] Allow **remapping keys**.
  - [ ] Allow selecting **Simple Rotation vs. SRS**.
  - [ ] **Test**: Config changes should **affect gameplay correctly**.

- [ ] **Implement High Score Saving**
  - [ ] Save **high scores to a file (`highscore.json`)**.
  - [ ] Load high scores on game start.
  - [ ] Display high score after **Game Over**.
  - [ ] **Test**: High scores should **persist between sessions**.

---

# ✅ **Phase 6: Final Refinements & Deployment**
- [ ] **Optimize Terminal Rendering**
  - [ ] Use `tcell` or `termbox-go` for **flicker-free rendering**.
  - [ ] **Update only changed parts** of the screen.

- [ ] **Ensure Cross-Platform Support**
  - [ ] Test on **Linux**.
  - [ ] Test on **macOS**.
  - [ ] Test on **Windows**.

- [ ] **Build & Package the Game**
  - [ ] Compile **standalone binaries (`go build`)**.
  - [ ] Provide a **README** with usage instructions.

- [ ] **Test Final Game**
  - [ ] Playtest **full game cycle**.
  - [ ] Ensure **no major bugs or crashes**.

---

### ✅ **Future Enhancements (Optional)**
- [ ] Add **Ghost Piece (shows where the tetromino will land)**.
- [ ] Implement **Combo Scoring** for consecutive line clears.
- [ ] Implement **Back-to-Back Tetris bonus**.
- [ ] Allow **customizable gravity/speed settings**.
- [ ] Add **multiplayer support** (competitive mode).

---
