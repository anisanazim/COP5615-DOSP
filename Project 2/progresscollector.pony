use "collections"

actor ProgressCollector
  let _total: USize
  let _network: Network tag
  let _progress: Array[String]
  var _collected: USize = 0

  new create(total: USize, network: Network tag) =>
    _total = total
    _network = network
    _progress = Array[String].init("", total)

  be collect(index: USize, status: String) =>
    try
      _progress(index)? = "Node " + index.string() + ": " + status
      _collected = _collected + 1
      if _collected == _total then
        let progress_str = _generate_progress_string()
        _network.log_progress(progress_str)
      end
    end

  fun _generate_progress_string(): String val =>
    let s = String(_progress.size() * 20)
    for (i, status) in _progress.pairs() do
      if i > 0 then
        s.append(", ")
      end
      s.append(status)
    end
    s.clone()