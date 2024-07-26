import { forwardRef, useState } from 'react'

import Card from '@mui/material/Card'
import TextField from '@mui/material/TextField'
import { useTheme } from '@mui/material/styles'
import CardHeader from '@mui/material/CardHeader'
import CardContent from '@mui/material/CardContent'
import InputAdornment from '@mui/material/InputAdornment'

import format from 'date-fns/format'
import { ApexOptions } from 'apexcharts'

import Icon from 'src/@core/components/icon'

import { DateType } from 'src/types/forms/reactDatepickerTypes'

import ReactApexcharts from 'src/@core/components/react-apexcharts'
import { PoolStatsData, PoolTimeSeriesDataResponse } from 'src/types/apps/minerDetailsTypes'

interface PickerProps {
  start: Date | number
  end: Date | number
}

type RecetPoolStatsChartProps = {
  poolStats: PoolTimeSeriesDataResponse
}
type AggregatedData = {
  accepted: number[]
  staled: number[]
  rejected: number[]
}

function generateChartDataFromPoolStats(poolTimeSeriesRecord: any[], timestamps: Date[]) {
  // Default to current time if timestamps is null/undefined or empty
  const currentTime =
    timestamps && timestamps.length > 0 ? new Date(timestamps[timestamps.length - 1]).getTime() : new Date().getTime()
  const sixHoursAgo = currentTime - 6 * 60 * 60 * 1000

  const hourlyData: { [hour: string]: any[] } = {}

  if (poolTimeSeriesRecord && timestamps) {
    poolTimeSeriesRecord.forEach((record, index) => {
      const timestamp = timestamps[index]
      let timestampMillis: number

      if (timestamp instanceof Date) {
        timestampMillis = timestamp.getTime()
      } else if (typeof timestamp === 'string') {
        timestampMillis = new Date(timestamp).getTime()
      } else if (typeof timestamp === 'number') {
        timestampMillis = timestamp
      } else {
        console.error('Invalid timestamp format:', timestamp)
        return
      }

      if (timestampMillis < sixHoursAgo) {
        return
      }

      const hour = new Date(timestampMillis).getUTCHours()
      const key = `${hour}`

      if (!hourlyData[key]) {
        hourlyData[key] = []
      }

      hourlyData[key].push(record)
    })
  }

  const acceptedData = []
  const rejectedData = []
  const staledData = []
  const hours = []

  for (let offset = 0; offset < 6; offset++) {
    const hourAgoMillis = currentTime - offset * 60 * 60 * 1000
    const hourAgo = new Date(hourAgoMillis).getUTCHours()
    const key = `${hourAgo}`
    const records = hourlyData[key] || []

    const acceptedSum = records.reduce((sum, record) => sum + record.accepted, 0)
    const rejectedSum = records.reduce((sum, record) => sum + record.rejected, 0)
    const staledSum = records.reduce((sum, record) => sum + record.stale, 0)

    const dataCount = records.length || 1

    acceptedData.unshift(Math.round(acceptedSum / dataCount))
    rejectedData.unshift(Math.round(rejectedSum / dataCount))
    staledData.unshift(Math.round(staledSum / dataCount))
    hours.unshift(`${hourAgo}:00`)
  }

  console.log('accepted data ---->>>> ', acceptedData)

  return {
    data: [
      {
        name: 'Accepted',
        data: acceptedData.reverse()
      },
      {
        name: 'Staled',
        data: staledData.reverse()
      },
      {
        name: 'Rejected',
        data: rejectedData.reverse()
      }
    ],
    hours: hours.reverse()
  }
}

const RecentPoolStatsChart = ({ poolStats }: RecetPoolStatsChartProps) => {
  const { data, hours } = generateChartDataFromPoolStats(poolStats?.pool_time_series_record, poolStats?.timestamps)
  const theme = useTheme()

  const options: ApexOptions = {
    chart: {
      parentHeightOffset: 0,
      toolbar: { show: false }
    },
    colors: [theme.palette.primary.main, theme.palette.secondary.main, '#FF4136'],
    dataLabels: { enabled: false },
    plotOptions: {
      bar: {
        borderRadius: 3,
        barHeight: '50%',
        horizontal: true,
        startingShape: 'rounded'
      }
    },
    grid: {
      borderColor: theme.palette.divider,
      xaxis: {
        lines: { show: false }
      },
      padding: {
        top: -10
      }
    },
    yaxis: {
      labels: {
        style: { colors: theme.palette.text.primary }
      }
    },
    xaxis: {
      axisBorder: { show: false },
      axisTicks: { color: theme.palette.divider },
      categories: hours,
      labels: {
        style: { colors: theme.palette.text.primary }
      }
    },
    legend: {
      labels: {
        colors: [theme.palette.text.primary]
      }
    },
    tooltip: {
      custom: function (data) {
        const { series, seriesIndex, dataPointIndex, w } = data
        const seriesName = w.config.series[seriesIndex].name
        const value = series[seriesIndex][dataPointIndex]

        // Determine theme-based colors
        const tooltipColors =
          theme.palette.mode === 'dark'
            ? {
                text: '#FFFFFF',
                background: '#333333'
              }
            : {
                text: '#333333',
                background: '#FFFFFF'
              }

        // Return HTML for tooltip
        return `
          <div style="padding: 10px; color: ${tooltipColors.text}; background-color: ${tooltipColors.background}; border-radius: 4px;">
            <strong>${seriesName}</strong>
            <div>${value}</div>
          </div>
        `
      }
    }
  }

  return (
    <Card>
      <CardHeader
        title='Recent Pool Statistics'
        subheader='For the last 6 hours'
        sx={{
          flexDirection: ['column', 'row'],
          alignItems: ['flex-start', 'center'],
          '& .MuiCardHeader-action': { mb: 0 },
          '& .MuiCardHeader-content': { mb: [0, 0] }
        }}
      />
      <CardContent>
        <ReactApexcharts type='bar' height={400} options={options} series={data} />{' '}
      </CardContent>
    </Card>
  )
}

export default RecentPoolStatsChart
