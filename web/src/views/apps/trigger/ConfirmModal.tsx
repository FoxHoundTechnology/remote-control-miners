import { SetStateAction, Dispatch } from 'react'

import Dialog from '@mui/material/Dialog'
import Button from '@mui/material/Button'
import Typography from '@mui/material/Typography'
import DialogContent from '@mui/material/DialogContent'
import DialogActions from '@mui/material/DialogActions'

type ConfirmModalProps = {
  title: string
  message?: string
  onConfirm: () => void
  show: boolean
  setShow: (show: boolean) => void
}

const ConfirmModal = ({ title, message, onConfirm, show, setShow }: ConfirmModalProps) => {
  return (
    <Dialog
      fullWidth
      open={show}
      maxWidth='sm'
      scroll='body'
      onClose={() => {
        setShow(false)
      }}
      // TransitionComponent={Transition}
    >
      <DialogContent
        sx={{
          position: 'relative',
          pb: theme => `${theme.spacing(8)} !important`,
          px: theme => [`${theme.spacing(5)} !important`, `${theme.spacing(15)} !important`],
          pt: theme => [`${theme.spacing(8)} !important`, `${theme.spacing(12.5)} !important`]
        }}
      >
        <Typography variant='h5' sx={{ mb: 4 }}>
          {title}
        </Typography>
        <Typography sx={{ mb: 3 }}>{message}</Typography>
      </DialogContent>
      <DialogActions
        sx={{
          justifyContent: 'center',
          px: theme => [`${theme.spacing(5)} !important`, `${theme.spacing(15)} !important`],
          pb: theme => [`${theme.spacing(8)} !important`, `${theme.spacing(12.5)} !important`]
        }}
      >
        <Button
          variant='contained'
          sx={{ mr: 1 }}
          onClick={() => {
            onConfirm()
            setShow(false)
          }}
        >
          Confirm
        </Button>
        <Button
          variant='outlined'
          color='secondary'
          onClick={() => {
            setShow(false)
          }}
        >
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  )
}

export default ConfirmModal
