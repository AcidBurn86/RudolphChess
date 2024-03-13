package main

import (
	"math/rand"
	"strconv"
	"sync"
	"time"

	"fmt"

	"bufio"
	"os"
	"strings"

	//import vendor notnil
	"github.com/notnil/chess"
)

var isInEco bool = true
var hide_output = false
var ecoLine string = ""

type ValidMoves struct {
	Moves []*chess.Move
}

type TestResult struct {
	ExpectedMove  string
	FoundMove     string
	IsCorrect     bool
	NodeCount     int
	ExpectedNodes int
	ElapsedTime   time.Duration
}

var moveCountInternal int = 0
var knightMoveValue int16 = 3
var bishopMoveValue int16 = 3
var rookMoveValue int16 = 4
var queenMoveValue int16 = 5

var (
	miMapa = make(map[string]int16)
	mu     sync.Mutex
)

var isReProcessBestLines bool = false

var game *chess.Game
var nodeCount int
var initEval [30]int16

func main() {

	cache = make(map[string]mobilityResult)

	bestCutMovesMap = make(map[string]*MostCutMoves)

	scanner := bufio.NewScanner(os.Stdin)

	//startTune()

	for scanner.Scan() {
		line := scanner.Text()
		tokens := strings.Fields(line)

		if len(tokens) == 0 {
			continue
		}

		switch tokens[0] {
		case "move":
			if len(tokens) > 1 {
				err := game.MoveStr(tokens[1])
				if err != nil {
					fmt.Println("Error al interpretar movimiento:", err)
				} else {
					game.MoveStr(tokens[1])
					if game.Position().Turn() == chess.Black {
						moveCountInternal++
					}
				}
			}

		case "uci":
			handleUCI()
		case "isready":
			handleIsReady()
		case "ucinewgame":
			handleUCINewGame()
		case "position5":
			handlePosition5(tokens)
		case "position6":
			handlePosition6(tokens)
		case "position":
			handlePosition(tokens)
		case "go":
			handleGo(tokens, "x", 0)
		case "numberOfMoves":
			handleNumberOfMoves()
		case "stop":
			//fmt.Println("bestmove", bestMove.Moves[0])
		case "d":
			fmt.Println(game.Position().Board().Draw())
			fmt.Println("moveCountInternal: ", moveCountInternal)
		case "dn":
			setupInitialBoard()
			printPrettyBoard()
			generateMoves()
		case "cut":
			hasToCut(game)
		case "getFen":
			fmt.Println(game.FEN())
		case "fen":
			setupFromFEN(strings.Join(tokens[1:], " "))
		case "fen5":
			setupFromFEN5(strings.Join(tokens[1:], " "))
		case "fen6":
			setupFromFEN6(strings.Join(tokens[1:], " "))
		case "perft":
			depth, _ := strconv.Atoi(tokens[1])
			fmt.Println("Perft:", Perft(game, depth))
		case "t":
			startTune()
		case "test":
			testsPos()
		case "qq":
			testsPos3()
		case "test2":
			testsPos2()
		case "li":
			testsPosLI()
		case "evalThis":
			evalThis = true
			fmt.Println("Eval:", EvaluateBoard(game.Position().Board(), game, game.Position().Turn() == chess.White, false, 1, 1))
			evalThis = false
		case "quit":
			return
		default:
			fmt.Println("Unknown command:", tokens[0])
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading input:", err)
	}

}

func Perft(game *chess.Game, depth int) int {
	if depth == 0 {
		return 1
	}

	var nodes int
	for _, move := range game.ValidMoves() {
		game2 := game.Clone()
		game2.Move(move)
		nodes += Perft(game2, depth-1)
	}

	return nodes
}

type reProcessBestLinesScores struct {
	score int16
	move  string
}

func handleGo(tokens []string, expectedMove string, expectedNodes int) TestResult {
	rand.Seed(time.Now().UnixNano())
	/*cpuProfileFile, err := os.Create("cpu.prof")
	if err != nil {
		panic(err)
	}
	defer cpuProfileFile.Close()

	// Start CPU profiling
	if err := pprof.StartCPUProfile(cpuProfileFile); err != nil {
		panic(err)
	}
	defer pprof.StopCPUProfile()*/

	/*memProfileFile, err := os.Create("mem.prof")
	if err != nil {
		panic(err)
	}
	defer memProfileFile.Close()

	// Write memory profile to file
	if err := pprof.WriteHeapProfile(memProfileFile); err != nil {
		panic(err)
	}
	fmt.Println("Memory profile written to mem.prof")*/

	//##############
	var newPositions []keyValue

	foundedHash = 0

	start := time.Now()
	cutLines = []string{}
	nodeCount = 0

	var timeLimit time.Duration = 0
	var score int16 = 0
	var sequence []string

	var reProcessBestLinesScoresVar []reProcessBestLinesScores

	var depth int8 = 1

	for i, token := range tokens {
		switch token {
		case "depth":
			depth = parseDepth(tokens[i+1])
		case "movetime":
			timeLimit = parseTime(tokens[i+1])

		case "wtime":
			timeLimit = parseTime(tokens[i+1])
		}
	}
	var isWhiteTurn = game.Position().Turn() == chess.White
	var isWhiteTurnBack = isWhiteTurn

	createBitBoard(game.Position().Board())

	countPawns := howManyPawns()
	countBigPieces := howManyBigPices()

	fmt.Println("countPawns: ", countPawns)
	fmt.Println("countBigPieces: ", countBigPieces)

	if countPawns <= 6 && areQueensInTheBoard() && countBigPieces <= 4 {
		depth = 6
	}

	if countPawns <= 4 && !areQueensInTheBoard() && countBigPieces <= 3 {
		depth = 7
	}

	if timeLimit == 0 {
		timeLimit = 5 * time.Second // 10 segundos
	}
	var resign bool = false
	var bestMove string = ""

	/*	fmt.Println("isInEco: ", isInEco)
		fmt.Println("ecoLine: ", ecoLine)
		fmt.Println("moveCountInternal: ", moveCountInternal)*/

	if isInEco && getNextMoveFromEco(ecoLine) != "" || moveCountInternal == 1 {

		bestMove = getNextMoveFromEco(ecoLine)
		//fmt.Println("Eco Best Move: ", bestMove)
	}
	if bestMove == "" {

		isInEco = false

		var mateFound bool = false
		cutedForDepth1 = []*chess.Move{}

		for i := 1; i <= int(depth); i++ {
			miMapa = make(map[string]int16)
			initEval[i] = EvaluateBoard(game.Position().Board(), game, isWhiteTurn, false, int8(i), int8(depth))

			if tryGoRoutines {
				score, sequence = SplitValidMoves(game, 1, -9999, 9999, isWhiteTurn, int8(i), []string{}, depth, true, false, false, false)
			} else {
				score, sequence = SearchBestMove(game, 1, -9999, 9999, isWhiteTurn, int8(i), []string{}, depth, true, false, false, false)
			}

			if len(sequence) > 0 {
				bestMovePrevDepth2 = bestMovePrevDepth
				bestMovePrevDepth = sequence[0]
			}
			sequenceStr := strings.Join(sequence, " ")
			if isWhiteTurn && score == 9999 || !isWhiteTurn && score == -9999 {
				mateFound = true
				break

			} else if !isWhiteTurn && score == 9999 || isWhiteTurn && score == -9999 {
				mateFound = true
				//resign = true

				fmt.Printf("info depth %d score cp %d pv %s nodes %d\n", depth, score, sequenceStr, nodeCount)
				//fmt.Println("resign")
				break
			}

			if !isWhiteTurnBack {
				score = score * -1
			}
			fmt.Printf("info depth %d score cp %d pv %s nodes %d\n", i, score, sequenceStr, nodeCount)

			AgregarElemento(sequenceStr, (score))
		}
		newPositions = reProcessBestLines(score, depth, isWhiteTurnBack)
		//fmt.Println("newPositions", newPositions)
		/*fmt.Println(newPositions[1].Key)
		fmt.Println(newPositions[1].Value)*/

		if !mateFound && depth > 22 {

			if len(newPositions) > 1 {
				for _, pos := range newPositions {
					isReProcessBestLines = true
					gameCloneReprocess := game.Clone()

					moves := strings.Split(pos.Key, " ")

					var isCapture bool = false

					//scoreLastMove := pos.Value

					for _, move := range moves[:len(moves)-1] {

						s2 := move[2:4]

						board := gameCloneReprocess.Position().Board()
						//	fmt.Println(move)
						pi := board.Piece(chess.Square(squareToIndex((s2))))
						gameCloneReprocess.MoveStr(move)

						//fmt.Println("Move: ", move)

						isCapture = false
						if pi != chess.NoPiece {
							//fmt.Println("Capture: ", move, pi)
							isCapture = true
							//scoreLastMove = pos.Value
						}

					}

					qscDepth := 6

					if !isCapture {
						//reProcessBestLinesScoresVar = append(reProcessBestLinesScoresVar, reProcessBestLinesScores{scoreLastMove, moves[0]})
						qscDepth = 6
					}

					alpha := int16(-9999)
					beta := int16(9999)
					if isWhiteTurn {
						alpha = score

					} else {
						beta = score
					}
					appendTo := true
					for i := 1; i <= int(qscDepth); i++ {
						isWhiteTurn = gameCloneReprocess.Position().Turn() == chess.White
						qs = true
						score, sequence = SearchBestMove(gameCloneReprocess, 1, alpha, beta, isWhiteTurn, int8(i), []string{}, depth, true, false, false, false)
						qs = false
						// append at the beggining fo sequence the moves
						sequence = append(moves, sequence...)
						appendTo = true
						if sequence != nil {
							sequenceStr := strings.Join(sequence, " ")
							if !hide_output && i == qscDepth {
								if score != 9999 && score != -9999 {
									fmt.Printf("info depth %d score cp %d pv %s nodes %d\n", int8(i), score, sequenceStr, nodeCount)
								} else {
									if isWhiteTurnBack && score == -9999 || !isWhiteTurnBack && score == 9999 {
										appendTo = false
									}
								}
							}

						}

						if isWhiteTurn && score == 9999 || !isWhiteTurn && score == -9999 {
							//	break
						} else if !isWhiteTurn && score == 9999 || isWhiteTurn && score == -9999 {
							//fmt.Println("resing")
							//resign = true
							//	break
						}
					}
					if appendTo {
						reProcessBestLinesScoresVar = append(reProcessBestLinesScoresVar, reProcessBestLinesScores{score, moves[0]})
					}
				}
			}
			//isReProcessBestLines = false

		}
		if isReProcessBestLines {

			//fmt.Println(reProcessBestLinesScoresVar)
			/*if depth == 6 && isWhiteTurnBack {
				isWhiteTurn = true
			}
			if depth == 6 && !isWhiteTurnBack {
				isWhiteTurn = false
			}
			if depth == 5 && isWhiteTurnBack {
				isWhiteTurn = true
			}
			if depth == 5 && !isWhiteTurnBack {
				isWhiteTurn = true
			}*/
			if isWhiteTurnBack {
				fmt.Println("isWhiteTurn")
				higherScore := int16(-9999)
				for _, score1 := range reProcessBestLinesScoresVar {
					if score1.score > higherScore {
						higherScore = score1.score
						bestMove = score1.move
						score = higherScore
					}
				}
			} else {
				fmt.Println("isBlackTurn")
				lowerScore := int16(9999)
				for _, score1 := range reProcessBestLinesScoresVar {
					if score1.score < lowerScore {
						lowerScore = score1.score
						bestMove = score1.move
						score = lowerScore
					}
				}
			}
		}
	}
	isReProcessBestLines = false
	elapsed := time.Since(start)
	fmt.Printf("Time taken: %s\n", elapsed)

	if len(sequence) > 0 {
		if sequence[0] != "" {
			if bestMove == "" {
				bestMove = sequence[0]
			}
		}
	}

	if bestMove == "" && len(newPositions) > 0 {
		bestMove = newPositions[0].Key[:5]
	}

	fmt.Println("foundedHash", foundedHash)

	if !resign {
		fmt.Println("bestScore", score)
		fmt.Println("bestmove", bestMove)
	}

	isCorrect := bestMove == expectedMove

	return TestResult{
		ExpectedMove:  expectedMove,
		FoundMove:     bestMove,
		IsCorrect:     isCorrect,
		NodeCount:     nodeCount,
		ExpectedNodes: expectedNodes,
		ElapsedTime:   elapsed,
	}

}

func parseTime(timeStr string) time.Duration {
	times, _ := strconv.Atoi(timeStr)

	return time.Duration(times) * time.Millisecond
}

func parseDepth(depthStr string) int8 {
	var depth int
	depth, _ = strconv.Atoi(depthStr)

	return int8(depth)
}

func handleUCINewGame() {
	positions = sync.Map{}
	cacheEvalCenterControl = make(map[gameTurnState]int16)
	cacheValidateCenterControl2 = make(map[gameTurnState]int16)
	memoizedResults = memoizeMap{cache: make(map[kingPositionKey][]int)}
	memoizedAttacks = memoizeAttacks{cache: make(map[attackKey]Bitboard)}
	cacheWhiteRookInOpenRank = make(map[boardState]bool)
	cacheBlackRookInOpenRank = make(map[boardState]bool)
	enemyCache = make(map[cacheKey]int)
	friendlyCache = make(map[cacheKey]int)
	nodeCount = 0
	cutedForDepth1 = []*chess.Move{}
	cutLines = []string{}
	bestCutMovesMap = make(map[string]*MostCutMoves)
	hashScores = sync.Map{}
	game = chess.NewGame(chess.UseNotation(chess.UCINotation{})) // Crear un nuevo juego
	handlePosition([]string{"startpos"})
	initBishopAttacks()
}

func handlePosition(tokens []string) {
	cacheEvalCenterControl = make(map[gameTurnState]int16)
	cacheValidateCenterControl2 = make(map[gameTurnState]int16)
	memoizedResults = memoizeMap{cache: make(map[kingPositionKey][]int)}
	memoizedAttacks = memoizeAttacks{cache: make(map[attackKey]Bitboard)}
	cacheWhiteRookInOpenRank = make(map[boardState]bool)
	cacheBlackRookInOpenRank = make(map[boardState]bool)
	enemyCache = make(map[cacheKey]int)
	friendlyCache = make(map[cacheKey]int)
	nodeCount = 0
	cutedForDepth1 = []*chess.Move{}
	cutLines = []string{}
	bestCutMovesMap = make(map[string]*MostCutMoves)
	hashScores = sync.Map{}
	//handleUCINewGame()
	var startIndex int // Índice para comenzar a aplicar los movimientos

	game = chess.NewGame(chess.UseNotation(chess.UCINotation{})) // Inicia un nuevo juego

OuterLoop:
	for i, token := range tokens {
		switch token {
		case "startpos":
			startIndex = i + 1 // Establecer startIndex al siguiente token después de "startpos"
		case "fen":
			fen := strings.Join(tokens[i+1:], " ")
			setupFromFEN(fen)
			startIndex = i + 2 // Establecer startIndex al siguiente token después del FEN
			break OuterLoop
		}
	}
	ecoLine = ""
	// Aplicar movimientos si existen después de "startpos" o "fen"
	if len(tokens) > startIndex && tokens[startIndex] == "moves" {

		moveCountInternal = 0
		for _, move := range tokens[startIndex+1:] {

			move = strings.TrimSpace(move) // Recorta los espacios antes de procesar el movimiento
			err := game.MoveStr(move)
			//	fmt.Println("Move:", move) // Eliminé el espacio extra después de "Move:" para mantener la consistencia
			if ecoLine == "" {
				ecoLine = move
			} else {
				ecoLine += " " + move // Utilicé el operador += para simplificar la concatenación
			}
			//		fmt.Println("ECO Line:", ecoLine) // Eliminé el espacio extra después de "ECO Line:" para mantener la consistencia
			moveCountInternal++
			if err != nil {
				return
			}
		}
	}

}

func handlePosition5(tokens []string) {
	hide_output = true
	var startIndex int // Índice para comenzar a aplicar los movimientos

	game = chess.NewGame(chess.UseNotation(chess.UCINotation{})) // Inicia un nuevo juego

OuterLoop:
	for i, token := range tokens {
		switch token {
		case "startpos":
			startIndex = i + 1 // Establecer startIndex al siguiente token después de "startpos"
		case "fen":
			fen := strings.Join(tokens[i+1:], " ")
			setupFromFEN(fen)
			startIndex = i + 2 // Establecer startIndex al siguiente token después del FEN
			break OuterLoop
		}
	}
	ecoLine = ""
	// Aplicar movimientos si existen después de "startpos" o "fen"
	if len(tokens) > startIndex && tokens[startIndex] == "moves" {

		moveCountInternal = 0
		for _, move := range tokens[startIndex+1:] {

			move = strings.TrimSpace(move) // Recorta los espacios antes de procesar el movimiento
			err := game.MoveStr(move)
			//	fmt.Println("Move:", move) // Eliminé el espacio extra después de "Move:" para mantener la consistencia
			if ecoLine == "" {
				ecoLine = move
			} else {
				ecoLine += " " + move // Utilicé el operador += para simplificar la concatenación
			}
			//		fmt.Println("ECO Line:", ecoLine) // Eliminé el espacio extra después de "ECO Line:" para mantener la consistencia
			moveCountInternal++
			if err != nil {
				return
			}
		}
	}

	handleGo([]string{"depth", "5"}, "x", 0)

}

func handlePosition6(tokens []string) {
	hide_output = true
	var startIndex int // Índice para comenzar a aplicar los movimientos

	game = chess.NewGame(chess.UseNotation(chess.UCINotation{})) // Inicia un nuevo juego

OuterLoop:
	for i, token := range tokens {
		switch token {
		case "startpos":
			startIndex = i + 1 // Establecer startIndex al siguiente token después de "startpos"
		case "fen":
			fen := strings.Join(tokens[i+1:], " ")
			setupFromFEN(fen)
			startIndex = i + 2 // Establecer startIndex al siguiente token después del FEN
			break OuterLoop
		}
	}
	ecoLine = ""
	// Aplicar movimientos si existen después de "startpos" o "fen"
	if len(tokens) > startIndex && tokens[startIndex] == "moves" {

		moveCountInternal = 0
		for _, move := range tokens[startIndex+1:] {

			move = strings.TrimSpace(move) // Recorta los espacios antes de procesar el movimiento
			err := game.MoveStr(move)
			//	fmt.Println("Move:", move) // Eliminé el espacio extra después de "Move:" para mantener la consistencia
			if ecoLine == "" {
				ecoLine = move
			} else {
				ecoLine += " " + move
			}
			moveCountInternal++
			if err != nil {
				return
			}
		}
	}

	handleGo([]string{"depth", "6"}, "x", 0)

}
func setupFromFEN(fen string) {
	handleUCINewGame()
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		fmt.Println("Error al interpretar FEN:", err)
		return
	}

	// replace in fen "\n" by nothing
	fen = strings.Replace(fen, "\n", "", -1)
	fen = strings.Replace(fen, " e3", "", -1)

	lastWord := strings.Split(fen, " ")
	num, _ := strconv.Atoi(lastWord[5])
	moveCountInternal = int(num)
	//fmt.Println("Move Count: ", moveCountInternal)
	// Aplicar la función FEN al juego.
	game = chess.NewGame(fenFunc, chess.UseNotation(chess.UCINotation{}))
}

func setupFromFEN5(fen string) {
	hide_output = true
	handleUCINewGame()
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		fmt.Println("Error al interpretar FEN:", err)
		return
	}

	lastWord := strings.Split(fen, " ")
	num, _ := strconv.Atoi(lastWord[5])
	moveCountInternal = int(num)

	game = chess.NewGame(fenFunc, chess.UseNotation(chess.UCINotation{}))

	handleGo([]string{"depth", "5"}, "x", 0)
}

func setupFromFEN6(fen string) {
	hide_output = true
	handleUCINewGame()
	fenFunc, err := chess.FEN(fen)
	if err != nil {
		fmt.Println("Error al interpretar FEN:", err)
		return
	}

	lastWord := strings.Split(fen, " ")
	num, _ := strconv.Atoi(lastWord[5])
	moveCountInternal = int(num)

	game = chess.NewGame(fenFunc, chess.UseNotation(chess.UCINotation{}))

	handleGo([]string{"depth", "6"}, "x", 0)
}

func handleIsReady() {
	fmt.Println("readyok")
}

func handleUCI() {
	fmt.Println("id name Rudolph")
	fmt.Println("id author MarcosPazzarelli")
	fmt.Println("id Country Argentina")

	fmt.Println("uciok")
}

func handleNumberOfMoves() {
	fmt.Println("numberOfMoves", len(game.MoveHistory()))
}
