use "time"

actor Main
  let _env: Env
  let _start_time: U64
  var network: (Network tag | None) = None
  let _timers: Timers = Timers
  var _timeout_timer: (Timer tag | None) = None

  new create(env: Env) =>
    _env = env
    _start_time = Time.millis()
    _init()

  be log(msg: String) =>
    _env.out.print(msg)

  be convergence_reached() =>
    let elapsed = (Time.millis() - _start_time)
    _env.out.print("Time to converge: " + elapsed.string() + " milliseconds")
    match network
    | let n: Network tag => 
      n.stop()
      network = None
    end
    match _timeout_timer
    | let t: Timer tag => _timers.cancel(t)
    end
    _env.exitcode(0)
    _env.out.print("Program finished successfully.")

  fun ref _init() =>
    try
      let args = _env.args
      if args.size() != 4 then
        _env.out.print("Usage: project2 <numNodes> <topology> <algorithm>")
        return
      end

      let numNodes = args(1)?.usize()?
      let topology = args(2)?
      let algorithm = args(3)?

      if (algorithm == "gossip") and (numNodes < 2) then
        _env.out.print("Gossip algorithm requires at least 2 nodes.")
        return
      end

      _env.out.print("Starting simulation with " + numNodes.string() + " nodes, topology: " + topology + ", algorithm: " + algorithm)

      let main: Main tag = this
      network = Network(numNodes, topology, algorithm, 
        {(msg: String) => main.log(msg) }, 
        {() => main.convergence_reached() })
      match network
      | let n: Network tag => n.start()
      end

      // Set a timeout
      let timer = Timer(
        object iso is TimerNotify
          let _main: Main tag = main
          
          fun ref apply(timer: Timer, count: U64): Bool =>
            _main.log("Timeout reached. Program may be stuck.")
            _main.convergence_reached()
            false
        end,
        300_000_000_000, // 300 seconds (5 minutes) timeout
        0
      )
      _timeout_timer = timer
      _timers(consume timer)
    else
      _env.out.print("Error parsing arguments")
    end