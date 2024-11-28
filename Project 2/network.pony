use "collections"
use "time"
use "random"

actor Network
  let _nodes: Array[Node]
  let _topology: String
  let _algorithm: String
  let _log: {(String)} val
  let _convergence_reached: {()} val
  var _converged_count: USize = 0
  let _total_nodes: USize
  let _converged_nodes: Set[USize] = Set[USize]
  var _is_stopping: Bool = false
  let _timers: Timers = Timers
  var _check_timer: (Timer tag | None) = None

  new create(numNodes: USize, topology: String, algorithm: String, log: {(String)} val, convergence_reached: {()} val) =>
    _nodes = Array[Node](numNodes)
    _topology = topology
    _algorithm = algorithm
    _log = log
    _convergence_reached = convergence_reached
    _total_nodes = numNodes
    
    for i in Range(0, numNodes) do
      _nodes.push(Node.create(this, i, algorithm, log))
    end
       
    _setup_topology()

  fun ref _setup_topology() =>
    match _topology
    | "full" => _setup_full()
    | "3D" => _setup_3d()
    | "line" => _setup_line()
    | "imp3D" => _setup_imp3D()
    else
      _log("Unknown topology: " + _topology)
    end

  fun ref _setup_full() =>
    for i in Range(0, _nodes.size()) do
      for j in Range(0, _nodes.size()) do
        if i != j then
          try
            _nodes(i)?.add_neighbor(_nodes(j)?)
          end
        end
      end
    end

  fun ref _setup_3d() =>
  let num_nodes = _nodes.size()
  let lx: USize = _approximate_cuberoot(num_nodes)  
  let ly: USize = lx                                
  let lz: USize = (num_nodes / (lx * ly)).usize()   

  for i in Range(0, _nodes.size()) do
    try
      let x = i % lx
      let y = (i / lx) % ly
      let z = i / (lx * ly)

      if x > 0 then
        _nodes(i)?.add_neighbor(_nodes(i - 1)?)
      end

      if x < (lx - 1) then
        _nodes(i)?.add_neighbor(_nodes(i + 1)?)
      end

      if y > 0 then
        _nodes(i)?.add_neighbor(_nodes(i - lx)?)
      end

      if y < (ly - 1) then
        _nodes(i)?.add_neighbor(_nodes(i + lx)?)
      end

      if z > 0 then
        _nodes(i)?.add_neighbor(_nodes(i - (lx * ly))?)
      end

      if z < (lz - 1) then
        _nodes(i)?.add_neighbor(_nodes(i + (lx * ly))?)
      end
    end
  end

  fun _approximate_cuberoot(num: USize): USize =>
    var x: USize = num
    var y: USize = 0
    while x > y do
      x = (x + y) / 2
      y = num / (x * x)
    end
    x

  fun ref _setup_line() =>
    for i in Range(0, _nodes.size()) do
      try
        if i > 0 then
          _nodes(i)?.add_neighbor(_nodes(i-1)?)
        end
        if i < (_nodes.size() - 1) then
          _nodes(i)?.add_neighbor(_nodes(i+1)?)
        end
      end
    end

   fun ref _setup_imp3D() =>
   let num_nodes = _nodes.size()
   let lx: USize = _approximate_cuberoot(num_nodes)  
   let ly: USize = lx                                
   let lz: USize = (num_nodes / (lx * ly)).usize()   
   
   for i in Range(0, _nodes.size()) do
     try
       let x = i % lx
       let y = (i / lx) % ly
       let z = i / (lx * ly)
   
       if x > 0 then
         _nodes(i)?.add_neighbor(_nodes(i - 1)?)
       end
   
       if x < (lx - 1) then
         _nodes(i)?.add_neighbor(_nodes(i + 1)?)
       end
   
       if y > 0 then
         _nodes(i)?.add_neighbor(_nodes(i - lx)?)
       end
   
       if y < (ly - 1) then
         _nodes(i)?.add_neighbor(_nodes(i + lx)?)
       end
   
       if z > 0 then
         _nodes(i)?.add_neighbor(_nodes(i - (lx * ly))?)
       end
   
       if z < (lz - 1) then
         _nodes(i)?.add_neighbor(_nodes(i + (lx * ly))?)
       end
     end
   end
   
   let random = Rand
   for (i, node) in _nodes.pairs() do
     try
       let rand_index = random.int(_nodes.size().u64()).usize()
       if rand_index != i then
         node.add_neighbor(_nodes(rand_index)?)
       end
     end
   end

  be node_converged(id: USize) =>
    if not _converged_nodes.contains(id) then
      _converged_nodes.set(id)
      _converged_count = _converged_count + 1
      _log("Node " + id.string() + " has converged. Total: " + _converged_count.string() + "/" + _total_nodes.string())
      if _converged_count == _total_nodes then
        _convergence_reached()
      end
    end

  be start() =>
    match _algorithm
    | "gossip" =>
      try
        _nodes(0)?.start("Initial rumor")
      end
    | "push-sum" =>
      for node in _nodes.values() do
        node.start("")
      end
    end
    _start_check_timer()

  fun ref _start_check_timer() =>
    let timer = Timer(
      object iso is TimerNotify
        let _network: Network = this
        fun ref apply(timer: Timer, count: U64): Bool =>
          _network.check_progress()
          true
      end,
      1_000_000_000, // 1 second
      1_000_000_000  // 1 second between checks
    )
    _check_timer = timer
    _timers(consume timer)

  be check_progress() =>
    if _converged_count < _total_nodes then
      try
        for i in Range(0, _total_nodes) do
          if not _converged_nodes.contains(i) then
            _nodes(i)?.start("Reminder")
          end
        end
      end
    else
      match _check_timer
      | let t: Timer tag => _timers.cancel(t)
      end
    end

  be log_progress(progress: String) =>
    _log(progress)

  be stop() =>
    if not _is_stopping then
      _is_stopping = true
      match _check_timer
      | let t: Timer tag => _timers.cancel(t)
      end
      for node in _nodes.values() do
        node.stop()
      end
    end