// Code reference: https://referencesource.microsoft.com/#mscorlib/system/random.cs

package csharp_random

import (
  "math"
)

const MSEED = 161803398

type CsRandom struct {
  seed      int32
  seedArray [56]int32
  inext     int32
  inextp    int32
}

func Random(seed int32) *CsRandom {
  r := &CsRandom{}
  r.init(seed)

  return r
}

func (self *CsRandom) init(seed int32) {
  var subtraction int32
  if seed == math.MinInt32 {
    subtraction = math.MaxInt32
  } else if seed < 0 {
    subtraction = - seed
  } else {
    subtraction = seed
  }

  var ii, mj, mk int32

  mj = MSEED - subtraction
  self.seedArray[55] = mj
  mk = 1
  for i:=1; i<55; i++ {
    ii = (int32)(21 * i) % 55
    self.seedArray[ii] = mk
    mk = mj - mk
    if mk < 0 {
      mk += math.MaxInt32
    }
    mj = self.seedArray[ii]
  }

  for k:=1; k<5; k++ {
    for i:=1; i<56; i++ {
      self.seedArray[i] -= self.seedArray[1 + (i + 30) % 55]
      if self.seedArray[i] < 0 {
        self.seedArray[i] = self.seedArray[i] + math.MaxInt32
      }
    }
  }

  self.inext = 0
  self.inextp = 21
  self.seed = 1
}

func (self *CsRandom) internalSample() int32 {
  var retVal int32
  var locINext int32 = self.inext + 1
  var locINextp int32 = self.inextp + 1

  if locINext >= 56 {
    locINext = 1
  }

  if locINextp >= 56 {
    locINextp = 1
  }

  retVal = self.seedArray[locINext] - self.seedArray[locINextp];

  if retVal == math.MaxInt32 {
    retVal--
  }
  if retVal < 0 {
    retVal += math.MaxInt32
  }

  self.seedArray[locINext] = retVal

  self.inext = locINext
  self.inextp = locINextp

  return retVal;
}

func (self *CsRandom) Sample() float64 {
  return float64(self.internalSample()) / float64(math.MaxInt32);
}

func (self *CsRandom) Next(minValue int32, maxValue int32) int32 {
  return (int32(self.Sample() * float64(maxValue - minValue)) + minValue);
}
