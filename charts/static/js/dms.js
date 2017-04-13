let interval = null
let intervalTimeOut = 3000

function showOneChart() {
  let id = document.getElementById('containerID')

  clearInterval(interval)
  clearCharts()
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
  clearInterval(interval)
  clearCharts()
  interval = allContainers()
}

function createChartDiv(parent, i, name) {
  let h2 = document.createElement('h2')
  h2.innerText = name
  h2.setAttribute('class', 'temp')
  parent.appendChild(h2)

  let div = document.createElement('div')
  div.setAttribute('id', 'chart' + i);
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
    fetch('http://localhost:8080/get/all').then(response => {
      return response.json()
    }).then(data => {
      if (isFirst) {
        for (let i = 0; i < data.length; i++) {
          createChartDiv(document.getElementById('wrapper'), i, data[i].Name)

          cpus.push(['cpu'])
          mems.push(['mem'])
          times.push(['time'])

          cpus[i] = setCPU(cpus[i], data[i])
          mems[i] = setMEM(mems[i], data[i])
          times[i] = setTime(times[i])

          charts.push(createChart('chart' + i, times[i], cpus[i], mems[i]))
        }

        isFirst = false

        console.log(cpus)
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
    fetch('http://localhost:8080/get/' + id).then(response => {
      return response.json()
    }).then(data => {
      cpu = setCPU(cpu, data)
      mem = setMEM(mem, data)
      time = setTime(time)

      if (isFirstChart) {
        createChartDiv(document.getElementById('wrapper'), 0, data.Name)

        chart = createChart('chart0', time, cpu, mem)

        isFirstChart = false
        return
      }

      chart.load({
        columns: [time, cpu, mem]
      })
    })
  }, intervalTimeOut)
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
