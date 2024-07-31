import React, { useRef, useState } from 'react'
import Box from '@mui/material/Box'
import Card from '@mui/material/Card'
import Grid from '@mui/material/Grid'
import CardHeader from '@mui/material/CardHeader'
import Typography from '@mui/material/Typography'
import CardContent from '@mui/material/CardContent'
import IconButton from '@mui/material/IconButton'
import ExpandMoreIcon from '@mui/icons-material/ExpandMore'
import ExpandLessIcon from '@mui/icons-material/ExpandLess'
import FullscreenIcon from '@mui/icons-material/Fullscreen'
import CloseIcon from '@mui/icons-material/Close'
import Modal from '@mui/material/Modal'
import { useTheme } from '@mui/material'

interface MinerLogProps {
  log: string
}

const MinerLog: React.FC<MinerLogProps> = ({ log }) => {
  const theme = useTheme()

  const logLines = log.split('\n')
  const cardBoxRef = useRef<HTMLDivElement>(null)
  const modalBoxRef = useRef<HTMLDivElement>(null)
  const [isModalOpen, setIsModalOpen] = useState(false)

  const scrollToBottom = (ref: React.RefObject<HTMLDivElement>) => {
    if (ref.current) {
      ref.current.scrollTop = ref.current.scrollHeight
    }
  }

  const scrollToTop = (ref: React.RefObject<HTMLDivElement>) => {
    if (ref.current) {
      ref.current.scrollTop = 0
    }
  }

  const openModal = () => setIsModalOpen(true)
  const closeModal = () => setIsModalOpen(false)

  const LogContent = ({ boxRef }: { boxRef: React.RefObject<HTMLDivElement> }) => (
    <Box
      ref={boxRef}
      sx={{
        height: '100%',
        overflowY: 'auto',
        padding: 2,
        backgroundColor: '#1e1e1e',
        fontFamily: 'Consolas, Monaco, "Courier New", monospace',
        whiteSpace: 'pre-wrap',
        wordBreak: 'break-all',
        '&::-webkit-scrollbar': {
          width: '10px'
        },
        '&::-webkit-scrollbar-track': {
          background: '#333333'
        },
        '&::-webkit-scrollbar-thumb': {
          background: '#666666',
          borderRadius: '5px'
        },
        '&::-webkit-scrollbar-thumb:hover': {
          background: '#888888'
        }
      }}
    >
      {logLines.map((line, index) => (
        <Typography
          key={index}
          variant='body2'
          component='div'
          sx={{
            fontSize: '0.9rem',
            lineHeight: 1.5,
            mb: 0.5,
            color: '#ffffff',
            fontFamily: 'Consolas, Monaco, "Courier New", monospace'
          }}
        >
          {line}
        </Typography>
      ))}
    </Box>
  )

  const ScrollButtons = ({ boxRef }: { boxRef: React.RefObject<HTMLDivElement> }) => (
    <>
      <IconButton onClick={() => scrollToTop(boxRef)} size='small' sx={{ color: theme.palette.grey[500] }}>
        <ExpandLessIcon />
      </IconButton>
      <IconButton onClick={() => scrollToBottom(boxRef)} size='small' sx={{ color: theme.palette.grey[500] }}>
        <ExpandMoreIcon />
      </IconButton>
    </>
  )

  return (
    <Grid container spacing={6}>
      <Grid item xs={12}>
        <Card>
          <CardHeader
            title={<Typography variant='h6'>Miner Logs</Typography>}
            sx={{
              // backgroundColor: '#2c2c2c',
              borderBottom: '1px solid #444444'
            }}
            action={
              <Box>
                <ScrollButtons boxRef={cardBoxRef} />
                <IconButton onClick={openModal} size='small'>
                  <FullscreenIcon />
                </IconButton>
              </Box>
            }
          />
          <CardContent
            sx={{
              height: '400px',
              padding: 2,
              backgroundColor: '#1e1e1e'
            }}
          >
            <LogContent boxRef={cardBoxRef} />
          </CardContent>
        </Card>
      </Grid>

      <Modal
        open={isModalOpen}
        onClose={closeModal}
        aria-labelledby='modal-modal-title'
        aria-describedby='modal-modal-description'
      >
        <Box
          sx={{
            position: 'absolute',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
            width: '90%',
            height: '90%',
            bgcolor: '#1e1e1e',
            border: '2px solid #000',
            boxShadow: 24,
            p: 4
          }}
        >
          <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 2 }}>
            <Typography id='modal-modal-title' variant='h6' component='h2' sx={{ color: '#ffffff' }}>
              Recent Miner Logs
            </Typography>
            <Box>
              <ScrollButtons boxRef={modalBoxRef} />
              <IconButton onClick={closeModal}>
                <CloseIcon sx={{ color: theme.palette.grey[500] }} />
              </IconButton>
            </Box>
          </Box>
          <Box sx={{ height: 'calc(100% - 40px)' }}>
            <LogContent boxRef={modalBoxRef} />
          </Box>
        </Box>
      </Modal>
    </Grid>
  )
}

export default MinerLog
