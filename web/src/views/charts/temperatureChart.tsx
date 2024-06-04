import { ApexOptions } from 'apexcharts'

import ReactApexcharts from 'src/@core/components/react-apexcharts'

import Card from '@mui/material/Card'
import { useTheme } from '@mui/material/styles'
import CardHeader from '@mui/material/CardHeader'
import CardContent from '@mui/material/CardContent'
import Grid from '@mui/material/Grid'

const lineColors = [
  '#FF5733',
  '#33FF57',
  '#5733FF',
  '#FF33FF',
  '#33FFFF',
  '#FFC133',
  '#8E44AD',
  '#3498DB',
  '#E74C3C',
  '#2ECC71',
  '#F39C12',
  '#D35400'
]

// TODO: Modify this after fixing the timeseries batch operation
function convertToClockFormat(timestamps: string[]): string[] {
  let lastTimeShown: Date | null = null // Track the last timestamp displayed

  return timestamps.map(ts => {
    const date = new Date(ts)
    const hours = String(date.getUTCHours()).padStart(2, '0')

    if (date.getUTCMinutes() >= 30) {
      return '' // Skip half-past times
    }

    if (lastTimeShown) {
      const diffMinutes = (date.getTime() - lastTimeShown.getTime()) / (1000 * 60)
      if (diffMinutes < 60) {
        return '' // Less than 60 minutes since last label, so skip
      }
    }

    lastTimeShown = date // Update the last shown time

    return `${hours}:00`
  })
}

function getYaxisBounds(data: number[], variancePercent = 30) {
  const maxVal = Math?.max(...data) || 0

  const variance = maxVal * (variancePercent / 100)

  return {
    max: maxVal + variance
  }
}

// ===== logic to create mock data set =======
// this includes the logic to convert the data to the format that the chart expects

type TemperatureChartProps = {
  tempSensorArr: number[][]
  timeStampArr: string[]
}

// ==========================================

const TemperatureChart = ({ tempSensorArr, timeStampArr }: TemperatureChartProps) => {
  const theme = useTheme()
  let seriesData = []
  let yaxisBounds = { max: 0 }
  if (tempSensorArr.length > 0) {
    // FIXME
    // TODO! need a logic to traverse the number of sensors
    // Determine the maximum number of sensors across all records
    const maxSensors = tempSensorArr?.reduce((max, record) => Math.max(max, record?.length), 0) || 0

    seriesData = []
    for (let i = 0; i < maxSensors; i++) {
      const sensorData = tempSensorArr.map(record => (i < record.length ? record[i] : null))
      seriesData.push({
        name: `Sensor ${i + 1}`,
        data: sensorData
      })
    }
  }

  yaxisBounds =
    tempSensorArr.length > 0
      ? getYaxisBounds(tempSensorArr.flat())
      : {
          max: 0
        } // Assuming getYaxisBounds expects a flat array of numbers

  const options: ApexOptions = {
    chart: {
      parentHeightOffset: 0,
      zoom: { enabled: true },
      toolbar: {
        show: true,
        offsetY: -15 // Introducing a gap of 15 pixels. Adjust this value based on your requirements.
      },
      animations: {
        enabled: false // This will disable the initial animation
      }
    },
    stroke: {
      width: [1], // Set width for each series. You can adjust the value as needed.
      curve: 'smooth' // This is optional and just demonstrates that you can set other stroke properties as well.
    },
    legend: {
      labels: {
        colors: [theme.palette.text.primary] // Set the color for legend labels using the theme object
      }
    },
    colors: lineColors,
    dataLabels: { enabled: false },
    markers: {
      strokeWidth: 0.3,
      strokeOpacity: 1,
      colors: [theme.palette.primary.main], // Change the dot's fill color to the theme's primary main color
      strokeColors: [theme.palette.divider],
      hover: {
        size: 5, // Optional: You can adjust the size of the dot on hover if needed.
        sizeOffset: 3 // Optional: This value changes the radius of the circle when hovering.
      }
    },
    grid: {
      padding: { top: -10 },
      borderColor: theme.palette.divider,
      xaxis: {
        lines: { show: true }
      }
    },
    tooltip: {
      shared: true,
      custom(data: any) {
        let tooltipContent = `
            <div class='bar-chart' style='background-color: ${theme.palette.background.paper}; padding: 8px; border-radius: 4px; color: ${theme.palette.text.primary}; font-size: 0.85em;'>
          `

        // Get the corresponding timestamp for the current data point
        const date = new Date(timeStampArr[data.dataPointIndex])
        const monthNames = [
          'January',
          'February',
          'March',
          'April',
          'May',
          'June',
          'July',
          'August',
          'September',
          'October',
          'November',
          'December'
        ]
        const formattedDate = `${String(date.getUTCHours()).padStart(2, '0')}:${String(date.getUTCMinutes()).padStart(
          2,
          '0'
        )}, ${monthNames[date.getUTCMonth()]} ${date.getUTCDate()}`

        // tooltipContent += `<strong>${formattedDate}</strong><br/>`

        data.series.forEach((series: number[], seriesIndex: number) => {
          const seriesColor = lineColors[seriesIndex % lineColors.length]
          tooltipContent += `<span style='color: ${seriesColor}'>▶</span> ${series[data.dataPointIndex]} (${
            data.w.globals.seriesNames[seriesIndex]
          })<br/>`
        })

        tooltipContent += `</div>`

        return tooltipContent
      }
    },
    yaxis: {
      labels: {
        style: { colors: theme.palette.text.primary } // Set the y-axis label color using the theme object
      },
      min: 0, // Updated min value
      max: yaxisBounds.max // Updated max value
    },
    xaxis: {
      axisBorder: { show: false },
      axisTicks: { color: theme.palette.divider },
      crosshairs: {
        stroke: { color: theme.palette.divider }
      },
      labels: {
        style: { colors: theme.palette.text.primary },
        rotate: -45, // Rotates labels by -10 degrees
        rotateAlways: true // Ensures labels are always rotated as specified
      },
      categories: convertToClockFormat(timeStampArr),
      // TODO: UNDO this after fixing the timeseries batch logic
      tooltip: {
        enabled: false // Disables x-axis tooltip
      }
    }
  }

  return (
    <Grid container spacing={6}>
      <Grid item xs={12}>
        <Card>
          <CardHeader title='Hashboard Temperature' />
          <CardContent>
            <ReactApexcharts type='line' height={400} options={options} series={seriesData} />
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  )
}

export default TemperatureChart
