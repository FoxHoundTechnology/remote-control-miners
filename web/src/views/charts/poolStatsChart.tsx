// slightly different from recentPoolStats chart
// this one is for the last 24 hours with an hourly breakdown
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
import { PoolStatsData } from 'src/types/apps/minerDetailsTypes'

interface PickerProps {
  start: Date | number
  end: Date | number
}

function generateChartDataFromPoolStats(poolTimeSeriesRecord: PoolStatsData[], timestamps: string[]) {
  const currentTime = new Date(timestamps[timestamps.length - 1]).getTime() || new Date().getTime()

  const acceptedData = []
  const rejectedData = []
  const staledData = []
  const hours = []

  for (let offset = 0; offset < 24; offset++) {
    const hourAgoMillis = currentTime - offset * 60 * 60 * 1000
    const hourAgo = new Date(hourAgoMillis).getUTCHours()

    acceptedData.unshift(poolTimeSeriesRecord[offset]?.accepted)
    rejectedData.unshift(poolTimeSeriesRecord[offset]?.rejected)
    staledData.unshift(poolTimeSeriesRecord[offset]?.stale)

    hours.unshift(`${hourAgo}:00`)
  }

  console.log('modified array for accept', acceptedData)

  return {
    data: [
      {
        name: 'Accepted',
        data: acceptedData
      },
      {
        name: 'Staled',
        data: staledData
      },
      {
        name: 'Rejected',
        data: rejectedData
      }
    ],
    hours
  }
}

type PoolStatsChartProps = {
  poolStatsArr: PoolStatsData[]
  timeStampArr: string[]
}

const PoolStatsChart = ({ poolStatsArr, timeStampArr }: PoolStatsChartProps) => {
  console.log('what is the pool chart data', poolStatsArr)
  console.log('what is the pool chart data', timeStampArr)

  const [endDate, setEndDate] = useState<DateType>(null)
  const [startDate, setStartDate] = useState<DateType>(null)

  const { data, hours } = generateChartDataFromPoolStats(poolStatsArr, timeStampArr)

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
        columnWidth: '80%', // This will create a gap between each bar.
        horizontal: false,
        startingShape: 'rounded'
      }
    },
    grid: {
      borderColor: theme.palette.divider,
      xaxis: {
        // Grid for x-axis.
        lines: {
          show: true
        }
      },
      yaxis: {
        // Enable the grid for y-axis.
        lines: {
          show: true
        }
      },
      padding: {
        top: -10
      }
    },
    xaxis: {
      // Categories will now be on the x-axis
      axisBorder: { show: false },
      axisTicks: { color: theme.palette.divider },
      categories: hours,

      labels: {
        style: { colors: theme.palette.text.primary }
      }
    },
    yaxis: {
      // Values will now be on the y-axis
      labels: {
        style: { colors: theme.palette.text.primary }
      },
      min: 0
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

  const CustomInput = forwardRef((props: PickerProps, ref) => {
    const startDate = props.start !== null ? format(props.start, 'MM/dd/yyyy') : ''
    const endDate = props.end !== null ? ` - ${format(props.end, 'MM/dd/yyyy')}` : null

    const value = `${startDate}${endDate !== null ? endDate : ''}`

    return (
      <TextField
        {...props}
        size='small'
        value={value}
        inputRef={ref}
        InputProps={{
          startAdornment: (
            <InputAdornment position='start'>
              <Icon icon='mdi:bell-outline' />
            </InputAdornment>
          ),
          endAdornment: (
            <InputAdornment position='end'>
              <Icon icon='mdi:chevron-down' />
            </InputAdornment>
          )
        }}
      />
    )
  })

  const handleOnChange = (dates: any) => {
    const [start, end] = dates
    setStartDate(start)
    setEndDate(end)
  }

  return (
    <Card>
      <CardHeader
        title='Pool Statistics'
        subheader='For the last 24 hours'
        sx={{
          flexDirection: ['column', 'row'],
          alignItems: ['flex-start', 'center'],
          '& .MuiCardHeader-action': { mb: 0 },
          '& .MuiCardHeader-content': { mb: [0, 0] }
        }}
      />
      <CardContent>
        <ReactApexcharts type='bar' height={400} options={options} series={data} />
      </CardContent>
    </Card>
  )
}

export default PoolStatsChart
