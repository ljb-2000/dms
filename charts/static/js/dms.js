let interval = null
let intervalTimeOut = 1000

function showOneChart() {
  clearInterval(interval)
  clearCharts()
  interval = oneContainer()
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
  let isFirst = true
  let charts = []
  let cpus = []
  let mems = []

  return setInterval(function() {
    fetch('http://localhost:8080/get/all').then(response => {
      return response.json()
    }).then(data => {
      if (isFirst) {
        for (let i = 0; i < data.length; i++) {
          createChartDiv(document.getElementById('wrapper'), i, data[i].Name)

          charts.push(createChart('chart' + i, data[i].CPUPercentage, data[i].MemoryPercentage))
          cpus.push(['CPU'])
          mems.push(['MEM'])
        }

        isFirst = false
      }

      for (let i = 0; i < data.length; i++) {
        if (cpus[i].length === 11 && mems[i].length === 11) {
          cpus[i].shift()
          cpus[i].shift()
          mems[i].shift()
          mems[i].shift()

          cpus[i].unshift('CPU')
          mems[i].unshift('MEM')
        }
        cpus[i].push(data[i].CPUPercentage)
        mems[i].push(data[i].MemoryPercentage)
      }

      for (let i = 0; i < data.length; i++) {
        charts[i].load({
          columns: [cpus[i], mems[i]]
        })
      }
    })
  }, intervalTimeOut)
}

function oneContainer() {
  let isFirstChart = true
  let cpu = ['CPU']
  let mem = ['MEM']
  let id = document.getElementById('containerID').value

  return setInterval(function() {
    fetch('http://localhost:8080/get/' + id).then(response => {
      return response.json()
    }).then(data => {
      if (cpu.length === 11 && mem.length === 11) {
        cpu.shift()
        cpu.shift()
        mem.shift()
        mem.shift()

        cpu.unshift('CPU')
        mem.unshift('MEM')
      }
      cpu.push(data.CPUPercentage)
      mem.push(data.MemoryPercentage)

      if (isFirstChart) {
        createChartDiv(document.getElementById('wrapper'), 0, data.Name)

        chart = createChart('chart0', cpu, mem)

        isFirstChart = false
        return
      }

      chart.load({
        columns: [cpu, mem]
      })
    })
  }, intervalTimeOut)
}

function createChart(id, cpu, mem) {
  return chart = c3.generate({
    bindto: '#' + id,
    data: {
      columns: [
        cpu,
        mem
      ]
    },
    axis: {
      x: {
        max: 9
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
