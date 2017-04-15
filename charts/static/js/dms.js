let interval = null

function show() {
  let id = document.getElementById('containerID')
  if (id.value === '') {
    id.value = 'all'
  }
  stop()
  interval = showCharts(id.value)
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

function showCharts(id) {
  let charts = [],
    cpus = [],
    mems = [],
    times = []

  return setInterval(function() {
    fetch('http://localhost:8080/stats/' + id).then(response => {
      return response.json()
    }).then(data => {
      if (data.error != undefined) {
        changeServerStatus('500', data.error)
        return
      }

      for (let i = 0; i < data.length; i++) {
        if (!document.getElementById(data[i].Name)) {
          createChartDiv(document.getElementById('chart'), data[i].Name)

          cpus.push(['cpu'])
          mems.push(['mem'])
          times.push(['time'])

          cpus[i] = setData(cpus[i], 'cpu', data[i].CPUPercentage)
          mems[i] = setData(mems[i], 'mem', data[i].MemoryPercentage)
          times[i] = setData(times[i], 'time', new Date())

          charts.push(createChart(data[i].Name, times[i], cpus[i], mems[i]))
        }

        for (let k = 0; k < charts.length; k++) {
          if (charts[k].element.id == data[i].Name) {
            cpus[k] = setData(cpus[k], 'cpu', data[i].CPUPercentage)
            mems[k] = setData(mems[k], 'mem', data[i].MemoryPercentage)
            times[k] = setData(times[k], 'time', new Date())

            charts[k].load({
              columns: [times[k], cpus[k], mems[k]]
            })
          }
        }
      }

      for (let i = 0; i < charts.length; i++) {
        isFound = false

        for (let k = 0; k < data.length; k++) {
          if (charts[i].element.id == data[k].Name) {
            isFound = true
          }
        }

        if (isFound) {
          continue
        }

        let div = document.getElementById(charts[i].element.id)
        let h2 = document.getElementById('h2' + charts[i].element.id)
        if (div && div.parentElement) {
          div.parentElement.removeChild(div)
          h2.parentElement.removeChild(h2)
        }
      }

      changeServerStatus('200')
    }).catch(error => {
      changeServerStatus('500', 'Internal server error')
    })
  }, 1000)
}

function setData(data, dataType, apiData) {
  if (data.length === 25) {
    data.shift()
    data.shift()
    data.unshift(dataType)
  }
  data.push(apiData)

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
