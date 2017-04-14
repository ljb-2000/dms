let interval = null,
  intervalTimeOut = 3000,
  apiHost = 'http://localhost:8080/stats/'

function showOneChart() {
  let id = document.getElementById('containerID')

  stop()
  interval = oneContainer(id.value)

  id.value = ''
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

function allContainers() {
  let isFirst = true,
    charts = [],
    cpus = [],
    mems = [],
    times = []

  return setInterval(function() {
    fetch(apiHost + 'all').then(response => {
      return response.json()
    }).then(data => {
      console.log(data)
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

function oneContainer(id) {
  let isFirstChart = true,
    cpu = ['cpu'],
    mem = ['mem'],
    time = ['time']

  return setInterval(function() {
    fetch(apiHost + id).then(response => {
      return response.json()
    }).then(data => {
      if (data.error != undefined) {
        changeServerStatus('500', data.message)
        return
      }

      cpu = setCPU(cpu, data)
      mem = setMEM(mem, data)
      time = setTime(time)

      if (isFirstChart) {
        createChartDiv(document.getElementById('chart'), data.Name)

        chart = createChart(data.Name, time, cpu, mem)

        isFirstChart = false
        changeServerStatus('200')
        return
      }

      chart.load({
        columns: [time, cpu, mem]
      })

      changeServerStatus('200')
    })
  }, intervalTimeOut)
}

function changeServerStatus(status, error) {
  let alertStatus = document.getElementById('alert-status')
  let alertErrorText = document.getElementById('alert-error-text')

  if (status === '500') {
    console.log('error: ', error)
    alertStatus.innerText = '500'
    alertStatus.setAttribute('class', 'alert alert-danger')

    alertErrorText.innerText = error
    alertErrorText.setAttribute('class', 'alert alert-danger')
    return
  }
  alertStatus.innerText = '200'
  alertStatus.setAttribute('class', 'alert alert-success')

  alertErrorText.innerText = ''
  alertErrorText.setAttribute('class', 'none')
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
