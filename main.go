package main // Package main, Do not change this line.

import (
	"context"
	"fmt"
	"math/rand/v2"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// Product represents the structure for a product entity

// SupplyChainContract defines the smart contract structure
type SupplyChainContract struct {
	// To Do
}

func main() {

	var prng_seed = 0
	var num_of_iterations = 0
	var input_prog string
	var output_file string
	var seedfile string
	if len(os.Args) < 5 {
		fmt.Println("Run the fuzzer with args :: 154890100 <num of iterations> <target prog path> <seed file> <crash output file>")
		os.Exit(1)
	}
	prng_seed, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Issue with inputs")
	}
	num_of_iterations, err = strconv.Atoi(os.Args[2])
	if err != nil {
		fmt.Printf("Issue with inputs")
	}
	input_prog = os.Args[3]
	seedfile = os.Args[4]
	output_file = os.Args[5]
	// fmt.Println("prng seed", prng_seed)
	// fmt.Println("iteration times", num_of_iterations)
	// fmt.Println("input prog ", input_prog)
	// fmt.Println("output prog ", output_file)

	inpbytes, err := os.ReadFile(seedfile)

	// inpbytes, _ := hex.DecodeString(string(inp))

	if err != nil {
		fmt.Printf("Seed File cannot be read Check permissions")
	}

	random := rand.NewPCG(uint64(prng_seed), uint64(prng_seed))
	randomizer := rand.New(random)

Iteration:
	for i := 0; i < num_of_iterations; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		// randomstr(rand, inpbytes)
		nc := exec.CommandContext(ctx, input_prog)
		nc.Stdin = strings.NewReader(string(inpbytes))
		err = nc.Start()
		// out, err = nc.CombinedOutput()
		if err != nil {
			fmt.Println("Error invoking target prog:", err)
		}
		// fmt.Println("output from prog", out)

		done := make(chan error, 1)
		go func() {
			done <- nc.Wait()
		}()

		select {
		case err := <-done:
			// fmt.Println(err)
			// fmt.Println(nc.Stderr)
			// fmt.Println(nc.Stdout)
			if err != nil {
				exitErr, ok := err.(*exec.ExitError)
				if ok {
					status, ok := exitErr.Sys().(syscall.WaitStatus)
					if ok {
						if status.Signaled() {
							// fmt.Println("Process exited due to error signal:", status.Signal())
							// fmt.Println("Status:", status)
							// fmt.Println(inpbytes)
							// os.MkdirAll(output_file, 0777)
							os.WriteFile(output_file, inpbytes, 0777)
							// fmt.Println("Iteration num-->", i)
							break Iteration
						}
					}
				}

			}
		case <-ctx.Done():
			// fmt.Println("Timeout, killing process...")
			cancel()
			<-done
		}

		if i != 0 && i%500 == 0 {
			for j := 0; j < 10; j++ {
				inpbytes = append(inpbytes, byte(rand.Int32N(256)))
			}
			// fmt.Printf("Appended Bytes ::: %s\n", hex.EncodeToString(f))

		}
		randomstr(randomizer, inpbytes)
	}

}

func randomstr(randomizer *rand.Rand, f []byte) []byte {
	flen := len(f)
	for i := 0; i < flen; i += 8 {
		j := randomizer.IntN(8)
		if (i + j) < flen {
			nbyte := byte(randomizer.Int32N(256))
			f[i+j] = nbyte
		}
	}
	return f

}
