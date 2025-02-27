package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"log"
	"math/rand"
	"os"
	"strings"
)

const (
	// Define the size of our playing field, add 2 for the borders.
	tableWidth  = 60 + 2
	tableHeight = 30 + 2

	// Define our symbols for representing the playing field.
	corner         = '+'
	lineVertical   = '|'
	lineHorizontal = '-'
	empty          = ' '
	player         = 'M'
	item           = '*'
	enemy          = 'T'
)

// model encapsulates our data for displaying and updating.
type model struct {
	table    [tableHeight][tableWidth]rune
	turtles  map[int]*object
	mario    *object
	mushroom *object
	//playerRow int
	//playerCol int
	moveCount int
	score     int
	gameOver  bool
}

type object struct {
	YY, XX int
}

func (m *model) spawnItem() {
	yy, xx := m.randomFreeCoordinates()
	m.table[yy][xx] = item
	m.mushroom = &object{YY: yy, XX: xx}
}

func (m *model) spawnEnemy() {
	yy, xx := m.randomFreeCoordinates()
	index := len(m.turtles) + 1
	m.turtles[index] = &object{YY: yy, XX: xx}
	m.table[yy][xx] = enemy
}

func (m *model) moveEnemies() {
	for _, item := range m.turtles {
		direction := rand.Intn(4)
		moveYY := item.YY
		moveXX := item.XX
		switch direction {
		case 1:
			moveYY = checkDown(moveYY) // Down. 1
		case 2:
			moveXX = checkLeft(moveXX) // Left. 0
		case 3:
			moveXX = checkRight(moveYY) // Right. 0
		default:
			moveYY = checkUp(moveXX) // Top. 0
		}
		m.table[item.YY][item.XX] = empty
		item.YY = moveYY
		item.XX = moveXX
		m.table[moveYY][moveXX] = enemy

		if m.mario.YY == moveYY && m.mario.XX == moveXX {
			m.gameOver = true
			return
		}
	}

	//for row := 1; row < tableHeight-1; row++ {
	//	for col := 1; col < tableWidth-1; col++ {
	//
	//		if m.table[row][col] == enemy {
	//			// Move enemy randomly to one of the four directly neighbooring cells if it is empty
	//			// or contains the player. Borders, items and other enemies will block enemy moves though.
	//			neighboors := [4][2]int{
	//				{row - 1, col}, // Top neighboor.
	//				{row, col + 1}, // Right neighboor.
	//				{row + 1, col}, // Bottom neighboor.
	//				{row, col - 1}, // Right neighboor.
	//			}
	//
	//			targetIndex := rand.Intn(len(neighboors))
	//			targetRow := neighboors[targetIndex][0]
	//			targetCol := neighboors[targetIndex][1]
	//			log.Println(row, col, targetRow, targetCol)
	//
	//			if m.table[targetRow][targetCol] == enemy {
	//				log.Println(targetRow, targetCol)
	//			} else if m.table[targetRow][targetCol] == empty {
	//				// Target cell is empty. Move enemy and clear the old one.
	//				m.table[targetRow][targetCol] = enemy
	//				m.table[row][col] = empty
	//			} else if m.table[targetRow][targetCol] == player {
	//				// Target cell contains the player. Attack and stop further processing, this game is over!
	//				m.gameOver = true
	//
	//				return
	//			}
	//		}
	//	}
	//}
}

func (m *model) randomFreeCoordinates() (row, col int) {
	// Generate some random coordinates.
	row, col = randomCoordinates()

	// Check that the random cell is empty.
	// If not repeat randomizing until we find an empty cell.
	for m.table[row][col] != empty {
		row, col = randomCoordinates()
	}
	return
}

// randomCoordinates returns a new set of random coordinates within the
// playing field excluding borders. However, it is not guaranteed that the
// cell under the returned coordinates is actually empty.
func randomCoordinates() (row, col int) {
	row = rand.Intn(tableHeight-2) + 1
	col = rand.Intn(tableWidth-2) + 1
	return
}

func (m *model) movePlayer(yy, xx int) {
	if m.gameOver {
		return
	}

	// Clear old player location.
	m.table[m.mario.YY][m.mario.XX] = empty

	if m.table[yy][xx] == enemy {
		// We ran into an enemy. Signal game over and skip further
		// processing.
		m.gameOver = true
		return
	}

	if m.table[yy][xx] == item {
		// We collected an item. A new item and enemy needs to be
		// spawned. Increase the score.
		m.score++
		m.spawnItem()
		m.spawnEnemy()
	}

	// Set new player location.
	m.table[yy][xx] = player
	m.mario.YY = yy
	m.mario.XX = xx
}

func checkUp(yy int) int {
	if yy <= 1 {
		return yy
	}
	return yy - 1
}

func checkDown(yy int) int {
	if yy >= tableHeight-2 {
		return yy
	}
	return yy + 1
}

func checkLeft(xx int) int {
	if xx <= 1 {
		return xx
	}
	return xx - 1
}

func checkRight(xx int) int {
	if xx >= tableWidth-2 {
		return xx
	}
	return xx + 1
}

func (m *model) playerUp() {
	newYY := checkUp(m.mario.YY)
	if m.mario.YY == newYY {
		// Do nothing as we are already at the border and cannot move.
		return
	}
	m.moveCount++
	m.movePlayer(newYY, m.mario.XX)
}

func (m *model) playerDown() {
	newYY := checkDown(m.mario.YY)
	if m.mario.YY == newYY {
		// Do nothing as we are already at the border and cannot move.
		return
	}
	m.moveCount++
	m.movePlayer(newYY, m.mario.XX)
}

func (m *model) playerLeft() {
	newXX := checkLeft(m.mario.XX)
	log.Println("checkLeft", m.mario.XX, newXX)
	if m.mario.XX == newXX {
		// Do nothing as we are already at the border and cannot move.
		return
	}
	m.moveCount++
	m.movePlayer(m.mario.YY, newXX)
}

func (m *model) playerRight() {
	newXX := checkRight(m.mario.XX)
	if m.mario.XX == newXX {
		// Do nothing as we are already at the border and cannot move.
		return
	}
	m.moveCount++
	m.movePlayer(m.mario.YY, newXX)
}

// init is responsible for initializing or resetting a model that is ready
// to use for a new game.
func (m *model) init() {
	// Clear and reset all fields as init() is also used for restarting
	// an existing game. Therefore, our model needs to be fresh.

	// Initially, set every cell to our empty symbol.
	for row := 0; row < tableHeight; row++ {
		for col := 0; col < tableWidth; col++ {
			m.table[row][col] = empty
		}
	}
	m.mario = &object{YY: 1, XX: 1}
	m.score = 0
	m.moveCount = 0
	m.gameOver = false
	m.turtles = make(map[int]*object)
	// Set the four corners.
	m.table[0][0] = corner
	m.table[0][tableWidth-1] = corner
	m.table[tableHeight-1][0] = corner
	m.table[tableHeight-1][tableWidth-1] = corner

	// Draw horizontal borders at the top and bottom.
	for col := 1; col < tableWidth-1; col++ {
		m.table[0][col] = lineHorizontal
		m.table[tableHeight-1][col] = lineHorizontal
	}

	// Draw vertical borders on the left and right side.
	for row := 1; row < tableHeight-1; row++ {
		m.table[row][0] = lineVertical
		m.table[row][tableWidth-1] = lineVertical
	}

	// Spawn our player near the top left corner.
	m.table[m.mario.YY][m.mario.XX] = player

	m.spawnItem()
	m.spawnEnemy()
}

// Init can be used to setup initial command to perform.
// We don't need anything here. Therefore we return nil.
func (m *model) Init() tea.Cmd {
	return nil
}

// Update is called whenever something happens like a key is pressed or
// another event occurs. Then, we have the option of reacting to it by
// modifying our model.
func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:

		// If our game is lost, any key shall restart the game.
		if m.gameOver {

			switch msg.String() {

			case "ctrl+c", "q":
				return m, tea.Quit
			case "enter":
				m.init()
				return m, nil
			}
		}

		switch msg.String() {

		// Exit program on ctrl+c or q typing.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			m.playerUp()
			m.moveEnemies()
		case "down":
			m.playerDown()
			m.moveEnemies()
		case "left":
			m.playerLeft()
			m.moveEnemies()
		case "right":
			m.playerRight()
			m.moveEnemies()
		}
	}

	return m, nil
}

// View is required for building what we want to show on the screen.
// That means we need to translate our model data into a string for displaying.
func (m *model) View() string {
	builder := strings.Builder{}

	if m.gameOver {
		// Just inform about game over and don't continue.
		builder.WriteString("\n\n\n\n\n")
		builder.WriteString("          You died, Game Over!")
		builder.WriteString("\n\n")
		builder.WriteString(fmt.Sprintf("          Your score: %d | Your moved: %d \n", m.score, m.moveCount))
		builder.WriteString("\n\n")
		builder.WriteString("          Press enter to restart or q to quit")

		return builder.String()
	}

	builder.WriteString(fmt.Sprintf("   Your score: %d | Your moved: %d \n", m.score, m.moveCount))
	// Iterate our table (2d array) and print the cells.
	var enemyYY = 0
	var enemyXX = 0
	for y, yy := range m.table {
		for x, xx := range yy {
			if xx == enemy {
				enemyYY = y
				enemyXX = x
			}
			builder.WriteRune(xx)
		}
		// Go to next line after each row.
		builder.WriteString("\n")
	}
	builder.WriteString(fmt.Sprintf("position :  %d ,  %d \n", enemyYY, enemyXX))
	builder.WriteString(fmt.Sprintf("position :  %d ,  %d", m.mario.YY, m.mario.XX))
	return builder.String()
}

func main() {
	// Create our initial model.
	model := &model{}
	model.init()

	// Program setup to initialize bubbletea and use full screen.
	program := tea.NewProgram(model, tea.WithAltScreen())

	// Run bubbletea and exist with a message if an error occurs.
	if _, err := program.Run(); err != nil {
		fmt.Println("Unexpected error: %v", err)
		os.Exit(1)
	}
}
