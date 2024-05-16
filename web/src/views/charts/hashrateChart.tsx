import { ChangeEvent, useState } from 'react'

import Link from 'next/link'

import CustomChip from 'src/@core/components/mui/chip'
import { ApexOptions } from 'apexcharts'
import ReactApexcharts from 'src/@core/components/react-apexcharts'

import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import { useTheme } from '@mui/material/styles'
import CardHeader from '@mui/material/CardHeader'
import Typography from '@mui/material/Typography'
import CardContent from '@mui/material/CardContent'
import Grid from '@mui/material/Grid'

type ChartDataType = {
  series: { data: number[] }[]
  timestamps: string[]
}

type HashrateChartProps = {
  hashrateArr: number[]
  timestampArr: string[]
}

// Util functions
function convertToDecimalFormat(data: ChartDataType): ChartDataType {
  const transformedSeries = data.series[0].data.map(value => Math.round((value / 1000) * 100) / 100)
  return {
    ...data,
    series: [{ data: transformedSeries }]
  }
}

function convertToClockFormat(timestamps: string[]): string[] {
  return timestamps.map(ts => {
    const date = new Date(ts)
    const hours = String(date.getUTCHours()).padStart(2, '0')
    const minutes = String(date.getUTCMinutes()).padStart(2, '0')

    // Construct time label based on hours and minutes
    // For better readability and precision on the graph, format to show minutes when not zero
    let timeLabel = `${hours}:00` // Default to full hour
    if (minutes !== '00') {
      timeLabel = `${hours}:${minutes}` // Include minutes when they are not '00'
    }

    console.log('time label', timeLabel)

    return timeLabel
  })
}

function getYaxisBounds(data: number[], variancePercent: number = 10) {
  const maxVal = Math.max(...data)
  const variance = maxVal * (variancePercent / 100)

  return {
    max: maxVal + variance
  }
}

const HashrateChart = ({ hashrateArr, timestampArr }: HashrateChartProps) => {
  console.log('hashrate chart timestamp', timestampArr)

  const chartData: ChartDataType = {
    series: [
      {
        data: [...hashrateArr]
      }
    ],
    timestamps: [...timestampArr]
  }

  console.log('hash arr in hashrate chart', hashrateArr)

  const transformedData = convertToDecimalFormat(chartData)
  const yaxisBounds = getYaxisBounds(transformedData.series[0].data)

  const theme = useTheme()

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
    colors: [theme.palette.primary.main],
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
      custom(data: any) {
        const date = new Date(timestampArr[data.dataPointIndex])
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

        return `
            <div class='bar-chart' style='background-color: ${
              theme.palette.background.paper
            }; padding: 8px; border-radius: 4px; color: ${theme.palette.text.primary};'>
              <span>${data.series[data.seriesIndex][data.dataPointIndex]} Th/s</span>
            </div>`
      }
    },
    yaxis: {
      labels: {
        style: { colors: theme.palette.text.disabled }
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
        style: { colors: theme.palette.text.disabled },
        rotate: -45, // Rotates labels by -10 degrees
        rotateAlways: true // Ensures labels are always rotated as specified
      },
      categories: convertToClockFormat(chartData?.timestamps),
      // TODO: UNDO this after fixing the timeseries batch logic
      tooltip: {
        enabled: false // Disables x-axis tooltip
      }
    }
  }

  console.log('time stamp array', chartData?.timestamps)

  return (
    <Grid container spacing={6}>
      <Grid item xs={12}>
        <Card>
          <CardHeader title='Hashrate' />
          <CardContent>
            <ReactApexcharts
              type='line'
              height={400}
              options={options}
              series={convertToDecimalFormat(chartData)?.series}
            />
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  )
}

export default HashrateChart
