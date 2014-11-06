package bloom

import (
	"fmt"
	"github.com/mtchavez/jenkins"
	"hash/fnv"
	"math"
)

type BloomFilter interface {
	Add(string)
	Test(string) bool
	Count() uint64
	FalsePositives() float64
}

/*
k : number of hashing functions
m : number of bits in the filter
n : number of elements inserted
*/

const (
	UINT64_SIZE = uint32(64)
)

type simpleBloomFilter struct {
	count uint64
	arr   []uint64
	bits  uint32
}

func NewSimpleBloomFilter(estimatedSize uint32, err float64) BloomFilter {

	return &simpleBloomFilter{
		count: 0,
		arr:   make([]uint64, size),
		bits:  UINT64_SIZE * size,
	}
}

func (t *simpleBloomFilter) Add(str string) {
	h1 := jenkinsHash(str)
	h2 := fnvHash(str)

	for i := 0; i < len(t.arr); i++ {
		x := h1 + (uint32(i)*h2)%UINT64_SIZE
		t.arr[x/UINT64_SIZE] |= 1 << (x % UINT64_SIZE)
	}
	t.count++
}

func (t *simpleBloomFilter) Test(str string) bool {
	h1 := jenkinsHash(str)
	h2 := fnvHash(str)

	for i := 0; i < len(t.arr); i++ {
		x := h1 + (uint32(i)*h2)%UINT64_SIZE
		if (t.arr[x/UINT64_SIZE] & (1 << (x % UINT64_SIZE))) == 0 {
			return false
		}
	}
	return true
}

func (t simpleBloomFilter) Count() uint64 {
	return t.count
}

func (t simpleBloomFilter) String() {
	for _, bits := range t.arr {
		fmt.Printf("%064b ", bits)
	}
	fmt.Println()
}

// x =  (1 - e^(-kn/m))^k
// ln(x) = k ln(1 - e^(-kn/m))
// ln(x) = k ln(1)ln(e^(-kn/m))
// ln(x) / k
// (1 - x^(1/k)) = e^(-kn/m)
//
func (t *simpleBloomFilter) FalsePositives() float64 {
	k := float64(len(t.arr)) // num of hashes
	m := float64(t.bits)     // num of bloom bits
	n := float64(t.count)
	return math.Pow(1-math.Exp(-k*n/m), k)
}

func jenkinsHash(str string) uint32 {
	j := jenkins.New()
	if _, err := j.Write([]byte(str)); err != nil {
		panic(err)
	}
	return j.Sum32() % UINT64_SIZE
}

func fnvHash(str string) uint32 {
	f := fnv.New32()
	if _, err := f.Write([]byte(str)); err != nil {
		panic(err)
	}
	return f.Sum32() % UINT64_SIZE
}

/*
 m = k * p (size)
 g_i(x) = h_1(x) + i * h_2(x) % p

 n: num elements
 k: num hash functions
 p: number bits
 i: [0, k-1]
 h_1(x)
 h_2(x)
 hash range = [0, p-1]

 False Positive

 P(Fp) = 1 - (1 k^2/m^2)^n

 m/n = c

 P(Fp) = (1 - exp^(-k/c))^k


  n -> ∞
  p: m/k is prime
  m: c * n

Double Hashing
  Array
  m bits paritioned into k disjoint arrays
  of m' = m/k bits (m%k==0)

  h_1(u) + i * h_2(u) mod m'

  P(Fp) = (1 - e^-(k/c))^k


Extended Double Hashing

  h_1(u) + i * h_2(u) + f(i) mod m
*/
