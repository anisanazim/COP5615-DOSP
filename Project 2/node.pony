use "time"
use "random"

actor Node
  let _id: USize
  let _neighbors: Array[Node] = Array[Node]
  let _rand: Rand
  let _network: Network
  let _log: {(String)} val
  let _algorithm: String
  var _converged: Bool = false
  let _timers: Timers = Timers

  // Gossip-specific fields
  var _rumor_count: USize = 0

  // Push-sum specific fields
  var _s: F64
  var _w: F64 = 1.0
  var _old_ratio: F64
  var _rounds_stable: USize = 0
  let _epsilon: F64 = 10e-10
  let _max_stable_rounds: USize = 3
  var _active: Bool = true

  new create(network: Network, id: USize, algorithm: String, log: {(String)} val) =>
    _network = network
    _id = id
    _algorithm = algorithm
    _rand = Rand(Time.nanos())
    _s = id.f64()
    _old_ratio = _s / _w
    _log = log

  be add_neighbor(neighbor: Node) =>
    _neighbors.push(neighbor)

  be start(initial_message: String) =>
    match _algorithm
    | "gossip" => receive_rumor(initial_message)
    | "push-sum" => 
      _send_push_sum()
      _start_push_sum_timer()
    end

  fun ref _start_push_sum_timer() =>
    let timer = Timer(PushSumNotify(this), 0, 10_000_000) // Check every 10ms
    _timers(consume timer)

  be receive_rumor(rumor: String) =>
    if (_algorithm == "gossip") and (not _converged) then
      _rumor_count = _rumor_count + 1
      _network.log_progress("Node " + _id.string() + " rumor count: " + _rumor_count.string())
      if _rumor_count >= 10 then
        _converge()
      else
        _send_rumor(rumor)
      end
    end

  fun ref _send_rumor(rumor: String) =>
    try
      let next = _neighbors(_rand.int(_neighbors.size().u64()).usize())?
      next.receive_rumor(rumor)
    end

  be receive_push_sum(s': F64, w': F64) =>
    if (_algorithm == "push-sum") and _active then
      _s = _s + s'
      _w = _w + w'
      _send_push_sum()
    end

  fun ref _send_push_sum() =>
    if _active and (_neighbors.size() > 0) then
      let s_to_send = _s / 2
      let w_to_send = _w / 2
      _s = _s / 2
      _w = _w / 2
      try
        let next = _neighbors(_rand.int(_neighbors.size().u64()).usize())?
        next.receive_push_sum(s_to_send, w_to_send)
      else
        _log("Node " + _id.string() + " has no neighbors to send to")
      end
    end

  be check_push_sum_termination() =>
    if _active then
      let current_ratio = _s / _w
      if (current_ratio - _old_ratio).abs() < _epsilon then
        _rounds_stable = _rounds_stable + 1
        if _rounds_stable >= _max_stable_rounds then
          _converge()
        end
      else
        _rounds_stable = 0
      end
      _old_ratio = current_ratio
      if not _converged then
        _send_push_sum()
      end
    end

  fun ref _converge() =>
    if not _converged then
      _converged = true
      _active = false
      _timers.dispose()
      _network.node_converged(_id)
      match _algorithm
      | "gossip" => _log("Node " + _id.string() + " has converged with rumor count: " + _rumor_count.string())
      | "push-sum" => _log("Node " + _id.string() + " has converged with ratio: " + (_s / _w).string())
      end
    end

  be stop() =>
    _converged = true
    _active = false
    _timers.dispose()
    //match _algorithm
    //| "gossip" => _log("Node " + _id.string() + " stopped with rumor count: " + _rumor_count.string())
    //| "push-sum" => _log("Node " + _id.string() + " stopped with ratio: " + (_s / _w).string())
    //end

class PushSumNotify is TimerNotify
  let _node: Node tag

  new iso create(node: Node tag) =>
    _node = node

  fun ref apply(timer: Timer, count: U64): Bool =>
    _node.check_push_sum_termination()
    true

  fun ref cancel(timer: Timer) =>
    None