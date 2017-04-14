let interval = null

function show(isAllContainers) {
  if (isAllContainers) {
    stop()
    interval = showCharts('all')
  } else {
    let id = document.getElementById('containerID')
    stop()
    interval = showCharts(id.value)
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

function showCharts(id) {
  let isFirst = true,
    charts = [],
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

      if (isFirst) {
        for (let i = 0; i < data.length; i++) {
          createChartDiv(document.getElementById('chart'), data[i].Name)

          cpus.push(['cpu'])
          mems.push(['mem'])
          times.push(['time'])

          cpus[i] = setData(cpus[i], 'cpu', data[i].CPUPercentage)
          mems[i] = setData(mems[i], 'mem', data[i].MemoryPercentage)
          times[i] = setData(times[i], 'time', new Date())

          charts.push(createChart(data[i].Name, times[i], cpus[i], mems[i]))
        }

        isFirst = false
        changeServerStatus('200')
        return
      }

      for (let i = 0; i < data.length; i++) {
        cpus[i] = setData(cpus[i], 'cpu', data[i].CPUPercentage)
        mems[i] = setData(mems[i], 'mem', data[i].MemoryPercentage)
        times[i] = setData(times[i], 'time', new Date())

        charts[i].load({
          columns: [times[i], cpus[i], mems[i]]
        })
      }

      changeServerStatus('200')
    }).catch(error => {
      changeServerStatus('500', error, false)
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
