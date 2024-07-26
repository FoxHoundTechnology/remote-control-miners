import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import Grid from '@mui/material/Grid'
import { styled } from '@mui/material/styles'
import TimelineDot from '@mui/lab/TimelineDot'
import TimelineItem from '@mui/lab/TimelineItem'
import CardHeader from '@mui/material/CardHeader'
import Typography from '@mui/material/Typography'
import CardContent from '@mui/material/CardContent'
import TimelineContent from '@mui/lab/TimelineContent'
import TimelineSeparator from '@mui/lab/TimelineSeparator'
import TimelineConnector from '@mui/lab/TimelineConnector'
import MuiTimeline, { TimelineProps } from '@mui/lab/Timeline'

// Custom Components Imports
import CustomChip from 'src/@core/components/mui/chip'

// Styled Timeline component
const Timeline = styled(MuiTimeline)<TimelineProps>(({ theme }) => ({
  margin: 0,
  padding: 0,
  marginLeft: theme.spacing(0.75),
  '& .MuiTimelineItem-root': {
    '&:before': {
      display: 'none'
    },
    '&:last-child': {
      minHeight: 60
    }
  }
}))

type Activity = {
  type: 'mode' | 'reboot' | 'config' | 'warning' | 'alert' | 'disabled'
  title: string
  description: string
  timestamp: Date
}

const mockData: Activity[] = [
  {
    type: 'mode',
    title: 'Power Mode',
    description: 'Switched to Power Mode at 2:12pm',
    timestamp: new Date('2023-08-15T14:12:00Z')
  },
  {
    type: 'reboot',
    title: 'Reboot',
    description: 'Rebooted 10:15am',
    timestamp: new Date('2023-08-15T10:15:00Z')
  },
  {
    type: 'reboot',
    title: 'Reboot',
    description: 'Rebooted at 9:30am',
    timestamp: new Date('2023-08-13T08:00:00Z')
  },
  {
    type: 'warning',
    title: 'Low Hashrate',
    description: 'Low Hashrate detected at 5:30pm',
    timestamp: new Date('2023-08-14T17:30:00Z')
  },
  {
    type: 'alert',
    title: 'Low Hashrate',
    description: 'Low hashrate detected at 8:45pm',
    timestamp: new Date('2023-08-14T20:45:00Z')
  }
]

const getActivityColor = (type: Activity['type']) => {
  switch (type) {
    case 'mode':
      return 'primary' // Color choice for 'mode'
    case 'reboot':
      return 'secondary' // Color choice for 'reboot'
    case 'config':
      return 'warning' // Color choice for 'config'
    case 'warning':
      return 'error' // Color choice for 'warning'
    case 'alert':
      return 'error' // Color choice for 'alert'
    case 'disabled':
      return 'grey' // Color choice for 'disabled'
    default:
      return 'grey'
  }
}

const getTimeDifference = (timestamp: Date) => {
  const diffInMs = Date.now() - timestamp.getTime()
  const diffInMinutes = Math.floor(diffInMs / (1000 * 60))
  if (diffInMinutes < 60) return `${diffInMinutes} min ago`
  const diffInHours = Math.floor(diffInMinutes / 60)
  if (diffInHours < 24) return `${diffInHours} hours ago`
  const diffInDays = Math.floor(diffInHours / 24)
  return `${diffInDays} days ago`
}

const renderTimelineItems = (activities: Activity[]) => {
  return activities.map((activity, index) => (
    <TimelineItem key={index}>
      <TimelineSeparator>
        {/* Reduce the size of the TimelineDot */}
        <TimelineDot color={getActivityColor(activity.type)} />
        {index !== activities.length - 1 && <TimelineConnector />}
      </TimelineSeparator>
      <TimelineContent>
        {/* Adjust font sizes using the sx prop */}
        <Typography component='span' sx={{ fontSize: '0.8rem', mt: -1 }}>
          {activity.title}
        </Typography>
        <Typography sx={{ fontSize: '0.8rem' }}>{activity.description}</Typography>
        <Typography variant='body2' sx={{ fontSize: '0.7rem' }}>
          {getTimeDifference(activity.timestamp)}
        </Typography>
      </TimelineContent>
    </TimelineItem>
  ))
}

const ActivityLogs = () => (
  <Grid container spacing={6}>
    <Grid item xs={12}>
      <Card>
        <CardHeader title='Activity Logs' />
        <CardContent>
          <CustomChip
            skin='light'
            size='small'
            color='primary'
            label={'Activity Logs (Under Development)'}
            sx={{ height: 20, fontSize: '0.75rem', fontWeight: 500, borderRadius: '10px' }}
          />
          <Box mb={3} />
          <Timeline>{renderTimelineItems(mockData)}</Timeline>
        </CardContent>
      </Card>
    </Grid>
  </Grid>
)

export default ActivityLogs
