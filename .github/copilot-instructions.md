Your goal is to implement an automatic scripting tool for turn-based games
The project uses the following technology stacks:
- Golang
- OpenCV
- robotgo
There are a few things to consider when outputting your code but don't overthink it:
- Consider loading time to switch to the next screen when the game moves to the next session
- The game may stall and cause unsuccessful clicks
- Appropriate delays need to be added to the code to ensure that the game has enough time to load the new screen
- You need to handle stalling situations, such as retrying clicks
- Only consider running on macOS M series
