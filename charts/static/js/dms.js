let interval = null,
  intervalTimeOut = 3000,
  apiHost = 'http://localhost:8080/stats/'

function showChart(type) {
  if (type) {
    stop()
    interval = oneOrMoreContainers('all')
  } else {
    let id = document.getElementById('containerID')

    stop()
    interval = oneOrMoreContainers(id.value)

    id.value = ''
  }
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

function showAllCharts() {
  stop()
  interval = allContainers()
}

function createChartDiv(parent, name) {
  let h2 = document.createElement('h2')
  h2.innerText = name
  h2.setAttribute('class', 'temp')
  parent.appendChild(h2)

  let div = document.createElement('div')
  div.setAttribute('id', name);
  div.setAttribute('class', 'temp');
  parent.appendChild(div)
}

function oneOrMoreContainers(id) {
  let isFirst = true,
    charts = [],
    cpus = [],
    mems = [],
    times = []

  return setInterval(function() {
    fetch(apiHost + id).then(response => {
      return response.json()
    }).then(data => {
      if (data.error != undefined) {
        changeServerStatus('500', data.error)
        return
      }

      if (isFirst) {
        for (let i = 0; i < data.length; i++) {
          createChartDiv(document.getElementById('chart'), data[i].Name)

          cpus.push(['cpu'])
          mems.push(['mem'])
          times.push(['time'])

          cpus[i] = setCPU(cpus[i], data[i])
          mems[i] = setMEM(mems[i], data[i])
          times[i] = setTime(times[i])

          charts.push(createChart(data[i].Name, times[i], cpus[i], mems[i]))
        }

        isFirst = false

        changeServerStatus('200')
        return
      }

      for (let i = 0; i < data.length; i++) {
        cpus[i] = setCPU(cpus[i], data[i])
        mems[i] = setMEM(mems[i], data[i])
        times[i] = setTime(times[i])

        charts[i].load({
          columns: [times[i], cpus[i], mems[i]]
        })
      }

      changeServerStatus('200')
    }).catch(error => {
      changeServerStatus('500', error, false)
    })
  }, intervalTimeOut)
}

function setCPU(cpu, data) {
  if (cpu.length === 11) {
    cpu.shift()
    cpu.shift()

    cpu.unshift('cpu')
  }
  cpu.push(data.CPUPercentage)

  return cpu
}

function setMEM(mem, data) {
  if (mem.length === 11) {
    mem.shift()
    mem.shift()

    mem.unshift('mem')
  }
  mem.push(data.MemoryPercentage)

  return mem
}

function setTime(time) {
  if (time.length === 11) {
    time.shift()
    time.shift()

    time.unshift('time')
  }
  time.push(new Date())

  return time
}

function changeServerStatus(status, error, isAlert) {
  let alertStatus = document.getElementById('alert-status')
  let alertErrorText = document.getElementById('alert-error-text')

  if (status === '500') {
    console.log('error: ', error)
    alertStatus.innerText = '500'
    alertStatus.setAttribute('class', 'alert alert-danger')

    if (!isAlert) {
      return
    }

    alertErrorText.innerText = error
    alertErrorText.setAttribute('class', 'alert alert-danger temp')
    return
  } else if (status === '404') {
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
        },
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
