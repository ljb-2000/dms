
let id = 0, chart = null, first = true

function newContainer() {
  id = document.getElementById("containerID").value
}

function createChart(cpu, mem) {
  chart = c3.generate({
    bindto: '#chart',
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

function work() {
  let cpu = ['CPU']
  let mem = ['MEM']
  let name = document.getElementById("containerName")

  setInterval(function () {
    if (id === 0) {
      return
    }

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

      name.innerText = data.Name

      if (first) {
        createChart(cpu, mem)

        first = false
        return
      }

      chart.load({
        columns: [
          cpu,
          mem
        ]
      })
    })
  }, 3000)
}

work()
