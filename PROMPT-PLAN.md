Here's a structured plan to implement the Classic Tetris Console Game step-by-step, ensuring safe, incremental progress with strong testing.

---

## **Phased Development Plan for Classic Tetris Console Game**

This plan follows a **layered, test-driven approach**, starting with the **core mechanics**, then building up **input handling, game logic, rendering, and UI**, ensuring each step is functional and testable before moving to the next.

---

## **Phase 1: Core Game Data Structures**
- **Goal**: Implement and test fundamental structures for the game grid, tetrominoes, and movement constraints.

### **Step 1.1: Define Game Grid**
- Create a **20x10 grid** stored as a 2D array.
- Implement a method to **print the grid** (for debugging).
- Write unit tests to confirm:
  - Grid initializes correctly.
  - Empty cells are displayed properly.

### **Step 1.2: Implement Tetromino Definitions**
- Define the **7 tetrominoes** using a standard representation.
- Store each tetromino in a struct with:
  - **Shape** (2D array or coordinate offsets).
  - **Current position** (row, column).
  - **Current rotation state** (0, 90, 180, 270 degrees).
- Write tests to ensure:
  - Tetrominoes are created correctly.
  - Rotation states are correctly stored.

### **Step 1.3: Implement Movement Constraints**
- Implement **collision detection**:
  - Ensure tetrominoes don’t move **out of bounds**.
  - Ensure tetrominoes don’t move **through settled blocks**.
- Write tests:
  - Moving left, right, down should be **blocked** at grid edges.
  - Moving into occupied cells should be **disallowed**.

---

## **Phase 2: Gameplay Mechanics**
- **Goal**: Implement core mechanics, including gravity, line clearing, and piece locking.

### **Step 2.1: Implement Gravity (Auto Fall)**
- Create a **game loop** that moves the tetromino down every **X milliseconds**.
- Ensure the piece stops when it reaches the **bottom** or lands on another piece.
- Write tests:
  - A piece placed above the bottom should move **down each frame**.
  - A piece reaching the bottom should **not move further**.

### **Step 2.2: Implement Line Clearing**
- Check if a row is **fully filled**.
- If full, remove the row and **shift all above rows downward**.
- Write tests:
  - A row filled with blocks should be **cleared**.
  - Rows above the cleared row should **shift down correctly**.

### **Step 2.3: Implement Locking & Next Piece**
- Once a piece **lands**, it should be **locked** into the grid.
- Spawn the **next piece** at the top.
- If the spawn position is **occupied**, trigger **Game Over**.
- Write tests:
  - A landed piece should be **converted into a permanent block**.
  - A new piece should **appear after locking**.

---

## **Phase 3: Player Input**
- **Goal**: Implement real-time player movement (left, right, soft drop, rotate, hard drop).

### **Step 3.1: Implement Left & Right Movement**
- Capture **arrow key presses** for **left** and **right**.
- Ensure movement respects **collision detection**.
- Write tests:
  - Moving left/right should be **blocked at walls**.
  - Moving left/right should be **blocked by other blocks**.

### **Step 3.2: Implement Rotation**
- Add **rotation support** (default: clockwise).
- Implement **simple rotation first**, then add **SRS kicks** later.
- Write tests:
  - Rotating near **walls** should prevent overlap.
  - Rotating should **change the piece orientation correctly**.

### **Step 3.3: Implement Soft Drop & Hard Drop**
- **Soft Drop**: Increase **fall speed** while holding **Down**.
- **Hard Drop**: Move piece to **lowest possible position** and **lock instantly**.
- Write tests:
  - Soft drop should **speed up gravity**.
  - Hard drop should **instantly place the piece**.

---

## **Phase 4: Scoring & Leveling**
- **Goal**: Implement score tracking, level progression, and display.

### **Step 4.1: Implement Scoring System**
- Award **points for line clears**:
  - 1 line = **100** points.
  - 2 lines = **300** points.
  - 3 lines = **500** points.
  - 4 lines = **800** points.
- Write tests:
  - Clearing lines should **increase the score**.
  - Multiple line clears should **apply the correct bonus**.

### **Step 4.2: Implement Level Progression**
- Increase **level every 10 lines** cleared.
- Increase **fall speed** as level increases.
- Write tests:
  - Clearing **10 lines** should **increase level**.
  - Higher levels should **increase gravity speed**.

### **Step 4.3: Display Score, Level, and Lines Cleared**
- Update the **console UI** to show:
  ```
  Score: 1200  Level: 3  Lines: 25
  ```
- Ensure the **score updates correctly** after each action.

---

## **Phase 5: Game Loop & UX Improvements**
- **Goal**: Polish the game loop, adding pause, game over handling, and settings.

### **Step 5.1: Implement Game Over Handling**
- If a new piece **can’t spawn**, display **"Game Over"**.
- Ask the player if they want to **restart**.
- Write tests:
  - Game should **end if the board fills up**.
  - Game should **reset when restarted**.

### **Step 5.2: Implement Pause Feature**
- Allow pausing/resuming using the **"P" key**.
- Display **"Paused"** when the game is paused.

### **Step 5.3: Add Configurable Settings**
- Allow **remapping keys**.
- Store **high scores in a file** (`highscore.json`).
- Allow toggling between **Simple Rotation and SRS**.
- Write tests:
  - High scores should **persist between sessions**.
  - Changing settings should **affect gameplay correctly**.

---

## **Phase 6: Final Refinements & Deployment**
- **Goal**: Ensure smooth terminal rendering, optimize performance, and prepare for distribution.

### **Step 6.1: Optimize Terminal Rendering**
- Use `tcell` or `termbox-go` for **flicker-free rendering**.
- Instead of **redrawing everything**, update **only changed parts**.

### **Step 6.2: Ensure Cross-Platform Support**
- Test on **Linux**, **macOS**, and **Windows**.
- Fix any **platform-specific rendering issues**.

### **Step 6.3: Package & Release**
- Build standalone binaries (`go build`).
- Provide a **README** with usage instructions.
- Distribute as a **single binary** (`tetris.exe`, `tetris`).

---

## **LLM Code Generation Prompts**
- Each of the above steps can be **implemented using LLM-generated code**.
- Below are **sample prompts** for each **incremental implementation**.

---

### **Prompt 1: Implement Game Grid**
```text
Write a Go function that initializes a 20x10 grid as a 2D array. Each cell should start as empty. Add a method to print the grid to the console in a structured format.
```

---

### **Prompt 2: Define Tetrominoes**
```text
Create a struct in Go to represent a Tetris tetromino. It should include:
- A shape (2D array or coordinate offsets).
- A position (row, column).
- A rotation state (0, 90, 180, 270 degrees).
Define all 7 classic Tetris tetrominoes using this struct.
```

---

### **Prompt 3: Implement Collision Detection**
```text
Write a Go function that checks if a tetromino's position is valid within the grid. The function should return false if:
- The piece is out of bounds.
- The piece overlaps with a settled block.
Write unit tests to verify correct behavior.
```

---

### **Prompt 4: Implement Line Clearing**
```text
Write a Go function that scans the grid for full rows. If a row is full, remove it and shift all rows above down. Ensure the function updates the grid correctly.
```

---

### **Final Prompt**
```text
Integrate all previous steps into a complete, playable Tetris game in the terminal. Ensure all input handling, game loop mechanics, scoring, and rendering work smoothly.
```

---

This **incremental, test-driven approach** ensures a **stable**, **maintainable**, and **fully functional** Tetris game. 🚀