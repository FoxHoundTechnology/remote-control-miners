import Box from '@mui/material/Box'

import { styled } from '@mui/material/styles'
import TimelineDot from '@mui/lab/TimelineDot'
import TimelineItem from '@mui/lab/TimelineItem'
import Typography from '@mui/material/Typography'
import TimelineContent from '@mui/lab/TimelineContent'
import TimelineSeparator from '@mui/lab/TimelineSeparator'
import TimelineConnector from '@mui/lab/TimelineConnector'
import MuiTimeline, { TimelineProps } from '@mui/lab/Timeline'

import { Trigger, TriggerHistory } from 'src/types/apps/triggerTypes'
import { useEffect, useState } from 'react'
import { FormControl, InputLabel, Select, MenuItem, Button, Modal } from '@mui/material'
import { generateActionSentence } from './utils'
import { useDispatch } from 'react-redux'
import { AppDispatch } from 'src/store'
import { DeleteTriggerHistories } from 'src/store/apps/trigger'

// Styled Timeline component
const Timeline = styled(MuiTimeline)<TimelineProps>({
  paddingLeft: 0,
  paddingRight: 0,
  '& .MuiTimelineItem-root': {
    width: '100%',
    '&:before': {
      display: 'none'
    }
  }
})

// Styled component for the image of a shoe
const ImgShoe = styled('img')(({ theme }) => ({
  borderRadius: theme.shape.borderRadius
}))

type TimelineEntry = {
  dotColor: 'error' | 'primary' | 'warning' | 'success'
  title: string
  date: string
  description: string
  details?: {
    iconSrc?: string
    iconAlt?: string
    fileName?: string
    imageSrc?: string
    imageAlt?: string
    person?: {
      avatarSrc: string
      name: string
      role?: string
    }
    actions?: {
      messageIcon?: string
      phoneIcon?: string
    }
    product?: {
      customer: string
      price: string
      quantity: number
    }
  }
}

interface TriggerTimelineProps {
  triggers: Trigger[]
}

const itemsPerPage = 5

const mapTriggerHistoryToTimelineEntry = (triggerHistory: TriggerHistory, trigger: Trigger): TimelineEntry => {
  const dateObj = new Date(triggerHistory.timestamp)

  // Extracting time in am/pm format
  let hours = dateObj.getHours()
  const minutes = dateObj.getMinutes()
  const ampm = hours >= 12 ? 'pm' : 'am'
  hours = hours % 12
  hours = hours ? hours : 12 // The hour '0' should be '12'
  const strMinutes = minutes < 10 ? '0' + minutes : minutes.toString()
  const strTime = `${hours}:${strMinutes} ${ampm}`

  // Extracting day, month, and year
  const days: string[] = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']
  const months: string[] = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec']
  const formattedDate = `${strTime}, ${days[dateObj.getDay()]} ${
    months[dateObj.getMonth()]
  } ${dateObj.getDate()} ${dateObj.getFullYear()}`

  return {
    dotColor: 'success',
    title: `${trigger.name}`,
    date: formattedDate,
    description: triggerHistory.message
  }
}

const TriggerTimeline = ({ triggers }: TriggerTimelineProps) => {
  // state for selected trigger
  const [selectedTriggerId, setSelectedTriggerId] = useState<number | null>(null)
  // modal component
  const [isConfirmModalOpen, setIsConfirmModalOpen] = useState(false)

  useEffect(() => {
    if (triggers.length > 0 && selectedTriggerId === null) {
      setSelectedTriggerId(triggers[0].ID)
    }
  }, [triggers])

  const selectedTrigger = triggers.find(trigger => trigger.ID === selectedTriggerId)
  const filteredHistories = selectedTrigger?.histories || []

  const timelineEntries = [...filteredHistories]
    .reverse()
    .map(history => mapTriggerHistoryToTimelineEntry(history, selectedTrigger!))

  const [currentPage, setCurrentPage] = useState<number>(1)

  // pagination logic based on the timeline entries
  const paginatedEntries = timelineEntries.slice((currentPage - 1) * itemsPerPage, currentPage * itemsPerPage)
  const totalPages = Math.ceil(timelineEntries.length / itemsPerPage)

  // redux logic
  const dispatch = useDispatch<AppDispatch>()

  return (
    <Box>
      <Box display='flex' justifyContent='flex-end' mb={2}>
        <FormControl hiddenLabel variant='outlined' style={{ width: '250px' }}>
          <Select
            labelId='trigger-filter-label'
            value={selectedTriggerId}
            onChange={event => setSelectedTriggerId(event.target.value as number)}
            label='Filter by Trigger'
          >
            {triggers.map(trigger => (
              <MenuItem key={trigger.ID} value={trigger.ID}>
                {trigger.name}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </Box>
      <Timeline>
        {paginatedEntries.map((entry, index) => (
          <TimelineItem key={index}>
            <TimelineSeparator>
              <TimelineDot color={entry.dotColor} />
              <TimelineConnector />
            </TimelineSeparator>
            <TimelineContent>
              <Box
                sx={{ mb: 2, display: 'flex', flexWrap: 'wrap', alignItems: 'center', justifyContent: 'space-between' }}
              >
                <Typography variant='body2' sx={{ mr: 2, fontWeight: 600, color: 'text.primary' }}>
                  {entry.title}
                </Typography>
                <Typography variant='caption'>{entry.date}</Typography>
              </Box>
              <Typography
                variant='body2'
                sx={{ color: 'text.primary' }}
                dangerouslySetInnerHTML={{ __html: entry.description }}
              />
              {entry.details && (
                <>
                  {entry.details.iconSrc && (
                    <Box sx={{ mt: 2, display: 'flex', alignItems: 'center' }}>
                      <img width={28} height={28} alt={entry.details.iconAlt} src={entry.details.iconSrc} />
                      <Typography variant='subtitle2' sx={{ ml: 2, fontWeight: 600 }}>
                        {entry.details.fileName}
                      </Typography>
                    </Box>
                  )}
                  {/* You can continue with similar conditional renderings for other details like `imageSrc`, `person`, `actions`, `product`, etc. */}
                </>
              )}
            </TimelineContent>
          </TimelineItem>
        ))}
      </Timeline>
      <Box display='flex' justifyContent='center' mt={2} alignItems='center'>
        <Button
          disabled={currentPage === 1}
          onClick={() => setCurrentPage(currentPage - 1)}
          sx={{
            border: '1px solid',
            borderColor: currentPage === 1 ? 'grey.300' : 'primary.main',
            mr: 2 // Add margin to the right
          }}
        >
          Previous
        </Button>
        <Box mx={2}>
          {currentPage} / {totalPages}
        </Box>
        <Button
          disabled={currentPage === totalPages}
          onClick={() => setCurrentPage(currentPage + 1)}
          sx={{
            border: '1px solid',
            borderColor: currentPage === totalPages ? 'grey.300' : 'primary.main',
            ml: 2 // Add margin to the left
          }}
        >
          Next
        </Button>
      </Box>
      <Box mt={{ xs: 4, sm: 0 }}>
        {/* Clear History Button */}
        <Button onClick={() => setIsConfirmModalOpen(true)} sx={{ border: '1px solid' }}>
          Clear History
        </Button>
      </Box>

      <Modal
        open={isConfirmModalOpen}
        onClose={() => setIsConfirmModalOpen(false)}
        aria-labelledby='modal-title'
        aria-describedby='modal-description'
      >
        <Box
          sx={{
            width: 400,
            p: 4,
            mx: 'auto',
            mt: '20vh',
            bgcolor: 'background.paper',
            borderRadius: 2,
            boxShadow: 3
          }}
        >
          <Typography id='modal-title' variant='h6' component='h2'>
            Clear History
          </Typography>
          <Typography id='modal-description' sx={{ mt: 2 }}>
            Are you sure you want to clear the history?
          </Typography>
          <Box mt={3} display='flex' justifyContent='flex-end'>
            <Button
              onClick={() => setIsConfirmModalOpen(false)}
              sx={{ mr: 2, border: '1px solid', borderColor: 'divider' }}
            >
              Cancel
            </Button>
            <Button
              onClick={() => {
                // Add logic to clear the history here
                dispatch(DeleteTriggerHistories(Number(selectedTriggerId)))
                setIsConfirmModalOpen(false)
              }}
              color='error'
              sx={{ border: '1px solid', borderColor: 'error.main' }} // Added border here
            >
              Confirm
            </Button>
          </Box>
        </Box>
      </Modal>
    </Box>
  )
}

const triggerHistoryMockData: TriggerHistory[] = [
  {
    ID: 1,
    timestamp: new Date('2023-01-05T12:00:00Z'),
    message: 'Triggered alert for high temperature',
    trigger_id: 101
  },
  {
    ID: 2,
    timestamp: new Date('2023-01-10T15:30:00Z'),
    message: 'Triggered emergency shutdown',
    trigger_id: 102
  },
  {
    ID: 3,
    timestamp: new Date('2023-01-15T09:45:00Z'),
    message: 'Triggered maintenance alert',
    trigger_id: 103
  },
  {
    ID: 4,
    timestamp: new Date('2023-01-20T17:20:00Z'),
    message: 'Triggered alert for low pressure',
    trigger_id: 104
  },
  {
    ID: 5,
    timestamp: new Date('2023-01-25T14:10:00Z'),
    message: 'Triggered security breach alert',
    trigger_id: 105
  }
]

export default TriggerTimeline
