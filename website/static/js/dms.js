let interval = null

function show() {
  let ids = document.getElementById('containerID')
  if (ids.value === '') {
    ids.value = 'all'
  }
  stop()
  interval = showCharts(ids.value)
  ids.value = ''
}

function clearCharts() {
  let elements = document.getElementsByClassName("temp")

  for (let i = elements.length - 1; i >= 0; i--) {
    if (elements[i] && elements[i].parentElement) {
      elements[i].parentElement.removeChild(elements[i])
    }
  }
}

function stop() {
  changeServerStatus('404')
  clearInterval(interval)
  clearCharts()
}

function createChartDiv(parent, name) {
  let h2 = document.createElement('h2')
  h2.innerText = name
  h2.setAttribute('class', 'temp')
  h2.setAttribute('id', 'h2' + name)
  parent.appendChild(h2)

  let div = document.createElement('div')
  div.setAttribute('id', name);
  div.setAttribute('class', 'temp');
  parent.appendChild(div)
}

function removeChartDiv(id) {
  let div = document.getElementById(id)
  let h2 = document.getElementById('h2' + id)
  if (div && div.parentElement) {
    div.parentElement.removeChild(div)
    h2.parentElement.removeChild(h2)
  }
}

function showCharts(ids) {
  let chart = new Map(),
    cpu = new Map(),
    mem = new Map(),
    time = new Map()

  return setInterval(function() {
    fetch('http://localhost:8080/stats/' + ids).then(response => {
      return response.json()
    }).then(data => {
      if (data.message) {
        throw data.error
      }

      if (data.stopped) {
        for (let i in data.stopped) {
          if (ids.includes(data.stopped[i]) || ids === 'all') {
            removeChartDiv(data.stopped[i])

            cpu.delete(data.stopped[i])
            mem.delete(data.stopped[i])
            time.delete(data.stopped[i])
            chart.delete(data.stopped[i])
          }
        }
      }

      for (let i in data.data) {
        id = data.data[i].Name

        if (chart.has(id)) {
          cpu.set(id, setData(cpu.get(id), 'cpu', data.data[i].CPUPercentage))
          mem.set(id, setData(mem.get(id), 'mem', data.data[i].MemoryPercentage))
          time.set(id, setData(time.get(id), 'time', new Date()))

          chart.get(id).load({
            columns: [time.get(id), cpu.get(id), mem.get(id)]
          })
        } else {
          createChartDiv(document.getElementById('chart'), id)

          cpu.set(id, setData(['cpu'], 'cpu', 0))
          mem.set(id, setData(['mem'], 'mem', 0))
          time.set(id, setData(['time'], 'time', new Date()))
          chart.set(id, createChart(
            id,
            time.get(id), cpu.get(id),
            mem.get(id)))
        }
      }

      changeServerStatus('200')
    }).catch(error => {
      cpu.clear()
      mem.clear()
      time.clear()
      chart.clear()

      changeServerStatus('500', error)
    })
  }, 3000)
}

function setData(data, type, value) {
  if (data.length === 10) {
    data.shift()
    data.shift()
    data.unshift(type)
  }
  data.push(value)

  return data
}

function changeServerStatus(status, error) {
  let alertStatus = document.getElementById('alert-status')
  let alertErrorText = document.getElementById('alert-error-text')

  alertErrorText.setAttribute('class', 'none')

  if (status === '500') {
    clearCharts()

    alertStatus.innerText = '500'
    alertStatus.setAttribute('class', 'alert alert-danger')

    alertErrorText.innerText = error
    alertErrorText.setAttribute('class', 'alert alert-danger')

    return
  } else if (status === '404') {
    clearCharts()

    alertStatus.innerText = '404'
    alertStatus.setAttribute('class', 'alert alert-warning')

    return
  }

  alertStatus.innerText = '200'
  alertStatus.setAttribute('class', 'alert alert-success')
}

function createChart(id, time, cpu, mem) {
  return chart = c3.generate({
    bindto: '#' + id,
    data: {
      x: 'time',
      columns: [
        time,
        cpu,
        mem
      ]
    },
    axis: {
      x: {
        type: 'timeseries',
        tick: {
          format: '%H:%M:%S'
        }
      },
      y: {
        tick: {
          format: d3.format(',.2f')
        }
      }
    },
    grid: {
      x: {
        show: true
      },
      y: {
        show: true
      }
    }
  })
}
