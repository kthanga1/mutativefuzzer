#!/usr/bin/env python3
import subprocess

prng_seed = 154890100
num_iterations = 100
target program path =  "./prog"
seed =  _seed__
outputpath =  level-1.crash


def executescript(scrpath, args):
    process = subprocess.Popen([scrpath, args], stdout=subprocess.PIPE, stderr=subprocess.PIPE, text=True)
    out, err = process.communicate()
    return process.returncode, out, err




if __name__ == '__main__':

    for i in range(10):
        print("Testing the fuzzer for level 1")
        testout = "./testlevel_"+i+"/level"+i+".crash"
        status, out, err = executescript("./fuzzy/fuzzy",154890100, 100, "./prog", "./seed/_seed_", testout)

        if status == 1:
            print("Process executed for level"+i)
            resultpath = "./results/level-"+i+"/level-"+i+".crash"
            status, out, err = executescript("diff", resultpath , testout)
            print(status)
            print(out)
            print(err)
        break

