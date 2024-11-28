use "collections"
use "time"

actor Main
  new create(env: Env) =>
    if env.args.size() != 3 then
      env.out.print("Usage: lukas <N> <k>")
      return
    end

    try
      let n = env.args(1)?.u64()?
      let k = env.args(2)?.u64()?
      
      let start_time = Time.nanos()
      let boss = Boss(n, k, env.out, start_time)
      boss.begin_computation()
    else
      env.out.print("Invalid input")
    end

actor Boss
  let _n: U64
  let _k: U64
  let _out: OutStream
  let _start_time: U64
  var _workers: Array[Worker] = Array[Worker]
  var _completed: USize = 0
  var _results: Array[U64] = Array[U64]

  new create(n: U64, k: U64, out: OutStream, start_time: U64) =>
    _n = n
    _k = k
    _out = out
    _start_time = start_time

  be begin_computation() =>
    let worker_count: U64 = 16
    let chunk_size = (_n / worker_count).max(1)

    for i in Range[U64](0, worker_count) do
      let start_val = (i * chunk_size) + 1
      let end_val = if i == (worker_count - 1) then _n else (i + 1) * chunk_size end
      let worker = Worker(this, start_val, end_val, _k)
      _workers.push(worker)
      worker.compute()
    end

  be report_result(value: U64) =>
    _results.push(value)

  be worker_done() =>
    _completed = _completed + 1
    if _completed == _workers.size() then
      let end_time = Time.nanos()
      let elapsed_time = end_time - _start_time
      _out.print("Computation complete")
      _out.print("Results: " + ", ".join(_results.values()))
     end

actor Worker
  let _boss: Boss
  let _start: U64
  let _end: U64
  let _k: U64

  new create(boss: Boss, start_val: U64, end_val: U64, k: U64) =>
    _boss = boss
    _start = start_val
    _end = end_val
    _k = k

  be compute() =>
    for i in Range[U64](_start, _end + 1) do
      var sum: U64 = 0
      for j in Range[U64](0, _k) do
        sum = sum + ((i + j) * (i + j))
      end
      let sqrt = (sum.f64().sqrt()).u64()
      if (sqrt * sqrt) == sum then
        _boss.report_result(i)
      end
    end
    _boss.worker_done()
