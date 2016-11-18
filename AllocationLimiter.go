package AllocationLimiter

import (
  "sync"
)

/** @fileoverview 
  Avoids unbounded allocation of memory by blocking until sufficient memory has been de-allocated.
  
  Use when program is memory-bound to prevent rquests from causing heap size to grow until OOM.
  
  Cooperative:
  Relies on application correctly report when it would like to request memory, and return memory.
  Typically used to limit a part of the program that creates a large allocation to control heap size by bounding it. 

  This is explicit memory management?
*/

type AllocationLimiter struct {
  
  // Modification allocated is atomic
  mutex *sync.Mutex

  capacity_bytes int64 // @final
  // Initialized to capacity
  idle_bytes int64 // Acquire lock first before reading or writing this file.

  deallocNoticeChan chan bool
}

func (altr *AllocationLimiter) RequestAllocation(size_bytes int64) {

  if (size_bytes > altr.capacity_bytes) {
    panic("AllocationLimiter.RequestAllocation(): size_bytes cannot exceed capacity_bytes.")
  }

  retry_alloc:

  // Atomic operation
  altr.mutex.Lock()

  if ( altr.idle_bytes-size_bytes > 0) {
    altr.idle_bytes -= size_bytes
    altr.mutex.Unlock()
  } else {
    altr.mutex.Unlock()
    // Wait deallocation and try again.
    <- altr.deallocNoticeChan
    goto retry_alloc
  }
}

func (altr *AllocationLimiter) Deallocated(size_bytes int64) {
  altr.mutex.Lock() // Prevent clobbering by multiple threads calling RequestAllocation()
  
  altr.idle_bytes += size_bytes 

  altr.mutex.Unlock()

  go altr.sendDeallocatedToChan() // Do not block current thread
}

func (altr *AllocationLimiter) sendDeallocatedToChan() {
  altr.deallocNoticeChan <- false // give, do not block if no one is waiting.
}

// constructor
// @param capacity_bytes -- Represents how much memory can be used concurrently by this part of the program represented.
func New(capacity_bytes int64) *AllocationLimiter {
  return &AllocationLimiter{mutex: new(sync.Mutex),
                            deallocNoticeChan: make(chan bool, 0),
                            capacity_bytes: capacity_bytes,
                            idle_bytes: capacity_bytes}
}
