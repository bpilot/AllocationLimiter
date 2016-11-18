# AllocationLimiter
  Avoids unbounded allocation of memory by blocking until sufficient memory has been de-allocated.      Use when program is memory-bound to prevent rquests from causing heap size to grow until OOM.      Cooperative:   Relies on application correctly report when it would like to request memory, and return memory.   Typically used to limit a part of the program that creates a large allocation to control heap size by bounding it.     This is explicit memory management?