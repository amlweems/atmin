package atmin

import (
	"bytes"
	"log"
	"sync"
)

const (
	SetSteps       = 128
	TrimStartSteps = 16
	SetMinSize     = 4
)

type Minimizer struct {
	ex  Executor
	val Validator

	in  []byte
	out []byte
}

type Executor interface {
	Execute(in []byte) []byte
}

type Validator interface {
	Validate(initial, current []byte) bool
}

func NewMinimizer(in []byte) Minimizer {
	inDup := make([]byte, len(in))
	copy(inDup, in)

	return Minimizer{
		in: inDup,
	}
}

func nextP2(val int) int {
	ret := 1
	for val > ret {
		ret <<= 1
	}
	return ret
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func memset(a []byte, v byte) {
	if len(a) == 0 {
		return
	}
	a[0] = v
	for bp := 1; bp < len(a); bp *= 2 {
		copy(a[bp:], a[:bp])
	}
}

func (m Minimizer) Minimize() []byte {
	// minimization algorithm adapted from afl-tmin by lcamtuf
	// see: http://lcamtuf.coredump.cx/afl

	var changed bool
	var pass, inLen, setPos, setLen int
	var wg sync.WaitGroup
	var mu sync.Mutex

	// Stage 0: Block Normalization
	log.Print("Stage 0: Block Normalization")
	inLen = len(m.in)
	setLen = nextP2(inLen / SetSteps)
	if setLen < SetMinSize {
		setLen = SetMinSize
	}

	for setPos < inLen {
		wg.Add(1)
		go func(pos int) {
			defer wg.Done()

			tmpBuf := make([]byte, inLen)
			useLen := min(setLen, inLen-pos)

			mu.Lock()
			copy(tmpBuf, m.in)
			mu.Unlock()

			memset(tmpBuf[pos:pos+useLen], '0')

			out := m.ex.Execute(tmpBuf)
			res := m.val.Validate(m.out, out)

			if res {
				mu.Lock()
				memset(m.in[pos:pos+useLen], '0')
				mu.Unlock()
			}
		}(setPos)

		setPos += setLen
	}
	wg.Wait()

next_pass:
	var delLen, delPos, tailLen int
	var prevDel bool

	changed = false
	pass++

	// Stage 1: Block Deletion
	log.Print("Stage 1: Block Deletion")
	delLen = nextP2(inLen / TrimStartSteps)

next_del_blksize:
	if delLen <= 0 {
		delLen = 1

	}
	delPos = 0
	prevDel = true

	for delPos < inLen {
		tmpBuf := make([]byte, inLen)

		tailLen = inLen - delPos - delLen
		if tailLen < 0 {
			tailLen = 0
		}

		if delPos+delLen >= len(m.in) {
			break
		}

		/* If we have processed at least one full block (initially, prev_del == 1),
		   and we did so without deleting the previous one, and we aren't at the
		   very end of the buffer (tailLen > 0), and the current block is the same
		   as the previous one... skip this step as a no-op. */
		if !prevDel && tailLen > 0 && !bytes.Equal(m.in[delPos-delLen:delPos], m.in[delPos:delPos+delLen]) {
			delPos += delLen
			continue
		}

		prevDel = false

		// Head
		copy(tmpBuf[:delPos], m.in[:delPos])

		// Tail
		copy(tmpBuf[delPos:delPos+tailLen], m.in[delPos+delLen:delPos+delLen+tailLen])

		out := m.ex.Execute(tmpBuf[:delPos+tailLen])
		res := m.val.Validate(m.out, out)

		if res {
			copy(m.in[:delPos+tailLen], tmpBuf[:delPos+tailLen])
			prevDel = true
			inLen = delPos + tailLen
			changed = true
		} else {
			delPos += delLen
		}
	}

	if delLen > 1 && inLen >= 1 {
		delLen /= 2
		goto next_del_blksize
	}

	// Stage 2: Alphabet Minimization
	log.Print("Stage 2: Alphabet Minimization")
	alphaMap := make(map[byte]int)
	for i := 0; i < inLen; i++ {
		alphaMap[m.in[i]]++
	}

	for i := 0; i < 256; i++ {
		wg.Add(1)

		go func(pos byte) {
			defer wg.Done()

			if pos == '0' || alphaMap[pos] == 0 {
				return
			}

			tmpBuf := make([]byte, inLen)
			copy(tmpBuf[:inLen], m.in[:inLen])

			for r := 0; r < inLen; r++ {
				if tmpBuf[r] == pos {
					tmpBuf[r] = '0'
				}
			}

			out := m.ex.Execute(tmpBuf[:inLen])
			res := m.val.Validate(m.out, out)

			if res {
				mu.Lock()
				copy(m.in[:inLen], tmpBuf[:inLen])
				changed = true
				mu.Unlock()
			}
		}(byte(i))
	}
	wg.Wait()

	// Stage 3: Character Minimization
	log.Print("Stage 3: Character Minimization")

	for i := 0; i < inLen; i++ {
		wg.Add(1)

		go func(pos int) {
			defer wg.Done()

			tmpBuf := make([]byte, inLen)

			mu.Lock()
			copy(tmpBuf[:inLen], m.in[:inLen])
			mu.Unlock()

			if tmpBuf[pos] == '0' {
				return
			}
			tmpBuf[pos] = '0'

			out := m.ex.Execute(tmpBuf[:inLen])
			res := m.val.Validate(m.out, out)

			if res {
				mu.Lock()
				m.in[pos] = '0'
				changed = true
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()

	if changed {
		goto next_pass
	}

	return m.in[:inLen]
}
