# pspterm — Soul

What this project is, what it values, and why it exists.

---

## The idea

The PSP XMB was a piece of UI design that punched well above its hardware. Running on a 333 MHz MIPS CPU behind a 4.3" LCD, it gave you silky horizontal navigation with spring-loaded bounce, a clock in the corner, dark navy gradients, and icons that dimmed as they receded. It felt premium when it had no right to.

pspterm takes that feeling and drops it into a terminal. Not a faithful pixel recreation — a faithful *feeling* recreation. The same bounce, the same dark, the same cold blue accent, the same clock in the corner. But underneath: a launcher. Navigate to a project directory. Open a shell. Launch vim. The XMB as productivity surface.

---

## What it values

**Feel over function.**
Any launcher can open a program. pspterm should feel like opening a program. The spring animation is not decoration — it communicates state. The overshoot tells you exactly where you are and where you came from.

**Dark and cold.**
The PSP color palette was deep navy against near-black. Not warm. Not colorful. High contrast where it mattered (selected item), invisible everywhere else. The UI should recede and let the content lead.

**Small and sharp.**
The codebase is intentionally small. Five packages. No framework beyond Bubble Tea. Every line exists because it needs to. Adding a feature should feel like carving, not piling on.

**The terminal is the platform.**
No Electron. No web view. No GUI framework. The PSP ran on constrained hardware and looked great anyway. The terminal is the constraint here — and the goal is to make something that looks great anyway.

**Respects the user's shell.**
The eval trick (printing a path to stdout, captured by a shell function) is the right primitive. pspterm does not try to be a shell. It navigates to a place and hands you back to your shell, which does the rest. Clean handoff.

---

## What it is not

Not a file manager. Not a process manager. Not a dashboard. Not a replacement for your shell.

It is a launcher with a very specific aesthetic and a very clean exit condition: you press enter on something, and either your environment changes (directory), a program runs (command), or a URL opens (browser). Then pspterm is done.

---

## The animation

The harmonica spring is the heart of the experience. The PSP did not snap categories — it flew them with momentum and let them settle. The damping ratio sits below 1.0 deliberately: underdamped, so there is a brief overshoot before resting. That overshoot is the soul of XMB navigation. Without it, this is just another fuzzy-finder with arrow keys.

If you tune the spring, aim to preserve three qualities:
1. **Immediate response** — the animation starts on keypress, no lag.
2. **Momentum** — the selected item clearly flew from somewhere.
3. **Settle** — it rests cleanly, no perpetual oscillation.

---

## On minimalism

Every removed feature is a gift. Waves: removed. Wave config fields: removed. The default config is short enough to read in thirty seconds. The source is small enough to understand in one sitting.

When adding something, ask: does this serve the core feeling, or does it add complexity in exchange for a feature nobody needed? Usually the answer is to leave it out.

---

## Lineage

- **PSP XMB** (2005) — the original. Sony's best UI work.
- **Bubble Tea** — the only TUI framework worth using for this kind of animation.
- **harmonica** — physics-based spring math, no fuss.
- **lipgloss** — styling without fighting the terminal.

The project stands on good shoulders. Keep it worthy of them.
