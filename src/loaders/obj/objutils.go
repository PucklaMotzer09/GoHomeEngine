package loader

import (
	"bufio"
	"strconv"
)

func readLine(rd *bufio.Reader) (string, error) {
	var str string
	var isPrefix = true
	var err error
	var buf []byte

	for err == nil && isPrefix {
		buf, isPrefix, err = rd.ReadLine()
		str += string(buf)
	}

	return str, err
}

func toTokens(line string) []string {
	var curByte byte
	var readToken bool = false
	var tokens []string

	for i := 0; i < len(line); i++ {
		curByte = line[i]
		if curByte == ' ' {
			readToken = false
		} else {
			if readToken {
				tokens[len(tokens)-1] += string(curByte)
			} else {
				tokens = append(tokens, "")
				tokens[len(tokens)-1] += string(curByte)
				readToken = true
			}
		}
	}

	return tokens
}

func processFaceData1(elements [][]string) (rv []int) {
	rv = make([]int, len(elements))
	for i := 0; i < len(rv); i++ {
		temp, _ := strconv.ParseInt(elements[i][0], 10, 32)
		rv[i] = int(temp)
	}
	return
}

func processFaceData2(elements [][]string) (pos []int, tex []int) {
	rv := make([]int, len(elements)*2)
	pos = make([]int, len(elements))
	tex = make([]int, len(elements))
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			temp, _ := strconv.ParseInt(elements[i][j], 10, 32)
			rv[i*2+j] = int(temp)
		}
	}

	for i := 0; i < len(elements); i++ {
		pos[i] = rv[i*2]
		tex[i] = rv[i*2+1]
	}

	return
}

func processFaceData3(elements [][]string) (pos []int, norm []int) {
	rv := make([]int, len(elements)*2)
	pos = make([]int, len(elements))
	norm = make([]int, len(elements))
	var readIndex int
	for i := 0; i < 3; i++ {
		for j := 0; j < 2; j++ {
			if j == 1 {
				readIndex = 2
			} else {
				readIndex = j
			}
			temp, _ := strconv.ParseInt(elements[i][readIndex], 10, 32)
			rv[i*2+j] = int(temp)
		}
	}

	for i := 0; i < len(elements); i++ {
		pos[i] = rv[i*2]
		norm[i] = rv[i*2+1]
	}

	return
}

func processFaceData4(elements [][]string) (pos []int, tex []int, norm []int) {
	rv := make([]int, len(elements)*3)
	pos = make([]int, len(elements))
	tex = make([]int, len(elements))
	norm = make([]int, len(elements))
	for i := 0; i < len(elements); i++ {
		for j := 0; j < 3; j++ {
			temp, _ := strconv.ParseInt(elements[i][j], 10, 32)
			rv[i*3+j] = int(temp)
		}
	}

	for i := 0; i < len(elements); i++ {
		pos[i] = rv[i*3]
		tex[i] = rv[i*3+1]
		norm[i] = rv[i*3+2]
	}

	return
}

func process3Floats(tokens []string) [3]float32 {
	var rv [3]float32
	var temp float64
	var err error

	for i := 0; i < 3; i++ {
		temp, err = strconv.ParseFloat(tokens[i], 32)
		if err != nil {
			return [3]float32{0.0, 0.0, 0.0}
		}
		rv[i] = float32(temp)
	}

	return rv
}

func process2Floats(tokens []string) [2]float32 {
	var rv [2]float32
	var temp float64
	var err error

	for i := 0; i < 2; i++ {
		temp, err = strconv.ParseFloat(tokens[i], 32)
		if err != nil {
			return [2]float32{0.0, 0.0}
		}
		rv[i] = float32(temp)
	}

	return rv
}

func process1Float(tokens string) float32 {
	var rv float32
	var temp float64
	var err error

	temp, err = strconv.ParseFloat(tokens, 32)
	if err != nil {
		return 0.0
	}
	rv = float32(temp)

	return rv
}

func addAllTokens(tokens []string, start int) (str string) {
	for i := start; i < len(tokens); i++ {
		if i == start {
			str += " "
		}
		str += tokens[i]
	}
	return
}
