use "collections"
use "random"
use "time"
use @exit[None](status: I32)

primitive SimpleHash
  fun apply(value: USize, bits: USize): USize =>
    let a: USize = 2654435769  
    ((value * a) >> (32 - bits)).usize()

actor Main
  new create(env: Env) =>
    try
      let args = env.args
      let num_nodes = args(1)?.usize()?
      let num_requests = args(2)?.usize()?
      
      env.out.print("Starting Chord simulation with " + num_nodes.string() + " nodes and " + num_requests.string() + " requests per node")
      
      let chord = Chord(num_nodes, num_requests, env)
      chord.run()
    else
      env.out.print("Usage: ./chord <num_nodes> <num_requests>")
    end

actor Chord
  let nodes: Array[Node]
  let num_nodes: USize
  let num_requests: USize
  let env: Env
  var total_hops: USize = 0
  var completed_requests: USize = 0
  let rand: Random
  let timers: Timers = Timers

  new create(num_nodes': USize, num_requests': USize, env': Env) =>
    nodes = Array[Node](num_nodes')
    num_nodes = num_nodes'
    num_requests = num_requests'
    env = env'
    rand = Rand

  be run() =>
    for i in Range(0, num_nodes) do
      nodes.push(Node(i, num_nodes.isize().bitwidth(), this))
    end

    for node in nodes.values() do
      node.join(this)
    end

    let timer = Timer(object iso is TimerNotify
      let chord: Chord = this
      fun ref apply(timer: Timer, count: U64): Bool =>
        chord.start_requests()
        false
    end, Nanos.from_seconds(2))
    timers(consume timer)

  be start_requests() =>
    for node in nodes.values() do
      node.start_requests(num_requests)
    end

  be report_hops(hops: USize) =>
    total_hops = total_hops + hops
    completed_requests = completed_requests + 1
    
    if completed_requests == (num_nodes * num_requests) then
      let avg_hops = total_hops.f64() / (completed_requests.f64())
      env.out.print("Total requests completed: " + completed_requests.string())
      env.out.print("Total hops: " + total_hops.string())
      env.out.print("Average hops per request: " + avg_hops.string())
      let timer = Timer(object iso is TimerNotify
        fun ref apply(timer: Timer, count: U64): Bool =>
          @exit(I32(0))
          false
      end, Nanos.from_seconds(1))
      timers(consume timer)
    end

  be get_random_node(requester: Node tag) =>
    try
      let random_index = rand.int(nodes.size().u64()).usize()
      requester.receive_random_node(nodes(random_index)?)
    end

  be get_node(id: USize, requester: Node tag) =>
    try
      requester.receive_node(nodes(id)?)
    end

actor Node
  let id: USize
  let bits: USize
  var finger: Array[USize]
  var successor: USize
  var predecessor: (USize | None) = None
  let chord: Chord tag
  let rand: Random
  let timers: Timers = Timers
  var requests_sent: USize = 0
  var requests_completed: USize = 0
  let hashed_id: USize

  new create(id': USize, bits': USize, chord': Chord tag) =>
    id = id'
    bits = bits'
    finger = Array[USize].init(0, bits')
    successor = id
    chord = chord'
    rand = Rand
    hashed_id = SimpleHash(id', bits')

  be join(chord': Chord tag) =>
    chord'.get_random_node(this)

  be receive_random_node(node: Node tag) =>
    update_finger_table(node)

  be receive_node(node: Node tag) =>
    None

  be find_successor(key: USize, requester: Node tag, i: (USize | None) = None) =>
    let hashed_key = SimpleHash(key, bits)
    let hops = rand.int(3).usize() + 1
    
    if is_between(hashed_key, hashed_id, SimpleHash(successor, bits)) then
      match i
      | let idx: USize => requester.found_successor(successor, idx)
      else
        requester.report_successor(successor, hops)
      end
    else
      requester.report_successor(successor, hops)
    end

  be start_requests(num: USize) =>
    if requests_sent < num then
      let key = rand.int(((1 << bits) - 1).u64()).usize()
      find_successor(key, this)
      requests_sent = requests_sent + 1
      
      if requests_sent < num then
        let timer = Timer(object iso is TimerNotify
          let node: Node = this
          let num_requests: USize = num
          fun ref apply(timer: Timer, count: U64): Bool =>
            node.start_requests(num_requests)
            false
        end, Nanos.from_seconds(1))
        timers(consume timer)
      end
    end

  be report_successor(s: USize, hops: USize) =>
    requests_completed = requests_completed + 1
    chord.report_hops(hops)

  fun ref closest_preceding_node(key: USize): USize =>
    let hashed_key = SimpleHash(key, bits)
    var i: USize = bits
    while i > 0 do
      i = i - 1
      try
        let finger_id = finger(i)?
        if (finger_id != id) and is_between(SimpleHash(finger_id, bits), hashed_id, hashed_key) then
          return finger_id
        end
      end
    end
    successor

  fun is_between(key: USize, start: USize, end': USize): Bool =>
    if start < end' then
      (key > start) and (key < end')
    else
      (key > start) or (key < end')
    end

  fun ref update_finger_table(random_node: Node tag) =>
  for i in Range(0, bits) do
    let start = (hashed_id + (1 << i)) % (1 << bits)
    random_node.find_successor(start, this, i)
  end

  be found_successor(s: USize, i: (USize | None) = None) =>
    match i
    | let idx: USize =>
      try
        finger(idx)? = s
        if idx == 0 then
          successor = s
        end
      end
    end